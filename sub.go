package main

import (
	"os"
	"strconv"
)

func generateSubLink(nodeName, nodeProto, nodeAddr, nodeId, nodeSec, nodeFlow string, nodePort int) string {
	var subLink string
	if nodeProto == "VLESS" {
		nodeProto = "vless://"
	}
	subLink += nodeProto + nodeId + "@" + nodeAddr + ":" + strconv.Itoa(nodePort) + "?security=" + nodeSec + "&flow=" + nodeFlow + "#" + nodeName
	return subLink
}

func createSubLinkFile(userId, subLink string) {
	f, err := os.Create(userId + ".html")
	check(err)
	defer f.Close()
	_, err = f.Write([]byte(subLink))
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
