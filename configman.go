package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/brendanrjohnson/configman/backends"
	"github.com/brendanrjohnson/configman/resources/appconf"
	"github.com/kelseyhightower/confd/log"
)

func main() {
	flag.Parse()
	if printVersion {
		fmt.Printf("configman %s\n", Version)
		os.Exit(0)
	}
	if err := initConfig(); err != nil {
		log.Fatal(err.Error())
	}
	log.Notice("Starting configman")
	storeClient, err := backends.New(backendsConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	storeClient.CreateDir(config.Prefix, config.DefaultTTL)
	configuration, err := appconf.NewAppConfDefault("etc/configman/conf.d/mariadb.toml")
	if err != nil {
		log.Fatal(err.Error())
	}
	storeClient.CreateDir((config.Prefix + "/" + configuration.Prefix), defaultTTL)

	jsonConfiguration, err := json.Marshal(configuration)
	if err != nil {
		fmt.Println(err)
	}

	storeClient.Set((config.Prefix + "/" + configuration.Prefix + "/" + configuration.Conffile), string(jsonConfiguration), defaultTTL)
}
