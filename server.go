package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/logutils"
	"github.com/zmb3/spotify"
)

var (
	config Config
	client *spotify.Client
	// myURI is the OAuth redirect URI for the application.
	// You must register an application at Spotify's developer portal
	// and enter this value.
	myURI = "https://localhost:5009/callback"
	// applicationURI is the application's uri where the final token is sent
	applicationURI = "https://localhost:5007/callback"
	auth           spotify.Authenticator
	ch             = make(chan *spotify.Client)
	certificate    = "cert.pem"
	key            = "key.pem"
)

func completeAuth(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	// "access_token":"BQAyVReY-wxNd3K2kKLuMntyArFsRLoW6ahCMh_BphpojTFZe_EbW4t9bUWgQIRr5mFKhdODqXL_pA6uIhGKOae3aKllzpQVA0H7RlCumaN2NJAaSmw7y13fTRwEnLyNIyp9HMuMx7b2Y2ze6aM"
	// "token_type":"Bearer"
	// "refresh_token":"AQDMcU9J_7SVspLyXXvn-HvgW-Ust2tGr2Wep4OU1bxbJHg9KCTrc9X2SCbQJsidn2Ye5SG9SXPPD4QF1c3rQggvD6_u_AGM891mBxnYXGgo3jBnAgwPBBL-eXUIM79FlIQ"
	// "expiry":"2017-08-11T21:21:40.806561311-04:00"
	http.Redirect(w, r, applicationURI+fmt.Sprintf("?access_token=%s&token_type=%s&refresh_token=%s&expiry=%d&state=%s", tok.AccessToken, tok.TokenType, tok.RefreshToken, int(time.Until(tok.Expiry).Seconds()), state), 302)
}

func main() {
	logFilter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(logFilter)

	conf, err := os.Open("config.json")
	if err != nil {
		log.Print("[DEBUG] No config file specified, ignoring.")
	} else {
		defer conf.Close()

		decoder := json.NewDecoder(conf)
		err = decoder.Decode(&config)
		if err != nil {
			log.Fatalf("Config file 'config.json could not be read, %v", err)
		}
		if config.LogLevel != "" {
			logFilter.SetMinLevel(config.LogLevel)
		}
		if config.ApplicationURI != "" {
			applicationURI = config.ApplicationURI
		}
		if config.MyURI != "" {
			myURI = config.MyURI
		}
		if config.CertificateFile != "" {
			certificate = config.CertificateFile
		}
		if config.KeyFile != "" {
			key = config.KeyFile
		}
	}

	router := mux.NewRouter()

	router.HandleFunc("/callback", completeAuth).Methods("GET")
	auth = spotify.NewAuthenticator(myURI)
	if config.ClientID != "" && config.ClientSecret != "" {
		auth.SetAuthInfo(config.ClientID, config.ClientSecret)
	}

	log.Fatal(http.ListenAndServeTLS(":5009", certificate, key, router))
}
