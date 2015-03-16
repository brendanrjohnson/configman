package backends

import (
	"errors"
	"strings"

	"github.com/brendanrjohnson/configman/backends/etcd"
	"github.com/kelseyhightower/confd/log"
)

// The StoreClient interface is implemented by objects that can retrieve key/value pairs from a backend store.
type StoreClient interface {
	GetValues(keys []string) (map[string]string, error)
	CreateDir(key string, ttl uint64) error
	Set(key string, value string, ttl uint64) error
}

// New is used to create a storage client based on our configuration.
func New(config Config) (StoreClient, error) {
	if config.Backend == "" {
		config.Backend = "etcd"
	}
	backendNodes := config.BackendNodes
	log.Notice("Backend nodes set to " + strings.Join(backendNodes, ", "))
	switch config.Backend {
	case "etcd":
		// Create the etcd client upfront and usi it for the life of the process.
		return etcd.NewEtcdClient(backendNodes, config.ClientCert, config.ClientKey, config.ClientCaKeys)
	}
	return nil, errors.New("Invalid backend")
}
