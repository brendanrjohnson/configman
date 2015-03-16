package appconf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/brendanrjohnson/configman/backends"
	"github.com/kelseyhightower/confd/log"
)

type Config struct {
	ConfDir      string
	ConfigDir    string
	Prefix       string
	StoreClient  backends.StoreClient
	BaseConfsDir string
}

type AppConfDefaultObject struct {
	AppConfDefault AppConfDefault `toml:"appconfdefault"`
}

type AppConfDefault struct {
	Conffile    string       `toml:"conffile"`
	SubSections []SubSection `toml:"subsection"`
	Prefix      string       `toml:"prefix"`
	Name        string       `toml:"name"`
}

type SubSection struct {
	Identifier  string       `toml:"identifier"`
	OptionPairs []OptionPair `toml:"pair"`
}

type OptionPair struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}

func NewAppConfDefault(path string) (*AppConfDefault, error) {
	var acdo *AppConfDefaultObject
	log.Debug("Loading template resource from " + path)
	_, err := toml.DecodeFile(path, &acdo)

	if err != nil {
		return nil, fmt.Errorf("Cannot process template resource % - %s", path, err.Error())
	}
	acd := acdo.AppConfDefault
	return &acd, nil

}
