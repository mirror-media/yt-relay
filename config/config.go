package config

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Address    string
	ApiKey     string `yaml:"apiKey"`
	Cache      Cache  `yaml:"cache"`
	Port       int
	Redis      *RedisService `yaml:"redis"`
	Whitelists Whitelists    `yaml:"whitelists"`
}

// Whitelist are maps, key is the whitelist string, value determines if it should be effective
type Whitelists struct {
	ChannelIDs  map[string]bool `yaml:"channelIDs"`
	PlaylistIDs map[string]bool `yaml:"playlistIDs"`
}

type Cache struct {
	IsEnabled    bool           `yaml:"isEnabled"`
	TTL          int            `yaml:"ttl"`
	OverwriteTTL []OverwriteTTL `yaml:"overwriteTtl"`
}

type OverwriteTTL struct {
	TTL       int    `yaml:"ttl"`
	PrefixAPI string `yaml:"apiPrefix"`
}

// RedisService defines the conf of redis for cache. User should find the right configuration according to the type
type RedisService struct {
	Type           RedisType              `yaml:"type"`
	Cluster        *RedisCluster          `yaml:"cluster"`
	SingleInstance *RedisSingleInstance   `yaml:"single"`
	Sentinel       *RedisSentinel         `yaml:"sentinel"`
	Replica        *RedisReplicaInstances `yaml:"replica"`
}

type RedisType string

const (
	Cluster  RedisType = "cluster"
	Single   RedisType = "single"
	Sentinel RedisType = "sentinel"
	Replica  RedisType = "replica"
)

type RedisCluster struct {
	Addrs    []RedisAddress `yaml:"addresses"`
	Password string         `yaml:"password"`
}

type RedisSingleInstance struct {
	Instance RedisAddress `yaml:"instance"`
	Password string       `yaml:"password"`
}

type RedisSentinel struct {
	Addrs    []RedisAddress `yaml:"addresses"`
	Password string         `yaml:"password"`
}

type RedisReplicaInstances struct {
	MasterAddrs []RedisAddress `yaml:"writers"`
	SlaveAddrs  []RedisAddress `yaml:"readers"`
	Password    string         `yaml:"password"`
}

type RedisAddress struct {
	Addr string `yaml:"address"`
	Port int    `yaml:"port"`
}

func (c *Conf) Valid() bool {

	if c.ApiKey == "" {
		return false
	}

	if len(c.Whitelists.ChannelIDs) == 0 {
		return false
	}

	if len(c.Whitelists.PlaylistIDs) == 0 {
		return false
	}

	if c.Redis != nil {
		redis := c.Redis
		switch redis.Type {
		case Cluster:
			if redis.Cluster == nil {
				return false
			}

			cluster := redis.Cluster

			if len(cluster.Addrs) == 0 {
				return false
			}
		case Single:
			if redis.SingleInstance == nil {
				return false
			}

			single := redis.SingleInstance

			if single.Instance.Addr == "" {
				return false
			}
		case Sentinel:
			if redis.Sentinel == nil {
				return false
			}

			sentinel := redis.Sentinel

			if len(sentinel.Addrs) == 0 {
				return false
			}
		case Replica:
			if redis.Replica == nil {
				return false
			}

			replica := redis.Replica

			if len(replica.MasterAddrs) == 0 {
				return false
			}

			if len(replica.SlaveAddrs) == 0 {
				return false
			}
		default:
			return false
		}
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
