package redis

import (
	"time"
	"errors"

	"github.com/go-eyas/toolkit/log"
	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Cluster  bool
	Addrs    []string
	Password string
	DB       int
}

// redisClientInterface redis 实例拥有的功能
type redisClientInterface interface {
	redis.Cmdable
	Subscribe(...string) *redis.PubSub
	Close() error
}

// RedisClient redis client wrapper
type RedisClient struct {
	isCluster bool
	Namespace string
	Client    redisClientInterface
}

// RedisTTL 默认有效期 24 小时
var RedisTTL = time.Hour * 24

// Redis 暴露的redis封装
var Redis *RedisClient

// redis 客户端实例
var Client redisClientInterface

// Init 初始化redis
func Init(redisConf *RedisConfig) error {
	r := &RedisClient{}
	r.isCluster = redisConf.Cluster

	if len(redisConf.Addrs) == 0 {
		return errors.New("empty addrs")
	}

	if redisConf.Cluster {
		r.Client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisConf.Addrs,
			Password: redisConf.Password,
		})
	} else {
		r.Client = redis.NewClient(&redis.Options{
			Addr:     redisConf.Addrs[0],
			Password: redisConf.Password,
			DB:       redisConf.DB,
		})
	}
	_, err := r.Client.Ping().Result()
	if err != nil {
		log.Errorf("redis 连接失败, err=%v", err)
		return err
	}
	Redis = r
	Client = r.Client
	return nil
}

// Close 关闭redis连接
func Close() {
	if Redis != nil && Redis.Client != nil {
		Redis.Client.Close()
	}
}
