package protocol

import (
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
)

type vmess struct {
	*conf.VMessOutboundConfig
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

func NewVmess(outboundSettings json.RawMessage) (*vmess, error) {
	settings := new(conf.VMessOutboundConfig)
	if err := json.Unmarshal(outboundSettings, settings); err != nil {
		return nil, &ParseSettingsError{"VMess", err}
	}
	return &vmess{VMessOutboundConfig: settings}, nil
}
