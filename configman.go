package main

import (
	//	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/brendanrjohnson/configman/backends"
	"github.com/brendanrjohnson/configman/resources/baseconfs"
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
	fmt.Println(storeClient)
}
