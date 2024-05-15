package subscription

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/bionicosmos/aegle/models"
	"github.com/bionicosmos/aegle/services/subscription/common"
	"github.com/xtls/xray-core/infra/conf"
)

func Generate(profile *models.Profile, user *models.User) (string, error) {
	outbound := conf.OutboundDetourConfig{}
	if err := json.Unmarshal([]byte(profile.Outbound), &outbound); err != nil {
		return "", err
	}
	// 1. basic information
	u := new(url.URL)
	protocolName := strings.ToLower(outbound.Protocol)
	if protocolName == "" {
		return "", &Error{ErrNoProtocol}
	}
	protocol, err := NewProtocol(protocolName, *outbound.Settings)
	if err != nil {
		return "", &Error{err}
	}
	// 1.1 protocol
	u.Scheme = protocolName
	// 1.2 uuid
	u.User = url.User(user.UUID)
	// 1.3 remote-host, 1.4 remote-port
	if host, err := protocol.Host(); err != nil {
		return "", &Error{err}
	} else {
		u.Host = host
	}
	// 1.5 descriptive-text
	u.Fragment = profile.Name

	// 2. protocols
	query := make(url.Values)
	stream := outbound.StreamSetting
	if stream == nil {
		stream = new(conf.StreamConfig)
	}
	network, err := buildNetwork(stream.Network)
	if err != nil {
		return "", &Error{err}
	}
	// 2.1 type
	query.Set("type", network)
	// 2.2 encryption
	if protocolName == "vmess" {
		query.Set("encryption", user.Security)
	}

	// 3. transport
	// 3.1 security
	query.Set("security", stream.Security)
	// 3.2, 3.4 path
	query.Set("path", getPath(network, stream))
	// 3.3, 3.5 host
	query.Set("host", getHost(network, stream))
	// 3.6, 3.10 headerType
	if headerType, err := getHeaderType(network, stream); err != nil {
		return "", &Error{err}
	} else {
		query.Set("headerType", headerType)
	}
	// 3.7 seed
	query.Set("seed", getSeed(stream.KCPSettings))
	// 3.8 quicSecurity
	query.Set("quicSecurity", getQuicSecurity(stream.QUICSettings))
	// 3.9 key
	if key, err := getKey(stream.QUICSettings); err != nil {
		return "", &Error{err}
	} else {
		query.Set("key", key)
	}
	// 3.11 serviceName
	query.Set("serviceName", getServiceName(stream.GRPCConfig))
	// 3.12 mode
	query.Set("mode", getMode(stream.GRPCConfig))

	// 4. TLS
	if stream.Security == "tls" || stream.Security == "reality" {
		// 4.0 fp
		if fp, err := getFp(stream); err != nil {
			return "", &Error{err}
		} else {
			query.Set("fp", fp)
		}
		// 4.1 sni
		if sni, err := getSni(stream); err != nil {
			return "", &Error{err}
		} else {
			query.Set("sni", sni)
		}
		// 4.2 alpn
		if alpn, err := getAlpn(stream); err != nil {
			return "", &Error{err}
		} else {
			query.Set("alpn", alpn)
		}
		// 4.4 flow
		query.Set("flow", user.Flow)
		if pbk, sid, spx, err := getReality(stream); err != nil {
			return "", &Error{err}
		} else {
			// 4.5 pbk
			query.Set("pbk", pbk)
			// 4.6 sid
			query.Set("sid", sid)
			// 4.7 spx
			query.Set("spx", spx)
		}
	}

	removeAllEmpty(query)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func buildNetwork(network *conf.TransportProtocol) (string, error) {
	if network == nil {
		return "", nil
	}
	net, err := network.Build()
	if err != nil {
		return "", err
	}
	switch net {
	case "mkcp":
		return "kcp", nil
	case "websocket":
		return "ws", nil
	default:
		return net, nil
	}
}

func getPath(network string, stream *conf.StreamConfig) string {
	switch network {
	case "http":
		if settings := stream.HTTPSettings; settings != nil {
			return settings.Path
		}
	case "ws":
		if settings := stream.WSSettings; settings != nil {
			return settings.Path
		}
	}
	return ""
}

func getHost(network string, stream *conf.StreamConfig) string {
	switch network {
	case "http":
		if settings := stream.HTTPSettings; settings != nil {
			if host := settings.Host; host != nil {
				return strings.Join(*host, ",")
			}
		}
	case "ws":
		if settings := stream.WSSettings; settings != nil {
			return settings.Headers["Host"]
		}
	}
	return ""
}

func getHeaderType(network string, stream *conf.StreamConfig) (string, error) {
	header := make(map[string]string)
	switch network {
	case "kcp":
		if settings := stream.KCPSettings; settings != nil {
			if err := json.Unmarshal(settings.HeaderConfig, &header); err != nil {
				return "", &ParseHeaderError{protocol: "mKCP", err: err}
			}
		}
	case "quic":
		if settings := stream.QUICSettings; settings != nil {
			if err := json.Unmarshal(settings.Header, &header); err != nil {
				return "", &ParseHeaderError{protocol: "QUIC", err: err}
			}
		}
	}
	return header["type"], nil
}

func getSeed(kcpSettings *conf.KCPConfig) string {
	if kcpSettings != nil {
		if seed := kcpSettings.Seed; seed != nil {
			return *seed
		}
	}
	return ""
}

func getQuicSecurity(quicSettings *conf.QUICConfig) string {
	if quicSettings == nil {
		return ""
	}
	return quicSettings.Security
}

func getKey(quicSettings *conf.QUICConfig) (string, error) {
	if quicSettings == nil {
		return "", nil
	}
	if common.IsNone(quicSettings.Security) && quicSettings.Key != "" {
		return "", ErrNoQuicSecurity
	}
	if !common.IsNone(quicSettings.Security) && quicSettings.Key == "" {
		return "", ErrNoQuicKey
	}
	return quicSettings.Key, nil
}

func getServiceName(grpcSettings *conf.GRPCConfig) string {
	if grpcSettings == nil {
		return ""
	}
	return grpcSettings.ServiceName
}

func getMode(grpcSettings *conf.GRPCConfig) string {
	if grpcSettings == nil {
		return ""
	}
	if grpcSettings.MultiMode {
		return "multi"
	}
	return ""
}

func getFp(stream *conf.StreamConfig) (string, error) {
	switch stream.Security {
	case "tls":
		if settings := stream.TLSSettings; settings != nil {
			return settings.Fingerprint, nil
		}
	case "reality":
		if settings := stream.REALITYSettings; settings == nil || settings.Fingerprint == "" {
			return "", ErrNoFingerprint
		} else {
			return settings.Fingerprint, nil
		}
	default:
		return "", common.UnknownSecurityError(stream.Security)
	}
	return "", nil
}

func getSni(stream *conf.StreamConfig) (string, error) {
	switch stream.Security {
	case "tls":
		if settings := stream.TLSSettings; settings != nil {
			return settings.ServerName, nil
		}
	case "reality":
		if settings := stream.REALITYSettings; settings != nil {
			return settings.ServerName, nil
		}
	default:
		return "", common.UnknownSecurityError(stream.Security)
	}
	return "", nil
}

func getAlpn(stream *conf.StreamConfig) (string, error) {
	if stream.Security != "tls" {
		return "", nil
	}
	if settings := stream.TLSSettings; settings != nil {
		if alpn := settings.ALPN; alpn != nil {
			return strings.Join(*alpn, ","), nil
		}
	}
	return "", nil
}

func getReality(stream *conf.StreamConfig) (string, string, string, error) {
	if stream.Security != "reality" {
		return "", "", "", nil
	}
	if settings := stream.REALITYSettings; settings == nil || settings.PublicKey == "" {
		return "", "", "", ErrNoPublicKey
	} else {
		return settings.PublicKey, settings.ShortId, settings.SpiderX, nil
	}
}

func removeAllEmpty(query url.Values) {
	for k, v := range query {
		if v[0] == "" {
			query.Del(k)
		}
	}
}
