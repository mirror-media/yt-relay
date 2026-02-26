package config

import (
	"errors"
	"io/ioutil"
	"regexp"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	// AppName is only allowed tt have alphanumeric, dash, and comma.
	AppName    string `yaml:"appName"`
	Address    string
	ApiKey     string `yaml:"apiKey"`
	Cache      Cache  `yaml:"cache"`
	Port       int
	Redis      *RedisService `yaml:"redis"`
	Whitelists Whitelists    `yaml:"whitelists"`
}

// Whitelists are maps, key is the whitelist string, value determines if it should be effective
type Whitelists struct {
	ChannelIDs  map[string]bool `yaml:"channelIDs"`
	PlaylistIDs map[string]bool `yaml:"playlistIDs"`
}

type Cache struct {
	IsEnabled    bool            `yaml:"isEnabled"`
	DisabledAPIs map[string]bool `yaml:"disabledApis"`
	TTL          int             `yaml:"ttl"`
	ErrorTTL     int             `yaml:"errorTtl"`
	OverwriteTTL map[string]int  `yaml:"overwriteTtl"`
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

	isValidAppName, _ := regexp.MatchString("^[a-zA-Z0-9.-]+$", c.AppName)
	if !isValidAppName {
		log.Errorf("appName(%s) can only contains alphanumeric, dot, and hyphen are allowed, and it cannot be empty", c.AppName)
		return false
	}

	if c.ApiKey == "" {
		log.Error("apiKey cannot be empty")
		return false
	}

	if len(c.Whitelists.ChannelIDs) == 0 {
		log.Error("whitelist's channel id cannot be empty")
		return false
	}

	if len(c.Whitelists.PlaylistIDs) == 0 {
		log.Error("whitelist's playlist id cannot be empty")
		return false
	}

	if c.Cache.IsEnabled {
		if c.Cache.TTL <= 0 {
			log.Errorf("enabled cache's default ttl(%d) cannot be zero or negative", c.Cache.TTL)
			return false
		}

		if c.Cache.ErrorTTL <= 0 {
			log.Errorf("enabled cache's default error ttl(%d) cannot be zero or negative", c.Cache.ErrorTTL)
			return false
		}

		for api, ttl := range c.Cache.OverwriteTTL {
			if ttl <= 0 {
				log.Errorf("enabled cache's ttl(%d) fot api(%s) cannot be zero or negative", ttl, api)
				return false
			}
		}
	}

	if c.Redis != nil {
		redis := c.Redis
		switch redis.Type {
		case Cluster:
			if redis.Cluster == nil {
				log.Error("redis type is set to %s but there is no %s configuration", Cluster, Cluster)
				return false
			}

			cluster := redis.Cluster

			if len(cluster.Addrs) == 0 {
				log.Errorf("%s addresses cannot be empty", Cluster)
				return false
			}
			for _, addr := range cluster.Addrs {
				if len(addr.Addr) == 0 {
					log.Errorf("one of the %s addresses is empty", Cluster)
					return false
				}
			}
		case Single:
			if redis.SingleInstance == nil {
				log.Error("redis type is set to %s but there is no %s configuration", Single, Single)
				return false
			}

			single := redis.SingleInstance

			if single.Instance.Addr == "" {
				log.Errorf("%s address cannot be empty", Single)
				return false
			}
		case Sentinel:
			if redis.Sentinel == nil {
				log.Error("redis type is set to %s but there is no %s configuration", Sentinel, Sentinel)
				return false
			}

			sentinel := redis.Sentinel

			if len(sentinel.Addrs) == 0 {
				log.Errorf("%s addresses cannot be empty", Sentinel)
				return false
			}
			for _, addr := range sentinel.Addrs {
				if len(addr.Addr) == 0 {
					log.Errorf("one of the %s addresses is empty", Sentinel)
					return false
				}
			}
		case Replica:
			if redis.Replica == nil {
				log.Error("redis type is set to %s but there is no %s configuration", Replica, Replica)
				return false
			}

			replica := redis.Replica

			if len(replica.MasterAddrs) == 0 {
				log.Errorf("%s writer addresses cannot be empty", Replica)
				return false
			}
			for _, addr := range replica.MasterAddrs {
				if len(addr.Addr) == 0 {
					log.Errorf("one of the %s writer addresses is empty", Replica)
					return false
				}
			}

			if len(replica.SlaveAddrs) == 0 {
				log.Errorf("%s reader addresses cannot be empty", Replica)
				return false
			}
			for _, addr := range replica.SlaveAddrs {
				if len(addr.Addr) == 0 {
					log.Errorf("one of the %s reader addresses is empty", Replica)
					return false
				}
			}
		default:
			log.Errorf("redis type(%s) is not supported", redis.Type)
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
