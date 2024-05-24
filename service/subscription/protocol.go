package subscription

import (
	"encoding/json"

	proto "github.com/bionicosmos/aegle/service/subscription/protocol"
)

type Protocol interface {
	Host() (string, error)
}

func NewProtocol(protocol string, settings json.RawMessage) (Protocol, error) {
	var anyProtocol Protocol
	var err error
	switch protocol {
	case "vless":
		anyProtocol, err = proto.NewVless(settings)
	case "vmess":
		anyProtocol, err = proto.NewVmess(settings)
	case "trojan":
		anyProtocol, err = proto.NewTrojan(settings)
	default:
		return nil, UnknownProtocolError(protocol)
	}
	return anyProtocol, err
}
