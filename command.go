package main

import "fmt"

func add(configPath, nodeName, email, id string, level int) {
	config := decodeConfiguration(configPath)
	nodeAddress, nodePort := obtainNode(nodeName, config)
	addVlessUser(nodeAddress, nodePort, config.CertificateFile, config.InboundTag, email, id, config.Flow, uint32(level))
}

func rm(configPath, nodeName, email string) {
	config := decodeConfiguration(configPath)
	nodeAddress, nodePort := obtainNode(nodeName, config)
	removeUser(nodeAddress, nodePort, config.CertificateFile, config.InboundTag, email)
}

func load(configPath, userInfoPath string) {
	config := decodeConfiguration(configPath)
	user := decodeUserInfo(userInfoPath)
	for k := 0; k < len(config.Node); k++ {
		for i := 0; i < len(user); i++ {
			fmt.Printf("Add %s to %s: ", user[i].Email, config.Node[k].Name)
			nodeAddress, nodePort := obtainNode(config.Node[k].Name, config)
			addVlessUser(nodeAddress, nodePort, config.CertificateFile, config.InboundTag, user[i].Email, user[i].Id, config.Flow, user[i].Level)
		}
	}
}

func empty(configPath, userInfoPath string) {
	config := decodeConfiguration(configPath)
	user := decodeUserInfo(userInfoPath)
	for k := 0; k < len(config.Node); k++ {
		for i := 0; i < len(user); i++ {
			fmt.Printf("Remove %s from %s: ", user[i].Email, config.Node[k].Name)
			nodeAddress, nodePort := obtainNode(config.Node[k].Name, config)
			removeUser(nodeAddress, nodePort, config.CertificateFile, config.InboundTag, user[i].Email)
		}
	}
}
