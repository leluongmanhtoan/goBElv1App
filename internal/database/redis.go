package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	Database int
}

type RedisDataSource struct {
	config *RedisConfig
	client *redis.Client
}

func NewRedisConnection(config *RedisConfig) IRedisConnection {
	return &RedisDataSource{
		config: config,
	}
}

func (d *RedisDataSource) Connect() (err error) {
	d.client = redis.NewClient(&redis.Options{
		Addr:     d.config.Addr,
		Password: d.config.Password,
		DB:       d.config.Database,
	})
	return d.client.Ping(context.Background()).Err()
}

func (d *RedisDataSource) Close() error {
	return d.client.Close()
}

func (d *RedisDataSource) GetDB() *redis.Client {
	return d.client
}
