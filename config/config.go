package config

import (
	"os"
	"log"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Token				string
	ClientID		string
	Owner				string
}

var Cfg Config

// Reads info from config file
func ReadConfig(configfile string) {
	log.Println("Reading from config file:", configfile)
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing:", configfile)
	}

	if _, err := toml.DecodeFile(configfile, &Cfg); err != nil {
		log.Fatal(err)
	}
}