// Package cache defines and implement the cache layer, especially for redis
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
	"github.com/mirror-media/yt-relay/config"
)

type HTTP struct {
	StatusCode int    `json:"code"`
	Response   []byte `json:"response"`
}

type Rediser interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd
	SetXX(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.BoolCmd
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.BoolCmd

	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

func GetCacheKey(namespace string, name string) (string, error) {
	if namespace == "" {
		err := errors.New("namespace cannot be empty")
		return "", err
	}

	if name == "" {
		err := errors.New("key cannot be empty")
		return "", err
	}

	return fmt.Sprintf("%s:cache:%s", namespace, name), nil
}

func NewRedis(c config.Conf) (rdb Rediser, err error) {
	switch c.Redis.Type {
	case config.Cluster:
		cluster := c.Redis.Cluster
		if len(cluster.Addrs) == 0 {
			return nil, errors.New("there's no cluster redis address provided")
		}
		addrs := make([]string, 0, len(cluster.Addrs))
		for _, a := range cluster.Addrs {
			addrs = append(addrs, fmt.Sprintf("%s:%d", a.Addr, a.Port))
		}
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: cluster.Password,
		})
	case config.Single:
		single := c.Redis.SingleInstance
		if single.Instance.Addr == "" {
			return nil, errors.New("there's no single instance redis address provided")
		}

		addr := fmt.Sprintf("%s:%d", single.Instance.Addr, single.Instance.Port)

		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: single.Password,
		})
	case config.Sentinel:
		sentinel := c.Redis.Sentinel
		if len(sentinel.Addrs) == 0 {
			return nil, errors.New("there's no Sentinel redis address provided")
		}

		addrs := make([]string, 0, len(sentinel.Addrs))
		for _, a := range sentinel.Addrs {
			addrs = append(addrs, fmt.Sprintf("%s:%d", a.Addr, a.Port))
		}
		rdb = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: addrs,
			Password:      sentinel.Password,
		})
	case config.Replica:
		replica := c.Redis.Replica
		if len(replica.MasterAddrs) == 0 {
			return nil, errors.New("there's no master redis address provided")
		}
		if len(replica.SlaveAddrs) == 0 {
			return nil, errors.New("there's no slave redis address provided")
		}
		if rdb, err = NewReplicaRedisService(replica.MasterAddrs, replica.SlaveAddrs, replica.Password); err != nil {
			err = errors.Wrap(err, "Cannot create Replica type Redis service")
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported redis type(%s)", c.Redis.Type)
	}
	return rdb, nil
}
