package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Listen             string `envconfig:"listen" default:""`
	Port               string `envconfig:"port" default:"2379"`
	ETCDVersion        string `envconfig:"etcd_version" default:"3.5.0"`
	ETCDClusterVersion string `envconfig:"etcd_cluster_version" default:"3.5.0"`
	Driver             string `envconfig:"driver" default:"badger"`
}

var conf config

func LoadConf() error {
	if err := envconfig.Process("", &conf); err != nil {
		return fmt.Errorf("config.LoadConf: failed to load conf: %w", err)
	}

	return nil
}

func Listen() string {
	return conf.Listen
}

func Port() string {
	return conf.Port
}

func ETCDVersion() string {
	return conf.ETCDVersion
}

func ETCDClusterVersion() string {
	return conf.ETCDClusterVersion
}

func Driver() string {
	return conf.Driver
}
