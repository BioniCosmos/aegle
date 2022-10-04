package models

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/url"
    "strings"

    "github.com/xtls/xray-core/infra/conf"
    "github.com/xtls/xray-core/proxy/trojan"
    "github.com/xtls/xray-core/proxy/vless"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
    Id       primitive.ObjectID            `json:"id"`
    Name     string                        `json:"name"`
    Email    string                        `json:"email"`
    Level    uint32                        `json:"level"`
    Account  map[string]json.RawMessage    `json:"account"`
    Profiles map[primitive.ObjectID]string `json:"profiles"`
}

var usersColl *mongo.Collection

func FindUser(id string) (*User, error) {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var user User
    return &user, usersColl.FindOne(context.TODO(), bson.D{
        {Key: "_id", Value: _id},
    }).Decode(&user)
}

func FindUsers(filter any, sort any, skip int64, limit int64) ([]User, error) {
    cursor, err := usersColl.Find(context.TODO(), filter, options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit))
    if err != nil {
        return nil, err
    }

    var users []User
    return users, cursor.All(context.TODO(), &users)
}

func (user *User) Insert() error {
    _, err := usersColl.InsertOne(context.TODO(), user)
    return err
}

func (user *User) Update(id string) error {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = usersColl.UpdateByID(context.TODO(), _id, bson.D{
        {Key: "$set", Value: user},
    })
    return err
}

func DeleteUser(id string) error {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = usersColl.DeleteOne(context.TODO(), bson.D{
        {Key: "_id", Value: _id},
    })
    return err
}

func (user *User) GenerateSubscription(profile *Profile) (string, error) {
    var u url.URL
    var query url.Values
    outbound := profile.Outbound
    stream := outbound.StreamSetting

    if string(*stream.Network) != "" && string(*stream.Network) != "tcp" {
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
                return "", err
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
                return "", err
            }
            if headerType := header["type"]; headerType != "" {
                query.Set("headerType", headerType)
            }
            if quicSecurity := stream.QUICSettings.Security; quicSecurity != "" {
                query.Set("quicSecurity", quicSecurity)
                if key := stream.QUICSettings.Key; key != "" {
                    query.Set("key", key)
                } else {
                    return "", errors.New("key for QUIC not specified")
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
            return "", errors.New("unknown security type")
        }
    }

    proto := strings.ToLower(outbound.Protocol)
    u.Scheme = proto
    u.Fragment = outbound.Tag
    switch proto {
    case "vless":
        var settings conf.VLessOutboundConfig
        if err := json.Unmarshal(*outbound.Settings, &settings); err != nil {
            return "", err
        }

        vnext := settings.Vnext[0]
        u.Host = fmt.Sprintf("%v:%v", vnext.Address, vnext.Port)

        var account vless.Account
        if err := json.Unmarshal(user.Account["vless"], &account); err != nil {
            return "", err
        }
        u.User = url.User(account.Id)
        if encryption := account.Encryption; encryption != "" {
            query.Set("encryption", encryption)
        }
        if stream.Security == "xtls" {
            if flow := account.Flow; flow != "" {
                query.Set("flow", flow)
            } else {
                return "", errors.New("flow not specified")
            }
        }
    case "vmess":
        var settings conf.VMessOutboundConfig
        if err := json.Unmarshal(*outbound.Settings, &settings); err != nil {
            return "", err
        }

        vnext := settings.Receivers[0]
        u.Host = fmt.Sprintf("%v:%v", vnext.Address, vnext.Port)

        var account conf.VMessAccount
        if err := json.Unmarshal(user.Account["vmess"], &account); err != nil {
            return "", err
        }
        u.User = url.User(account.ID)
        if encryption := account.Security; encryption != "" && encryption != "auto" {
            query.Set("encryption", encryption)
        }
    case "trojan":
        var settings conf.TrojanClientConfig
        if err := json.Unmarshal(*outbound.Settings, &settings); err != nil {
            return "", err
        }

        server := settings.Servers[0]
        u.Host = fmt.Sprintf("%v:%v", server.Address, server.Port)

        var account trojan.Account
        if err := json.Unmarshal(user.Account["trojan"], &account); err != nil {
            return "", err
        }
        u.User = url.User(account.Password)
        if stream.Security == "xtls" {
            if flow := account.Flow; flow != "" {
                query.Set("flow", flow)
            } else {
                return "", errors.New("flow not specified")
            }
        }
    default:
        return "", errors.New("unknown protocol")
    }
    u.RawQuery = query.Encode()
    return u.String(), nil
}
