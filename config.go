package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/brendanrjohnson/configman/backends"
	"github.com/kelseyhightower/confd/log"
)

var (
	configFile        = ""
	defaultConfigFile = "etc/configman/configman.toml"
	backend           string
	backendsConfig    backends.Config
	clientCaKeys      string
	clientCert        string
	clientKey         string
	confdir           string
	config            Config
	debug             bool
	defaultTTL        uint64
	nodes             Nodes
	onetime           bool
	prefix            string
	printVersion      bool
	quiet             bool
	scheme            string
	verbose           bool
)

// A config structure is used to configure configman
type Config struct {
	Backend      string   `toml:"backend"`
	BackendNodes []string `toml:"nodes"`
	ClientCakeys string   `toml:"client_cakeys"`
	ClientCert   string   `toml:"client_cert"`
	ClientKey    string   `toml:"client_key"`
	ConfDir      string   `toml:"confdir"`
	Debug        bool     `toml:"debug"`
	DefaultTTL   uint64   `toml:"defaultTTL"`
	Prefix       string   `toml:"prefix"`
	Quiet        bool     `toml:"quiet"`
	Scheme       string   `toml:"scheme"`
	Verbose      bool     `toml:"verbose"`
}

func init() {
	flag.StringVar(&backend, "backend", "etcd", "backend to use")
	flag.StringVar(&clientCaKeys, "client-ca-keys", "", "client ca keys")
	flag.StringVar(&clientCert, "client-cert", "", "the client cert")
	flag.StringVar(&clientKey, "client-key", "", "the client key")
	flag.StringVar(&confdir, "confdir", "/etc/configman", "configman conf directory")
	flag.StringVar(&configFile, "config-file", "", "the configman conf file")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Uint64Var(&defaultTTL, "defaultTTL", 300, "default etcd TTL")
	flag.Var(&nodes, "node", "list of backend nodes")
	flag.BoolVar(&onetime, "onetime", false, "run once and exit")
	flag.StringVar(&prefix, "prefix", "/configman", "key path prefix")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	flag.BoolVar(&quiet, "quiet", false, "enable quiet logging")
	flag.StringVar(&scheme, "scheme", "http", "the backend URI scheme (http or https)")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
}

// initConfig initializes the configman configuration by first setting defaults,
// then overriding settings from the configman config file, and finally overriding
// settings from the flags set on the command line.
// It returns and error if any.
func initConfig() error {
	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			configFile = defaultConfigFile
		}
	}
	// Set default
	config = Config{
		Backend:    "etcd",
		ConfDir:    "/etc/configman",
		DefaultTTL: 300,
		Prefix:     "/configman",
		Scheme:     "http",
	}
	// update config from the TOML configuration file
	if configFile == "" {
		log.Warning("Skipping configman config file")
	} else {
		log.Debug("Loading " + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			return err
		}
	}
	// Update config from commandline flags.
	processFlags()

	// Configure Logging.
	log.SetQuiet(config.Quiet)
	log.SetVerbose(config.Verbose)
	log.SetDebug(config.Debug)

	if len(config.BackendNodes) == 0 {
		switch config.Backend {
		case "consul":
			config.BackendNodes = []string{"127.0.0.1:8500"}
		case "etcd":
			peerstr := os.Getenv("ETCDCTL_PEERS")
			if len(peerstr) > 0 {
				config.BackendNodes = strings.Split(peerstr, ",")
			} else {
				config.BackendNodes = []string{"http://127.0.0.1:4001"}
			}
		}
	}
	// Initialize the storage client
	log.Notice("Backend set to " + config.Backend)
	backendsConfig = backends.Config{
		Backend:      config.Backend,
		ClientCaKeys: config.ClientCakeys,
		ClientCert:   config.ClientCert,
		ClientKey:    config.ClientKey,
		BackendNodes: config.BackendNodes,
		Scheme:       config.Scheme,
	}
	return nil
}

// processFlags iterates through each flag set on the command line and overrides corresponding configuration settings.
func processFlags() {
	flag.Visit(setConfigFromFlag)
}

func setConfigFromFlag(f *flag.Flag) {
	switch f.Name {
	case "backend":
		config.Backend = backend
	case "client-cakeys":
		config.ClientCakeys = clientCaKeys
	case "client-cert":
		config.ClientCert = clientCert
	case "client-key":
		config.ClientKey = clientKey
	case "confdir":
		config.ConfDir = confdir
	case "debug":
		config.Debug = debug
	case "defaultTTL":
		config.DefaultTTL = defaultTTL
	case "node":
		config.BackendNodes = nodes
	case "prefix":
		config.Prefix = prefix
	case "quiet":
		config.Quiet = quiet
	case "scheme":
		config.Scheme = scheme
	case "verbose":
		config.Verbose = verbose
	}
}
