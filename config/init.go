package config

import (
    "encoding/json"
    "errors"
    "flag"
    "log"
    "os"
)

type config struct {
    Listen      string `json:"listen"`
    DatabaseURL string `json:"databaseURL"`
    Cert        string `json:"cert"`
    Static      struct {
        Home      string `json:"home"`
        Dashboard string `json:"dashboard"`
    }   `json:"static"`
}

var Conf config

func Init() {
    configFlag := flag.String("config", "./config.json", "the configuration file")
    flag.Parse()
    file, err := os.ReadFile(*configFlag)
    if err != nil {
        log.Fatal(err)
    }
    if err := json.Unmarshal(file, &Conf); err != nil {
        log.Fatal(err)
    }
    if Conf.DatabaseURL == "" {
        log.Fatal(errors.New("parsing config: no databaseURL specified"))
    }
}
