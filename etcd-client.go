package main

import (
	"errors"
	"strings"
	"time"

	goetcd "github.com/coreos/go-etcd/etcd"
)

type Client struct {
	client *goetcd.Client
}

func NewEtcdClient(machines []string, cert, key, caCert string) (*Client, error) {
	var c *goetcd.Client
	var err error
	if cert != "" && key != "" {
		c, err = goetcd.NewTLSClient(machines, cert, key, caCert)
		if err != nil {
			return &Client{c}, err
		}
	} else {
		c = goetcd.NewClient(machines)
	}

	// Configure the DialTimeOut, since 1 second is often to short
	c.SetDialTimeout(time.Duration(3) * time.Second)
	success := c.SetCluster(machines)
	if !success {
		return &Client{c}, errors.New("cannot connect to etcd cluster: " + strings.Join(machines, ","))
	}
	return &Client{c}, nil
}
