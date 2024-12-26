package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

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

type MySqlDataSource struct {
	config *SqlConfig
	db     *bun.DB
}

func NewMySqlConnection(config *SqlConfig) ISqlConnection {
	return &MySqlDataSource{
		config: config,
	}
}

func (d *MySqlDataSource) Connect() (err error) {
	sqlConnectionInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?readTimeout=%ds&writeTimeout=%ds", d.config.Username, d.config.Password, d.config.Host, d.config.Port, d.config.Database, d.config.ReadTimeout, d.config.WriteTimeout)
	sqldb, err := sql.Open("mysql", sqlConnectionInfo)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sqldb.SetMaxIdleConns(d.config.MaxIdleConns)
	sqldb.SetMaxOpenConns(d.config.MaxOpenConns)
	db := bun.NewDB(sqldb, mysqldialect.New(), bun.WithDiscardUnknownColumns())
	d.db = db
	return d.db.Ping()
}

func (d *MySqlDataSource) Close() error {
	return d.db.Close()
}

func (d *MySqlDataSource) GetDB() *bun.DB {
	return d.db
}
