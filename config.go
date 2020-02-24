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
	CertificateFile string            `json:"cert"`
	KeyFile         string            `json:"key"`
	LogLevel        logutils.LogLevel `json:"log_level"`
	Port            string            `json:"port"`
	Env             string            `json:"env"`
}

const (
	SandboxServer = "https://sandbox.dev.clover.com"
	ProductionServer = "https://www.clover.com"
	TokenEndpoint = "/oauth/token"
	Port = "5009"
)

func TokenURL() string {
	tokenURL := SandboxServer + TokenEndpoint
	if config.Env == "production" {
		tokenURL = ProductionServer + TokenEndpoint
	}
	return tokenURL
}

func LoadConfig() {
	conf, err := os.Open("config.json")
	if err != nil {
		if os.Getenv("LOG_LEVEL") != "" {
			logFilter.SetMinLevel(logutils.LogLevel(os.Getenv("LOG_LEVEL")))
		}
		log.Print("[DEBUG] No config file specified, reading environment variables.")

		if os.Getenv("ENV") != "" {
			config.Env = os.Getenv("ENV")
		}
		
		if os.Getenv("PORT") != "" {
			config.Port = os.Getenv("PORT")
		} else {
			config.Port = Port
		}
		config.ClientID = os.Getenv("CLIENT_ID")
		config.ClientSecret = os.Getenv("CLIENT_SECRET")
		if os.Getenv("CERTIFICATE") != "" {
			config.CertificateFile = os.Getenv("CERTIFICATE")
		}
		if os.Getenv("KEY") != "" {
			config.KeyFile = os.Getenv("KEY")
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
	}
}
