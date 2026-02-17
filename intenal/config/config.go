package config

import (
	"andrew_chat/intenal/domain"
	"encoding/json"
	"errors"
	"os"
)

var configPath string
var globalConfig *Config

type Config struct {
	Servers []domain.Server `json:"servers"`
}

func InitConfig(path string){
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err.Error())
    }

	configPath = path
    var cfg Config

    if err = json.Unmarshal(data, &cfg); err != nil {
        panic("Config corrupted, creating new one")
        
    }

    globalConfig = &cfg
}

func save() error {
	b, err := json.MarshalIndent(globalConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, b, 0644)
}

func AddServer(server domain.Server) error {
	if server.ID == "" {
		return errors.New("server id is empty")
	}

	for _, s := range globalConfig.Servers {
		if s.ID == server.ID {
			return errors.New("server already exists")
		}
	}

	globalConfig.Servers = append(globalConfig.Servers, server)

	return save()
}

func DeleteServer(id string) error {
	for i, s := range globalConfig.Servers {
		if s.ID == id {
			globalConfig.Servers =
				append(globalConfig.Servers[:i], globalConfig.Servers[i+1:]...)
			return save()
		}
	}

	return errors.New("server not found")
}

func UpdateServer(server domain.Server) error {
	for i, s := range globalConfig.Servers {
		if s.ID == server.ID {
			globalConfig.Servers[i] = server
			return save()
		}
	}

	return errors.New("server not found")
}

func GetServers() []domain.Server {
	servers := make([]domain.Server, len(globalConfig.Servers))
	copy(servers, globalConfig.Servers)
	return servers
}