package config

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

const (
	defaultConfigFilePath = "res/configuration.toml"
)

var (
	configFilePath string
	conf           config
	ServerConf     Service
	DBConf         Database
	ScheduleConf   Schedule
	EdgexConf      Edgex
)

type config struct {
	Server   Service  `toml:"Service"`
	DB       Database `toml:"Database"`
	Edgex    Edgex    `toml:"Edgex"`
	Schedule Schedule `toml:"Schedule"`
}

type Service struct {
	Host                string
	Port                int64
	Labels              []string
	OpenMsg             string
	StaticResourcesPath string
}

type Scheme struct {
	Event          string
	Application    string
	Node           string
	TaskEvent      string
	Device         string
	ScheduleResult string
	EventToExecute string
	EventExecuted  string
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
type Scheduleelement struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
	Path string `json:"path"`
}
type Schedule struct {
	GetSchedule Scheduleelement
	AppSchedule Scheduleelement
}

type Edgex struct {
	Gateway string
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
	configFilePath = absPath
	if _, err := toml.DecodeFile(absPath, &conf); err != nil {
		log.Printf("Decode Config File Error:%v", err)
		return err
	}
	ServerConf = conf.Server
	DBConf = conf.DB
	ScheduleConf = conf.Schedule
	EdgexConf = conf.Edgex
	return nil
}

func ModifyAppSchedule(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var schedule Scheduleelement
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	ScheduleConf.AppSchedule = schedule
	conf.Schedule = ScheduleConf
	var newConfBuffer bytes.Buffer
	e := toml.NewEncoder(&newConfBuffer)
	if err := e.Encode(conf); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	newConf := make([]byte, newConfBuffer.Len())
	newConfBuffer.Read(newConf)
	if err := ioutil.WriteFile(configFilePath, newConf, 0666); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.Write([]byte("OK"))
}

func ListAppSchedule(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(&ScheduleConf.AppSchedule)
	w.Write(result)
}
