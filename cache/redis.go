package cache

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mirror-media/yt-relay/config"
)

// replicaTypeRedis implements Rediser
type replicaTypeRedis struct {
	writeCount uint32
	readCount  uint32
	writers    []*redis.Client
	readers    []*redis.Client
}

func (r *replicaTypeRedis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
	wc := atomic.AddUint32(&r.writeCount, 1)
	i := int(wc) % len(r.writers)
	return r.writers[i].Set(ctx, key, value, ttl)
}

func (r *replicaTypeRedis) SetXX(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.BoolCmd {
	wc := atomic.AddUint32(&r.writeCount, 1)
	i := int(wc) % len(r.writers)
	return r.writers[i].SetXX(ctx, key, value, ttl)
}

func (r *replicaTypeRedis) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.BoolCmd {
	wc := atomic.AddUint32(&r.writeCount, 1)
	i := int(wc) % len(r.writers)
	return r.writers[i].SetNX(ctx, key, value, ttl)
}

func (r *replicaTypeRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	rc := atomic.AddUint32(&r.readCount, 1)
	i := int(rc) % len(r.writers)
	return r.readers[i].Get(ctx, key)
}
func (r *replicaTypeRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	wc := atomic.AddUint32(&r.writeCount, 1)
	i := int(wc) % len(r.writers)
	return r.writers[i].Del(ctx, keys...)
}

func NewReplicaRedisService(MasterAddrs []config.RedisAddress, SlaveAddrs []config.RedisAddress, Password string) (Rediser, error) {
	instance := replicaTypeRedis{}
	writers := make([]*redis.Client, 0, len(MasterAddrs))
	for _, a := range MasterAddrs {
		writers = append(writers, redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", a.Addr, a.Port),
			Password:     Password,
			PoolSize:     20,
			MaxRetries:   0,
			DialTimeout:  time.Second,
			IdleTimeout:  10 * time.Second,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}))
	}
	instance.writers = writers
	readers := make([]*redis.Client, 0, len(SlaveAddrs))
	for _, a := range SlaveAddrs {
		readers = append(readers, redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", a.Addr, a.Port),
			Password:     Password,
			PoolSize:     20,
			MaxRetries:   0,
			DialTimeout:  time.Second,
			IdleTimeout:  10 * time.Second,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}))
	}
	instance.readers = readers
	return &instance, nil
}
