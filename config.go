package main

import (
	"github.com/hashicorp/logutils"
)

type Config struct {
	Ignored         []string          `json:"ignored"`
	ClientID        string            `json:"client_id"`
	ClientSecret    string            `json:"client_secret"`
	ApplicationURI  string            `json:"application_uri"`
	MyURI           string            `json:"my_uri"`
	CertificateFile string            `json:"cert"`
	KeyFile         string            `json:"key"`
	LogLevel        logutils.LogLevel `json:"log_level"`
}
