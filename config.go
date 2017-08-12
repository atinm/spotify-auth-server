package main

import (
	"github.com/hashicorp/logutils"
)

type Config struct {
	Ignored      []string          `json:"ignored"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
	LogLevel     logutils.LogLevel `json:"log_level"`
}
