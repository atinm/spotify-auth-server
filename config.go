package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

type Config struct {
	Ignored         []string          `json:"ignored"`
	ClientID        string            `json:"client_id"`
	ClientSecret    string            `json:"client_secret"`
	ApplicationURI  string            `json:"application_uri"`
	BaseURI         string            `json:"base_uri"`
	CertificateFile string            `json:"cert"`
	KeyFile         string            `json:"key"`
	LogLevel        logutils.LogLevel `json:"log_level"`
	Port            string            `json:"port"`
}

func LoadConfig() {
	conf, err := os.Open("config.json")
	if err != nil {
		if os.Getenv("LOG_LEVEL") != "" {
			logFilter.SetMinLevel(logutils.LogLevel(os.Getenv("LOG_LEVEL")))
		}
		log.Print("[DEBUG] No config file specified, reading environment variables.")

		if os.Getenv("APPLICATION_URI") != "" {
			applicationURI = os.Getenv("APPLICATION_URI")
		}
		if os.Getenv("BASE_URI") != "" {
			baseURI = os.Getenv("BASE_URI")
		}
		// if os.Getenv("CERTIFICATE") != "" {
		// 	certificate = os.Getenv("CERTIFICATE")
		// }
		// if os.Getenv("KEY") != "" {
		// 	key = os.Getenv("KEY")
		// }
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}
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
		if config.BaseURI != "" {
			baseURI = config.BaseURI
		}
		// if config.CertificateFile != "" {
		// 	certificate = config.CertificateFile
		// }
		// if config.KeyFile != "" {
		// 	key = config.KeyFile
		// }
		if config.Port != "" {
			port = config.Port
		}
	}
}
