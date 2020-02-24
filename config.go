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
	BaseURI         string            `json:"base_uri"`
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
			baseURI = os.Getenv("ENV")
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
		config.ClientID = os.Getenv("CLOVER_ID")
		config.ClientSecret = os.Getenv("CLOVER_SECRET")

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
		if config.Env != "" {
			env = config.Env
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
