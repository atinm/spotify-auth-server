package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/atinm/spotify"
	"github.com/gorilla/mux"
	"github.com/hashicorp/logutils"
)

var (
	config Config
	client *spotify.Client
	// myURI is the OAuth redirect URI for the application.
	// You must register an application at Spotify's developer portal
	// and enter this value.
	baseURI = "https://localhost"
	// applicationURI is the application's uri where the final token is sent
	applicationURI = "https://localhost:5007/callback"
	// TokenURL is the URL to the Spotify Accounts Service's OAuth2
	// token endpoint.
	spotifyTokenURL = "https://accounts.spotify.com/api/token"
	auth            spotify.Authenticator
	ch              = make(chan *spotify.Client)
	certificate     = "cert.pem"
	key             = "key.pem"
	port            = "5009"
	logFilter       *logutils.LevelFilter
)

func completeAuth(w http.ResponseWriter, r *http.Request) {
	log.Print("[DEBUG] Received callback")
	state := r.FormValue("state")
	tok, err := auth.Token(state, r)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		return
	}
	// "access_token":"BQAyVReY-wxNd3K2kKLuMntyArFsRLoW6ahCMh_BphpojTFZe_EbW4t9bUWgQIRr5mFKhdODqXL_pA6uIhGKOae3aKllzpQVA0H7RlCumaN2NJAaSmw7y13fTRwEnLyNIyp9HMuMx7b2Y2ze6aM"
	// "token_type":"Bearer"
	// "refresh_token":"AQDMcU9J_7SVspLyXXvn-HvgW-Ust2tGr2Wep4OU1bxbJHg9KCTrc9X2SCbQJsidn2Ye5SG9SXPPD4QF1c3rQggvD6_u_AGM891mBxnYXGgo3jBnAgwPBBL-eXUIM79FlIQ"
	// "expiry":"2017-08-11T21:21:40.806561311-04:00"
	redirect := applicationURI + fmt.Sprintf("?access_token=%s&token_type=%s&refresh_token=%s&expiry=%d&state=%s", tok.AccessToken, tok.TokenType, tok.RefreshToken, int(time.Until(tok.Expiry).Seconds()), state)
	log.Print("[DEBUG] ", redirect)
	http.Redirect(w, r, redirect, 302)
}

func refreshTokenReq(w http.ResponseWriter, r *http.Request) {
	log.Print("[DEBUG] Received refreshToken request")
	clientId, clientSecret, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Couldn't get basicAuth", http.StatusForbidden)
		return
	}

	if clientId != config.ClientID || clientSecret != "" {
		http.Error(w, "Couldn't get refresh token: clientId("+clientId+") != "+config.ClientID+", clientSecret("+clientSecret+") != \"\"", http.StatusBadRequest)
		return
	}
	form := url.Values{
		"grant_type":    {r.FormValue("grant_type")},
		"refresh_token": {r.FormValue("refresh_token")},
	}

	req, err := http.NewRequest("POST", spotifyTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("[ERROR] Request failed: %v", err)
		http.Error(w, "Couldn't get refresh token", http.StatusBadRequest)
		return
	}
	req.SetBasicAuth(config.ClientID, config.ClientSecret)
	b := config.ClientID + ":" + config.ClientSecret
	enc := base64.StdEncoding.EncodeToString([]byte(b))
	req.Header.Set("Authorization", "Basic "+enc)
	log.Printf("[DEBUG] Set (%s) Header Authorization: %s", b, req.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Request failed: %v", err)
		http.Error(w, "Couldn't get refresh token", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Request failed, could not read response: %v", err)
		http.Error(w, "Couldn't get read response", resp.StatusCode)
		return
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func main() {
	logFilter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(logFilter)

	LoadConfig()

	router := mux.NewRouter()

	router.HandleFunc("/callback", completeAuth).Methods("GET")
	auth = spotify.NewAuthenticator(baseURI + "/callback")
	if config.ClientID != "" && config.ClientSecret != "" {
		auth.SetAuthInfo(config.ClientID, config.ClientSecret)
	}

	router.HandleFunc("/token", refreshTokenReq).Methods("POST")

	log.Printf("[DEBUG] listening on %s, with internal port %s", baseURI+"/callback", port)
	//log.Fatal(http.ListenAndServeTLS(":"+port, certificate, key, router))
	log.Fatal(http.ListenAndServe(":"+port, router))
}
