package Mongo


import (
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"Edgex-Export_Receiver/config"
	mgo "gopkg.in/mgo.v2"
)

var (
	database      string
	dbHost        string
	dbPort        int64
	dbUserName    string
	dbPassword    string
	eventScheme string
	taskScheme string
	nodeScheme string
)

type DataStore struct {
	S *mgo.Session
}

var DS DataStore

func (ds DataStore) DataStore() *DataStore {
	return &DataStore{ds.S.Copy()}
}

func loadConf() {
	database = config.DBConf.Name//edgex-ui-go
	dbHost = config.DBConf.Host//localhost
	dbPort = config.DBConf.Port//27017
	dbUserName = config.DBConf.Username//su
	dbPassword = config.DBConf.Password//su
	eventScheme = config.DBConf.Scheme.Event
	taskScheme = config.DBConf.Scheme.Task
	nodeScheme = config.DBConf.Scheme.Node
	log.Println(fmt.Sprintf("mongoDB connection info %s in %s:%d with credential (%s / %x), with scheme: %s %s, %s.",
		database, dbHost, dbPort, dbUserName, md5.Sum([]byte(dbPassword)), taskScheme, nodeScheme, eventScheme))
}

func DBConnect() bool {
	loadConf()

	mongoAddress := fmt.Sprintf("%s:%d", dbHost, dbPort)
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoAddress},
		Timeout:  time.Duration(5000) * time.Millisecond,
		Database: database,
		Username: dbUserName,
		Password: dbPassword,
	}
	s, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Println("Connect to mongoDB failed !")
		return false
	}
	s.SetSocketTimeout(time.Duration(5000) * time.Millisecond)
	DS.S = s
	log.Println("Success connect to mongoDB !")

	return true
}
