package main

import (
	"Edgex-Export_Receiver/app/Controller"
	"Edgex-Export_Receiver/app/config"
	"Edgex-Export_Receiver/app/db/Mongo"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := Controller.InitRestRoutes()

	var confFilePath string
	confFilePath = "res/configuration.toml"
	err := config.LoadConfig(confFilePath)
	if err != nil {
		log.Printf("Load config failed. Error:%v\n", err)
		return
	}

	ok := Mongo.DBConnect()
	if !ok {
		log.Println("Mongo Connection Failed")
	}

	server := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.FormatInt(config.ServerConf.Port, 10),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Edgex-Export-Receiver Listen At " + server.Addr)
	log.Fatal(server.ListenAndServe())
}


