package main

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Node []struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"node"`
	InboundTag      string `json:"inboundTag"`
	Flow            string `json:"flow"`
	CertificateFile string `json:"certificateFile"`
}

type UserInfo []struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	Level uint32 `json:"level"`
}

func decodeConfiguration(filePath string) Configuration {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var config Configuration
	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		panic(err)
	}
	return config
}

func decodeUserInfo(filePath string) UserInfo {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var user UserInfo
	err = json.Unmarshal([]byte(content), &user)
	if err != nil {
		panic(err)
	}
	return user
}
