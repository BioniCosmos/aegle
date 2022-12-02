package services

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/url"
    "strings"

    "github.com/bionicosmos/submgr/models"
    "github.com/xtls/xray-core/infra/conf"
    "github.com/xtls/xray-core/proxy/trojan"
    "github.com/xtls/xray-core/proxy/vless"
)

type subscriptionError struct {
    message error
}

func (err subscriptionError) Error() string {
    return fmt.Sprintf("failed to generate subscription: %v", err.message)
}

func (err subscriptionError) Unwrap() error {
    return err.message
}

func subscriptionErrorNew(err string) error {
    return subscriptionError{errors.New(err)}
}

func subscriptionErrorWrap(err error) error {
    return subscriptionError{err}
}

var SubscriptionError *subscriptionError

func GenerateSubscription(user *models.User, profile *models.Profile) (string, error) {
    var u url.URL
    var query url.Values
    outbound := profile.Outbound
    stream := outbound.StreamSetting
    if stream == nil {
        stream = new(conf.StreamConfig)
    }

    if stream.Network != nil && string(*stream.Network) != "" && string(*stream.Network) != "tcp" {
        network, err := stream.Network.Build()
        if err != nil {
            return "", nil
        }
        query.Set("type", network)
        switch network {
        case "http":
            query.Set("path", func() string {
                if path := stream.HTTPSettings.Path; path != "" {
                    return path
                }
                return "/"
            }())
            if host := stream.HTTPSettings.Host; host != nil {
                query.Set("host", strings.Join(*host, ","))
            }
        case "websocket":
            query.Set("path", func() string {
                if path := stream.WSSettings.Path; path != "" {
                    return path
                }
                return "/"
            }())
            if host := stream.WSSettings.Headers["Host"]; host != "" {
                query.Set("host", host)
            }
        case "mkcp":
            var header map[string]string
            if err := json.Unmarshal(stream.KCPSettings.HeaderConfig, &header); err != nil {
                return "", subscriptionErrorWrap(err)
            }
            if headerType := header["type"]; headerType != "" {
                query.Set("headerType", headerType)
            }
            if seed := stream.KCPSettings.Seed; seed != nil {
                query.Set("seed", *seed)
            }
        case "quic":
            var header map[string]string
            if err := json.Unmarshal(stream.QUICSettings.Header, &header); err != nil {
                return "", subscriptionErrorWrap(err)
            }
            if headerType := header["type"]; headerType != "" {
                query.Set("headerType", headerType)
            }
            if quicSecurity := stream.QUICSettings.Security; quicSecurity != "" {
                query.Set("quicSecurity", quicSecurity)
                if key := stream.QUICSettings.Key; key != "" {
                    query.Set("key", key)
                } else {
                    return "", subscriptionErrorNew("key for QUIC not specified")
                }
            }
        case "grpc":
            if serviceName := stream.GRPCConfig.ServiceName; serviceName != "" {
                query.Set("serviceName", serviceName)
            }
            if stream.GRPCConfig.MultiMode {
                query.Set("mode", "multi")
            }
        }
    }

    if security := stream.Security; security != "" && security != "tcp" {
        query.Set("security", security)
        switch security {
        case "tls":
            query.Set("sni", stream.TLSSettings.ServerName)
            query.Set("alpn", strings.Join(*stream.TLSSettings.ALPN, ","))
        case "xtls":
            query.Set("sni", stream.XTLSSettings.ServerName)
            query.Set("alpn", strings.Join(*stream.XTLSSettings.ALPN, ","))
        default:
            return "", subscriptionErrorNew("unknown security type")
        }
    }

    proto := strings.ToLower(outbound.Protocol)
    u.Scheme = proto
    u.Fragment = outbound.Tag
    switch proto {
    case "vless":
        var settings conf.VLessOutboundConfig
        if err := json.Unmarshal(*outbound.Settings, &settings); err != nil {
            return "", subscriptionErrorWrap(err)
        }

        vnext := settings.Vnext[0]
        u.Host = fmt.Sprintf("%v:%v", vnext.Address, vnext.Port)

        var account vless.Account
        if err := json.Unmarshal(user.Account["vless"], &account); err != nil {
            return "", subscriptionErrorWrap(err)
        }
        u.User = url.User(account.Id)
        if encryption := account.Encryption; encryption != "" {
            query.Set("encryption", encryption)
        }
        if stream.Security == "xtls" {
            if flow := account.Flow; flow != "" {
                query.Set("flow", flow)
            } else {
                return "", subscriptionErrorNew("flow not specified")
            }
        }
    case "vmess":
        var settings conf.VMessOutboundConfig
        if err := json.Unmarshal(*outbound.Settings, &settings); err != nil {
            return "", subscriptionErrorWrap(err)
        }

        vnext := settings.Receivers[0]
        u.Host = fmt.Sprintf("%v:%v", vnext.Address, vnext.Port)

        var account conf.VMessAccount
        if err := json.Unmarshal(user.Account["vmess"], &account); err != nil {
            return "", subscriptionErrorWrap(err)
        }
        u.User = url.User(account.ID)
        if encryption := account.Security; encryption != "" && encryption != "auto" {
            query.Set("encryption", encryption)
        }
    case "trojan":
        var settings conf.TrojanClientConfig
        if err := json.Unmarshal(*outbound.Settings, &settings); err != nil {
            return "", subscriptionErrorWrap(err)
        }

        server := settings.Servers[0]
        u.Host = fmt.Sprintf("%v:%v", server.Address, server.Port)

        var account trojan.Account
        if err := json.Unmarshal(user.Account["trojan"], &account); err != nil {
            return "", subscriptionErrorWrap(err)
        }
        u.User = url.User(account.Password)
        if stream.Security == "xtls" {
            if flow := account.Flow; flow != "" {
                query.Set("flow", flow)
            } else {
                return "", subscriptionErrorNew("flow not specified")
            }
        }
    default:
        return "", subscriptionErrorNew("unknown protocol")
    }
    u.RawQuery = query.Encode()
    return u.String(), nil
}
