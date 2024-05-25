package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Listen       string `json:"listen,omitempty"`
	DatabaseURL  string `json:"databaseURL,omitempty"`
	DatabaseName string `json:"databaseName,omitempty"`
	Home         string `json:"home,omitempty"`
	Dashboard    string `json:"dashboard,omitempty"`
	XrayConfig   string `json:"xrayConfig,omitempty"`
}

var C Config

func Init(path string, server string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(file, &C); err != nil {
		log.Fatal(err)
	}
	if server == "center" {
		if C.DatabaseURL == "" {
			log.Fatal("parsing config: no databaseURL specified")
		}
		if C.DatabaseName == "" {
			C.DatabaseName = "aegle"
		}
	} else {
		if C.XrayConfig == "" {
			C.XrayConfig = "/usr/local/etc/xray/inbounds.json"
		}
	}
}
