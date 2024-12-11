package dbclient

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type IClientConnection interface {
	Connect() (err error)
	Ping() error
}

type ISqlClientConnection interface {
	IClientConnection
	GetDB() *bun.DB
}

type IRedisClientConnection interface {
	IClientConnection
	GetClient() *redis.Client
	SetTTL(key string, value any, t time.Duration) (string, error)
	Get(key string) (string, error)
	Set(key string, value any) (string, error)
	HGetAll(list string) (map[string]string, error)
	HSet(key string, values []any) (int64, error)
	HMSet(key string, values ...any) error
	Expire(key string, t time.Duration) error
	HGet(list, key string) (string, error)
	SAdd(key string, member ...any) (int64, error)
	SRem(key string, member ...any) (int64, error)
	IsMemberInSet(setkey, member string) (bool, error)
	Del(key []string) error
	IsExisted(key string) (bool, error)
}

// Config
type SqlConfig struct {
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	Timeout      int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
	MaxIdleConns int
	MaxOpenConns int
}

type RedisConfig struct {
	Addr     string
	Password string
	Database int
}

// DB container
type MySqlClientConn struct {
	*SqlConfig
	DB *bun.DB
}

type PostgresSqlClientConn struct {
	*SqlConfig
	DB *bun.DB
}

type RedisClientConn struct {
	*RedisConfig
	DB  *redis.Client
	Ctx context.Context
}

// Method
// MySQL
func (d *MySqlClientConn) Ping() error {
	if err := d.DB.Ping(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (d *MySqlClientConn) GetDB() *bun.DB {
	return d.DB
}

func (c *MySqlClientConn) Connect() (err error) {
	sqlConnectionInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?readTimeout=%ds&writeTimeout=%ds", c.Username, c.Password, c.Host, c.Port, c.Database, c.ReadTimeout, c.WriteTimeout)
	sqldb, err := sql.Open("mysql", sqlConnectionInfo)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sqldb.SetMaxIdleConns(c.MaxIdleConns)
	sqldb.SetMaxOpenConns(c.MaxOpenConns)
	db := bun.NewDB(sqldb, mysqldialect.New(), bun.WithDiscardUnknownColumns())
	c.DB = db
	return nil
}

// Postresql
func (d *PostgresSqlClientConn) Ping() error {
	if err := d.DB.Ping(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (d *PostgresSqlClientConn) GetDB() *bun.DB {
	return d.DB
}

func (c *PostgresSqlClientConn) Connect() (err error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", c.Host, c.Port)),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithUser(c.Username),
		pgdriver.WithPassword(c.Password),
		pgdriver.WithDatabase(c.Database),
		pgdriver.WithTimeout(time.Duration(c.Timeout)*time.Second),
		pgdriver.WithDialTimeout(time.Duration(c.DialTimeout)*time.Second),
		pgdriver.WithReadTimeout(time.Duration(c.ReadTimeout)*time.Second),
		pgdriver.WithWriteTimeout(time.Duration(c.WriteTimeout)*time.Second),
		pgdriver.WithInsecure(true),
	)
	sqldb := sql.OpenDB(pgconn)
	sqldb.SetMaxIdleConns(c.MaxIdleConns)
	sqldb.SetMaxOpenConns(c.MaxOpenConns)
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	c.DB = db
	return nil
}

// Redis
func (r *RedisClientConn) GetClient() *redis.Client {
	return r.DB
}

func (r *RedisClientConn) Connect() (err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.Database,
	})
	_, err = client.Ping(r.Ctx).Result()
	if err != nil {
		log.Fatal(err)
		return err
	}
	r.DB = client
	return nil
}

func (r *RedisClientConn) Ping() error {
	_, err := r.DB.Ping(r.Ctx).Result()
	return err
}

func (r *RedisClientConn) SetTTL(key string, value any, t time.Duration) (string, error) {
	ret, err := r.DB.Set(r.Ctx, key, value, t).Result()
	return ret, err
}

func (r *RedisClientConn) Get(key string) (string, error) {
	ret, err := r.DB.Get(r.Ctx, key).Result()
	return ret, err
}

func (r *RedisClientConn) Set(key string, value any) (string, error) {
	ret, err := r.DB.Set(r.Ctx, key, value, 0).Result()
	return ret, err
}

func (r *RedisClientConn) IsExisted(key string) (bool, error) {
	res, err := r.DB.Exists(r.Ctx, key).Result()
	if res == 0 || err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisClientConn) HGetAll(list string) (map[string]string, error) {
	ret, err := r.DB.HGetAll(r.Ctx, list).Result()
	return ret, err
}

func (r *RedisClientConn) HSet(key string, values []any) (int64, error) {
	ret, err := r.DB.HSet(r.Ctx, key, values...).Result()
	return ret, err
}

func (r *RedisClientConn) HMSet(key string, values ...any) error {
	ret, err := r.DB.HMSet(r.Ctx, key, values...).Result()
	if err != nil {
		return err
	}
	if !ret {
		err = errors.New("HashMap Set failed")
	}
	return err
}

func (r *RedisClientConn) Expire(key string, t time.Duration) error {
	ret, err := r.DB.Expire(r.Ctx, key, t).Result()
	if err != nil {
		return err
	}
	if !ret {
		err = errors.New("TTL Set failed")
	}
	return err
}

func (r *RedisClientConn) HGet(list, key string) (string, error) {
	ret, err := r.DB.HGet(r.Ctx, list, key).Result()
	return ret, err
}

func (r *RedisClientConn) SAdd(key string, member ...any) (int64, error) {
	ret, err := r.DB.SAdd(r.Ctx, key, member...).Result()
	return ret, err
}

func (r *RedisClientConn) SRem(key string, member ...any) (int64, error) {
	ret, err := r.DB.SRem(r.Ctx, key, member...).Result()
	return ret, err
}

func (r *RedisClientConn) IsMemberInSet(setkey, member string) (bool, error) {
	ret, err := r.DB.SIsMember(r.Ctx, setkey, member).Result()
	return ret, err
}

func (r *RedisClientConn) Del(key []string) error {
	err := r.DB.Del(r.Ctx, key...).Err()
	return err
}
