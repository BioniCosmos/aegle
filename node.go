package main

import (
    "fmt"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    "github.com/xtls/xray-core/app/proxyman/command"
)

func obtainNode(nodeName string, config Configuration) (string, int) {
    for i := 0; i < len(config.Node); i++ {
        if nodeName == config.Node[i].Name {
            return config.Node[i].Address, config.Node[i].Port
        }
    }
    return "", 0
}

func dail(nodeAddress string, nodePort int, certificateFile string) command.HandlerServiceClient {
    creds, _ := credentials.NewClientTLSFromFile(certificateFile, nodeAddress)
    cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", nodeAddress, nodePort), grpc.WithTransportCredentials(creds))
    if err != nil {
        panic(err)
    }
    return command.NewHandlerServiceClient(cmdConn)
}
