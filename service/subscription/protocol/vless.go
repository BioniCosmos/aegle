package protocol

import (
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
)

type vless struct {
	*conf.VLessOutboundConfig
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

func NewVless(outboundSettings json.RawMessage) (*vless, error) {
	settings := new(conf.VLessOutboundConfig)
	if err := json.Unmarshal(outboundSettings, settings); err != nil {
		return nil, &ParseSettingsError{"VLESS", err}
	}
	return &vless{VLessOutboundConfig: settings}, nil
}
