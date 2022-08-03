package config

import (
	"io/ioutil"
	"os"

	"github.com/creasty/defaults"
	"github.com/tablesheep233/hook-service/logger"
	"gopkg.in/yaml.v3"
)

type DeployCmd struct {
	Port    int               `default:"9797"`
	Scripts map[string]string `yaml:"scripts"`
}

var Config DeployCmd

func init() {
	file, err := os.Open("config.yaml")
	if err != nil {
		logger.Error.Println("open config err", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Warning.Println("close config file err", err)
		}
	}()

	configBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Error.Println("read config err", err)
		return
	}

	if err := defaults.Set(&Config); err != nil {
		logger.Error.Println("Parse default config failed:", err)
		return
	}

	if err := yaml.Unmarshal(configBuffer, &Config); err != nil {
		logger.Error.Println("Parse config failed:", err)
		return
	}

	logger.Info.Println("load config succ")
}
