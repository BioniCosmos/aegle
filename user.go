package main

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vless"
)

func addVlessUser(nodeAddress string, nodePort int, certificateFile string, inboundTag string, email string, id string, flow string, level uint32) {
	client := dail(nodeAddress, nodePort, certificateFile)
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: level,
				Email: email,
				Account: serial.ToTypedMessage(&vless.Account{
					Id:   id,
					Flow: flow,
				}),
			},
		}),
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK.")
	}
}

func removeUser(nodeAddress string, nodePort int, certificateFile string, inboundTag string, email string) {
	client := dail(nodeAddress, nodePort, certificateFile)
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: email,
		}),
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK.")
	}
}
