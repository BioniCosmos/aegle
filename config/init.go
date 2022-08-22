package config

import (
    "encoding/json"
    "errors"
    "log"
    "os"
)

type config struct {
    Listen      string `json:"listen"`
    DatabaseURL string `json:"databaseURL"`
}

var Conf config

func Init() {
    file, err := os.ReadFile("./config.json")
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
