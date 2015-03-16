package appconf

import (
	"fmt"
	"path/filepath"
	//	"sync"
	//	"time"

	"github.com/kelseyhightower/confd/log"
)

func getBaseConfs(config Config) ([]*AppConfDefault, error) {
	var lastError error
	appconfs := make([]*AppConfDefault, 0)
	log.Debug("Loading template resources from confdir " + config.ConfDir)
	if !isFileExist(config.ConfDir) {
		log.Warning(fmt.Sprintf("Cannot load template resources confdir '%s' does not exist", config.ConfDir))
		return nil, nil
	}
	paths, err := filepath.Glob(filepath.Join(config.ConfigDir, "*.toml"))
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		a, err := NewAppConfDefault(p)
		if err != nil {
			lastError = err
			continue
		}
		appconfs = append(appconfs, a)
	}
	return appconfs, lastError
}
