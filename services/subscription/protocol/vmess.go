package protocol

import (
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
)

type vmess struct {
	*conf.VMessOutboundConfig
	*conf.VMessAccount
}

func (vmess *vmess) Id() (string, error) {
	id := vmess.ID
	if id == "" {
		return "", ErrNoId
	}
	return id, nil
}

func (vmess *vmess) Host() (string, error) {
	servers := vmess.Receivers
	if len(servers) == 0 {
		return "", ErrNoHost
	}
	server := servers[0]
	if server == nil {
		return "", ErrNoHost
	}
	return getHost(server.Address, server.Port)
}

func (vmess *vmess) Encryption() string {
	return vmess.Security
}

func (*vmess) Flow(string) (string, error) {
	return "", nil
}

func NewVmess(outboundSettings json.RawMessage, userAccount json.RawMessage) (*vmess, error) {
	settings := new(conf.VMessOutboundConfig)
	if err := json.Unmarshal(outboundSettings, settings); err != nil {
		return nil, &ParseSettingsError{"VMess", err}
	}
	account := new(conf.VMessAccount)
	if err := json.Unmarshal(userAccount, account); err != nil {
		return nil, &ParseAccountError{"VMess", err}
	}
	return &vmess{VMessOutboundConfig: settings, VMessAccount: account}, nil
}
