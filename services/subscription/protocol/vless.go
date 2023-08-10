package protocol

import (
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
	_vless "github.com/xtls/xray-core/proxy/vless"
)

type vless struct {
	*conf.VLessOutboundConfig
	*_vless.Account
}

func (vless *vless) Id() (string, error) {
	id := vless.Account.Id
	if id == "" {
		return "", ErrNoId
	}
	return id, nil
}

func (vless *vless) Host() (string, error) {
	servers := vless.Vnext
	if len(servers) == 0 {
		return "", ErrNoHost
	}
	server := servers[0]
	if server == nil {
		return "", ErrNoHost
	}
	return getHost(server.Address, server.Port)
}

func (vless *vless) Encryption() string {
	encryption := vless.Account.Encryption
	if encryption == "none" {
		return ""
	}
	return encryption
}

func (vless *vless) Flow() string {
	return vless.Account.Flow
}

func NewVless(outboundSettings json.RawMessage, userAccount json.RawMessage) (*vless, error) {
	settings := new(conf.VLessOutboundConfig)
	if err := json.Unmarshal(outboundSettings, settings); err != nil {
		return nil, &ParseSettingsError{"VLESS", err}
	}
	account := new(_vless.Account)
	if err := json.Unmarshal(userAccount, account); err != nil {
		return nil, &ParseAccountError{"VLESS", err}
	}
	return &vless{VLessOutboundConfig: settings, Account: account}, nil
}
