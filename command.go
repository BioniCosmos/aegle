package main

import (
	"encoding/base64"
	"fmt"
)

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
	sub := decodeSubscriptions(userInfoPath)
	for k := 0; k < len(sub); k++ {
		for i := 0; i < len(sub[k].UserNodes); i++ {
			subUserNode := sub[k].UserNodes[i]
			fmt.Printf("Add %s to %s: ", sub[k].UserInfo.Email, subUserNode.Name)
			nodeAddress, nodePort := obtainNode(subUserNode.Name, config)
			addVlessUser(nodeAddress, nodePort, config.CertificateFile, config.InboundTag, sub[k].UserInfo.Email, sub[k].UserInfo.ID, config.Flow, sub[k].UserInfo.Level)
		}
	}
}

func empty(configPath, userInfoPath string) {
	config := decodeConfiguration(configPath)
	sub := decodeSubscriptions(userInfoPath)
	for k := 0; k < len(sub); k++ {
		for i := 0; i < len(sub[k].UserNodes); i++ {
			subUserNode := sub[k].UserNodes[i]
			fmt.Printf("Remove %s from %s: ", sub[k].UserInfo.Email, subUserNode.Name)
			nodeAddress, nodePort := obtainNode(subUserNode.Name, config)
			removeUser(nodeAddress, nodePort, config.CertificateFile, config.InboundTag, sub[k].UserInfo.Email)
		}
	}
}

func sub(userInfoPath string) {
	sub := decodeSubscriptions(userInfoPath)
	var subLink string
	for i := 0; i < len(sub); i++ {
		for j := 0; j < len(sub[i].UserNodes); j++ {
			subUserNode := sub[i].UserNodes[j]
			subLink += generateSubLink(subUserNode.Name, subUserNode.Protocol, subUserNode.Address, sub[i].UserInfo.ID, subUserNode.Security, subUserNode.Flow, subUserNode.Port)
			subLink += "\n"
		}
		encodedSubLink := base64.StdEncoding.EncodeToString([]byte(subLink))
		createSubLinkFile(sub[i].UserInfo.ID[:8], encodedSubLink)
		subLink = ""
	}
}

func server(userInfoPath, serverListen, serverPath string) {
	sub := decodeSubscriptions(userInfoPath)
	subLinks := make(map[string]string)
	for i := range sub {
		var subLink string
		for j := range sub[i].UserNodes {
			subUserNode := sub[i].UserNodes[j]
			subLink += generateSubLink(subUserNode.Name, subUserNode.Protocol, subUserNode.Address, sub[i].UserInfo.ID, subUserNode.Security, subUserNode.Flow, subUserNode.Port)
			subLink += "\n"
		}
		encodedSubLink := base64.StdEncoding.EncodeToString([]byte(subLink))
		subLinks[sub[i].UserInfo.ID[:8]] = encodedSubLink
	}
	serveSubscriptions(subLinks, serverListen, serverPath)
}
