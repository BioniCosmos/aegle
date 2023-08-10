package subscription

import (
	"encoding/json"

	proto "github.com/bionicosmos/submgr/services/subscription/protocol"
)

type Protocol interface {
	Id() (string, error)
	Host() (string, error)
	Encryption() string
	Flow() string
}

func NewProtocol(protocol string, settings json.RawMessage, account map[string]json.RawMessage) (Protocol, error) {
	var anyProtocol Protocol
	var err error
	switch protocol {
	case "vless":
		anyProtocol, err = proto.NewVless(settings, account["vless"])
	case "vmess":
		anyProtocol, err = proto.NewVmess(settings, account["vmess"])
	case "trojan":
		anyProtocol, err = proto.NewTrojan(settings, account["trojan"])
	default:
		return nil, UnknownProtocolError(protocol)
	}
	return anyProtocol, err
}
