package protocol

import (
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
)

type trojan struct {
	*conf.TrojanClientConfig
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

func NewTrojan(outboundSettings json.RawMessage) (*trojan, error) {
	settings := new(conf.TrojanClientConfig)
	if err := json.Unmarshal(outboundSettings, settings); err != nil {
		return nil, &ParseSettingsError{"Trojan", err}
	}
	return &trojan{TrojanClientConfig: settings}, nil
}
