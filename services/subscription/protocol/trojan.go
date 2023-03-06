package protocol

import (
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
	_trojan "github.com/xtls/xray-core/proxy/trojan"
)

type trojan struct {
	*conf.TrojanClientConfig
	*_trojan.Account
}

func (trojan *trojan) Id() (string, error) {
	password := trojan.Password
	if password == "" {
		return "", ErrNoId
	}
	return password, nil
}

func (trojan *trojan) Host() (string, error) {
	servers := trojan.Servers
	if len(servers) == 0 {
		return "", ErrNoHost
	}
	server := servers[0]
	if server == nil {
		return "", ErrNoHost
	}
	return getHost(server.Address, server.Port)
}

func (*trojan) Encryption() string {
	return ""
}

func (trojan *trojan) Flow(security string) (string, error) {
	if security != "xtls" {
		return "", nil
	}
	flow := trojan.Account.Flow
	if flow == "" {
		return "", FlowError("XTLS")
	}
	return flow, nil
}

func NewTrojan(outboundSettings json.RawMessage, userAccount json.RawMessage) (*trojan, error) {
	settings := new(conf.TrojanClientConfig)
	if err := json.Unmarshal(outboundSettings, settings); err != nil {
		return nil, &ParseSettingsError{"Trojan", err}
	}
	account := new(_trojan.Account)
	if err := json.Unmarshal(userAccount, account); err != nil {
		return nil, &ParseAccountError{"Trojan", err}
	}
	return &trojan{TrojanClientConfig: settings, Account: account}, nil
}
