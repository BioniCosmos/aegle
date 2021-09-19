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

type Subscriptions []struct {
	UserInfo struct {
		ID    string `json:"id"`
		Level uint32 `json:"level"`
		Email string `json:"email"`
	} `json:"userInfo"`
	UserNodes []struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
		Address  string `json:"address"`
		Port     int    `json:"port"`
		Security string `json:"security"`
		Flow     string `json:"flow"`
	} `json:"userNodes"`
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

func decodeSubscriptions(filePath string) Subscriptions {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var sub Subscriptions
	err = json.Unmarshal([]byte(content), &sub)
	if err != nil {
		panic(err)
	}
	return sub
}
