package etcd

import (
	"errors"
	"strings"
	"time"

	goetcd "github.com/coreos/go-etcd/etcd"
)

type Client struct {
	client *goetcd.Client
}

// NewEtcdClient returns an *etcd.Client with a connection to named machines.
// It returns an error if a connection to the cluster cannot be made.
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

// GetValues queries etcd for keys prefixed by prefix.
func (c *Client) GetValues(keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, key := range keys {
		resp, err := c.client.Get(key, true, true)
		if err != nil {
			return vars, err
		}
		err = nodeWalk(resp.Node, vars)
		if err != nil {
			return vars, err
		}
	}
	return vars, nil
}

func (c *Client) CreateDir(key string, ttl uint64) error {
	_, err := c.client.CreateDir(key, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Set(key string, value string, ttl uint64) error {
	_, err := c.client.Set(key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

// nodeWalk recursively descends nodes, updating vars.
func nodeWalk(node *goetcd.Node, vars map[string]string) error {
	if node != nil {
		key := node.Key
		if !node.Dir {
			vars[key] = node.Value
		} else {
			for _, node := range node.Nodes {
				nodeWalk(node, vars)
			}
		}
	}
	return nil
}
