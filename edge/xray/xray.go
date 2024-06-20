package xray

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"slices"
	"sync"
)

type Server struct {
	UnimplementedXrayServer
}

var ErrNoInbound = errors.New("inbound not found")

func (*Server) AddInbound(
	_ context.Context,
	req *AddInboundRequest,
) (*Response, error) {
	store, err := storeRead()
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	if err := json.Unmarshal([]byte(req.GetInbound()), &m); err != nil {
		return nil, err
	}
	m["tag"] = req.GetName()
	store.inbounds = append(store.inbounds, m)
	return store.apply()
}

func (*Server) RemoveInbound(
	_ context.Context,
	req *RemoveInboundRequest,
) (*Response, error) {
	store, err := storeRead()
	if err != nil {
		return nil, err
	}
	store.inbounds = slices.DeleteFunc(
		store.inbounds,
		func(inbound map[string]any) bool {
			return inbound["tag"] == req.GetName()
		},
	)
	return store.apply()
}

func (*Server) AddUser(
	_ context.Context,
	req *AddUserRequest,
) (*Response, error) {
	store, err := storeRead()
	if err != nil {
		return nil, err
	}

	user := req.GetUser()
	email := user.GetEmail()
	level := user.GetLevel()
	uuid := user.GetUuid()
	flow := user.GetFlow()

	inbound, err := findInbound(&store, req.GetProfileName())
	if err != nil {
		return nil, err
	}
	proto := inbound["protocol"]
	var client map[string]any
	switch proto {
	case "vless":
		client = map[string]any{"id": uuid, "level": level, "email": email, "flow": flow}
	case "vmess":
		client = map[string]any{"id": uuid, "level": level, "email": email}
	case "trojan":
		client = map[string]any{"password": uuid, "level": level, "email": email}
	}
	fillNil(inbound)
	settings := inbound["settings"].(map[string]any)
	settings["clients"] = append(settings["clients"].([]any), client)
	return store.apply()
}

func (*Server) RemoveUser(
	_ context.Context,
	req *RemoveUserRequest,
) (*Response, error) {
	store, err := storeRead()
	if err != nil {
		return nil, err
	}
	inbound, err := findInbound(&store, req.GetProfileName())
	if err != nil {
		return nil, err
	}
	fillNil(inbound)
	settings := inbound["settings"].(map[string]any)
	settings["clients"] = slices.DeleteFunc(
		settings["clients"].([]any),
		func(client any) bool {
			return client.(map[string]any)["email"] == req.GetEmail()
		},
	)
	return store.apply()
}

type store struct {
	inbounds []map[string]any
	path     string
}

var mutex sync.Mutex

func storeRead() (store, error) {
	path := os.Getenv("XRAY_CONFIG")
	mutex.Lock()
	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return store{}, err
		}
		data = []byte(`{"inbounds":[]}`)
	}
	m := struct{ Inbounds []map[string]any }{}
	if err := json.Unmarshal(data, &m); err != nil {
		return store{}, err
	}
	return store{inbounds: m.Inbounds, path: path}, nil
}

func (store *store) apply() (*Response, error) {
	defer mutex.Unlock()
	data, err := json.Marshal(map[string]any{"inbounds": store.inbounds})
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(store.path, data, 0644); err != nil {
		return nil, err
	}
	message, err := exec.
		Command("xray", "run", "-confdir", store.path, "-test").
		Output()
	if err != nil {
		return nil, errors.New(string(message))
	}
	_, err = exec.Command("systemctl", "restart", "xray.service").Output()
	return nil, err
}

func findInbound(store *store, name string) (map[string]any, error) {
	i := slices.IndexFunc(store.inbounds, func(inbound map[string]any) bool {
		return inbound["tag"] == name
	})
	if i == -1 {
		return nil, ErrNoInbound
	}
	return store.inbounds[i], nil
}

func fillNil(inbound map[string]any) {
	if inbound["settings"] == nil {
		inbound["settings"] = make(map[string]any)
	}
	if settings := inbound["settings"].(map[string]any); settings["clients"] == nil {
		settings["clients"] = make([]any, 0)
	}
}
