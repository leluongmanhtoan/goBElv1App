package main

import (
	"context"
	"log"
	"os"
	server "program/api"
	apiv1 "program/api/v1"
	"program/internal/dbclient"
	"program/middleware"
	"program/repository"
	"program/repository/db"
	"program/service"

	"github.com/joho/godotenv"
)

func init() {
	// Read .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalln("env file load failed")
	}

	// Create sql connection

	sqlClientConfig := dbclient.SqlConfig{
		Host:         "localhost",
		Port:         3306,
		Database:     "testdb",
		Username:     "root",
		Password:     "ManhToan0123!",
		ReadTimeout:  30,
		WriteTimeout: 30,
		MaxIdleConns: 10,
		MaxOpenConns: 10,
	}

	//Select sql connect to Mysql Database with config
	repository.SqlClientConnection = &dbclient.MySqlClientConn{
		SqlConfig: &sqlClientConfig,
	}

	//Run connect to Mysql database
	err := repository.SqlClientConnection.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = repository.SqlClientConnection.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create redis connection
	redisClientConfig := dbclient.RedisConfig{
		Addr:     "localhost:6379",
		Password: "ManhToan0123",
		Database: 0,
	}
	//Select redis connect to Redis database with config
	repository.RedisClientConnection = &dbclient.RedisClientConn{
		RedisConfig: &redisClientConfig,
		Ctx:         context.Background(),
	}

	//Run connect to Redis database
	repository.RedisClientConnection.Connect()
	if err != nil {
		log.Fatal(err)
	}
	//Ping Redis connect
	err = repository.RedisClientConnection.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//Init repository
	repository.UserRepo = db.NewUserRepo()
	repository.RelationshipsRepo = db.NewRelationshipsRepo()
	repository.NewsfeedRepo = db.NewNewsFeedRepo()

}
func main() {
	//Init http server
	server := server.NewServer()

	//Init password handler & JWT auth variable
	passHandler := &service.PasswordHandler{SaltSize: 16}
	auth := &service.JwtAuth{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
		Issuer:    os.Getenv("JWT_ISSUER")}

	middleware.AuthMdw = middleware.NewAuthorMdw(auth)
	//Init API collections
	apiv1.NewUserAPI(server.Engine, service.NewUser(passHandler, auth))
	apiv1.NewRelationshipsAPI(server.Engine, service.NewRelationships())
	apiv1.NewNewsFeedAPI(server.Engine, service.NewNewsFeed())
	//Start http server
	server.Start("8080")
}
