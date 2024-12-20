package database

import (
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type ISqlConnection interface {
	Connect() (err error)
	Close() error
	GetDB() *bun.DB
}

type IRedisConnection interface {
	Connect() (err error)
	Close() error
	GetDB() *redis.Client
}
