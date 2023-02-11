package api

import (
	"encoding/json"
	"fmt"

	"github.com/bionicosmos/submgr/models"
	//lint:ignore SA1019 Needed by the upstream
	"github.com/golang/protobuf/proto"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/proxy/trojan"
	"github.com/xtls/xray-core/proxy/vless"
)

func AddInbound(inbound *conf.InboundDetourConfig, apiAddress string) error {
	conn, ctx, cancelFunc, err := dialAPIServer(apiAddress)
	defer cancelFunc()
	if err != nil {
		return err
	}

	client := command.NewHandlerServiceClient(conn)
	inboundHandlerConfig, err := inbound.Build()
	if err != nil {
		return err
	}

	req := command.AddInboundRequest{
		Inbound: inboundHandlerConfig,
	}
	if _, err := client.AddInbound(ctx, &req); err != nil {
		return err
	}
	return conn.Close()
}

func RemoveInbound(tag string, apiAddress string) error {
	conn, ctx, cancelFunc, err := dialAPIServer(apiAddress)
	defer cancelFunc()
	if err != nil {
		return err
	}

	client := command.NewHandlerServiceClient(conn)
	req := command.RemoveInboundRequest{
		Tag: tag,
	}
	if _, err := client.RemoveInbound(ctx, &req); err != nil {
		return err
	}
	return conn.Close()
}

func AddUser(inbound *conf.InboundDetourConfig, user *models.User, apiAddress string) error {
	conn, ctx, cancelFunc, err := dialAPIServer(apiAddress)
	defer cancelFunc()
	if err != nil {
		return err
	}

	client := command.NewHandlerServiceClient(conn)
	var accountToTypedMessage proto.Message
	switch inbound.Protocol {
	case "vless":
		var account vless.Account
		if err := json.Unmarshal(user.Account["vless"], &account); err != nil {
			return err
		}
		accountToTypedMessage = &account
	case "vmess":
		var account conf.VMessAccount
		if err := json.Unmarshal(user.Account["vmess"], &account); err != nil {
			return err
		}
		accountToTypedMessage = account.Build()
	case "trojan":
		var account trojan.Account
		if err := json.Unmarshal(user.Account["trojan"], &account); err != nil {
			return err
		}
		accountToTypedMessage = &account
	default:
		return fmt.Errorf("nonsupported protocol: %v", inbound.Protocol)
	}
	req := command.AlterInboundRequest{
		Tag: inbound.Tag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level:   user.Level,
				Email:   user.Email,
				Account: serial.ToTypedMessage(accountToTypedMessage),
			},
		}),
	}
	if _, err := client.AlterInbound(ctx, &req); err != nil {
		return err
	}
	return conn.Close()
}

func RemoveUser(inbound *conf.InboundDetourConfig, user *models.User, apiAddress string) error {
	conn, ctx, cancelFunc, err := dialAPIServer(apiAddress)
	defer cancelFunc()
	if err != nil {
		return err
	}

	client := command.NewHandlerServiceClient(conn)
	req := command.AlterInboundRequest{
		Tag: inbound.Tag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: user.Email,
		}),
	}
	if _, err := client.AlterInbound(ctx, &req); err != nil {
		return err
	}
	return conn.Close()
}
