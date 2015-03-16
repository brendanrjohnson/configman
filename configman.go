package main

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/confd/log"
)

func main() {

	configuration, err := NewAppConfDefault("etc/configman/conf.d/mariadb.toml")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(configuration)

	machines := []string{"http://192.168.39.21:4001", "http://192.168.39.22:4001", "http://192.168.39.23:4001"}
	cert := ""
	key := ""
	caCert := ""

	etcdClient, err := NewEtcdClient(machines, cert, key, caCert)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ttl uint64

	ttl = 300
	rootkey := "configman"
	_, err = etcdClient.client.CreateDir(rootkey, ttl)
	if err != nil {
		fmt.Println(err)
	}

	jsonConfiguration, err := json.Marshal(configuration)
	if err != nil {
		fmt.Println(err)
	}

	_, err = etcdClient.client.Set("/configman/mariadb", string(jsonConfiguration), ttl)
	if err != nil {
		fmt.Println(err)
	}

}
