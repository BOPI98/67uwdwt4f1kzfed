package main

import (
	"log"

	"gopkg.in/ini.v1"
)

type configType struct {
	ApiPort    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
	DbPort     string
}

var config configType

func initConfig() {
	f, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	err = f.MapTo(&config)
	if err != nil {
		log.Fatal(err)
	}
}
