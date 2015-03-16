package main

import (
	//	"encoding/json"
	"fmt"

	//	"github.com/brendanrjohnson/configman/backends"
	"github.com/kelseyhightower/confd/log"
)

func main() {

	configuration, err := NewAppConfDefault("etc/configman/conf.d/mariadb.toml")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(configuration)

	//machines := []string{"http://192.168.39.21:4001", "http://192.168.39.22:4001", "http://192.168.39.23:4001"}
	//cert := ""
	//key := ""
	//caCert := ""

}
