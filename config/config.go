package config

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	ApiKey      string `yaml:"apiKey"`
	Address     string
	ChannelID   string   `yaml:"channelID"`
	PlaylistIDs []string `yaml:"playlistIDs"`
	Port        int
}

func (c *Conf) Valid() bool {

	if c.ApiKey == "" {
		return false
	}

	if c.ChannelID == "" {
		return false
	}
	return true
}

// LoadFile attempts to load the configuration file stored at the path
// and returns the configuration. On error, it returns nil.
func LoadFile(path string) (*Conf, error) {
	log.Printf("loading configuration file from %s", path)
	if path == "" {
		return nil, errors.New("invalid path")
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("could not read configuration file")
	}

	return LoadConfig(body)
}

// LoadConfig attempts to load the configuration from a byte slice.
// On error, it returns nil.
func LoadConfig(config []byte) (*Conf, error) {
	var cfg = &Conf{}
	err := yaml.Unmarshal(config, &cfg)
	if err != nil {
		return nil, errors.New("failed to unmarshal configuration: " + err.Error())
	}

	if !cfg.Valid() {
		return nil, errors.New("invalid configuration")
	}

	log.Println("configuration ok")
	return cfg, nil
}
