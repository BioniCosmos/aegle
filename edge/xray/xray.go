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

func (*Server) AddInbound(
	_ context.Context,
	req *AddInboundRequest,
) (*Response, error) {
	store, err := storeRead("TODO")
	if err != nil {
		return nil, err
	}
	m := make(object)
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
	store, err := storeRead("TODO")
	if err != nil {
		return nil, err
	}
	store.inbounds = slices.DeleteFunc(
		store.inbounds,
		func(inbound object) bool {
			return inbound["tag"] == req.GetName()
		},
	)
	return store.apply()
}

func (*Server) AddUser(
	_ context.Context,
	req *AddUserRequest,
) (*Response, error) {
	store, err := storeRead("TODO")
	if err != nil {
		return nil, err
	}

	user := req.GetUser()
	email := user.GetEmail()
	level := user.GetLevel()
	uuid := user.GetUuid()
	flow := user.GetFlow()

	inbound := findInbound(&store, req.GetProfileName())
	proto := inbound["protocol"]
	var client object
	switch proto {
	case "vless":
		client = object{"id": uuid, "level": level, "email": email, "flow": flow}
	case "vmess":
		client = object{"id": uuid, "level": level, "email": email}
	case "trojan":
		client = object{"password": uuid, "level": level, "email": email}
	}
	settings := inbound["settings"].(object)
	settings["clients"] = append(settings["clients"].([]object), client)
	return store.apply()
}

func (*Server) RemoveUser(
	_ context.Context,
	req *RemoveUserRequest,
) (*Response, error) {
	store, err := storeRead("TODO")
	if err != nil {
		return nil, err
	}
	inbound := findInbound(&store, req.GetProfileName())
	settings := inbound["settings"].(object)
	settings["clients"] = slices.DeleteFunc(
		settings["clients"].([]object),
		func(client object) bool {
			return client["email"] == req.GetEmail()
		},
	)
	return store.apply()
}

type store struct {
	inbounds []object
	path     string
}
type object map[string]any

var mutex sync.Mutex

func storeRead(path string) (store, error) {
	mutex.Lock()
	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return store{}, err
		}
		data = []byte("{}")
	}
	m := make(object)
	if err := json.Unmarshal(data, &m); err != nil {
		return store{}, err
	}
	return store{inbounds: m["inbounds"].([]object), path: path}, nil
}

func (store *store) apply() (*Response, error) {
	defer mutex.Unlock()
	data, err := json.Marshal(object{"inbounds": store.inbounds})
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
	_, err = exec.Command("systemctl", "reload", "xray.service").Output()
	return nil, err
}

func findInbound(store *store, name string) object {
	i := slices.IndexFunc(store.inbounds, func(inbound object) bool {
		return inbound["tag"] == name
	})
	return store.inbounds[i]
}
