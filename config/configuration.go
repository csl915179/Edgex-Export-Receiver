package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
)



const (
	defaultConfigFilePath = "res/configuration.toml"
)

var (
	ServerConf   Service
	DBConf       Database
)

type config struct {
	Server Service      `toml:"Service"`
	DB     Database     `toml:"Database"`
}

type Service struct {
	Host                string
	Port                int64
	Labels              []string
	OpenMsg             string
	StaticResourcesPath string
}

type Scheme struct {
	Task string
	Node string
	Event string
}

type Database struct {
	Host     string
	Name     string
	Port     int64
	Username string
	Password string
	Timeout  int64
	Type     string
	Scheme   Scheme
}


func LoadConfig(confFilePath string) error {
	if len(confFilePath) == 0 {
		confFilePath = defaultConfigFilePath
	}

	absPath, err := filepath.Abs(confFilePath)
	if err != nil {
		log.Printf("Could not create absolute path to load configuration: %s; %v", absPath, err.Error())
		return err
	}
	log.Printf("Loading configuration from: %s\n", absPath)
	var conf config
	if _, err := toml.DecodeFile(absPath, &conf); err != nil {
		log.Printf("Decode Config File Error:%v", err)
		return err
	}
	ServerConf = conf.Server
	DBConf = conf.DB
	return nil
}
