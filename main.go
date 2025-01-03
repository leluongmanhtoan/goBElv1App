package main

import (
	"log"
	"os"
	httpServer "program/internal/api"
	apiv1 "program/internal/api/v1"
	"program/internal/database"
	"program/internal/middleware"
	"strconv"

	authenticationRepo "program/internal/repositories/auth"
	newsfeedRepo "program/internal/repositories/newfeed"
	relationshipsRepo "program/internal/repositories/relationships"
	userRepo "program/internal/repositories/user"
	"program/internal/services"

	"github.com/joho/godotenv"
)

// Init unexported global variables
var mySqlConn database.ISqlConnection
var myRedisConn database.IRedisConnection

func init() {
	// Read .env file
	if err := godotenv.Load("./config/.env"); err != nil {
		log.Fatalln("env file load failed")
	}

	port, _ := strconv.Atoi(os.Getenv("SQLPort"))
	// Create sql connection
	sqlClientConfig := database.SqlConfig{
		Host:         os.Getenv("SQLHost"),
		Port:         port,
		Database:     os.Getenv("SQLDb"),
		Username:     os.Getenv("SQLUser"),
		Password:     os.Getenv("SQLPass"),
		ReadTimeout:  30,
		WriteTimeout: 30,
		MaxIdleConns: 10,
		MaxOpenConns: 10,
	}

	// Select sql connect to Mysql Database with config
	mySqlConn = database.NewMySqlConnection(&sqlClientConfig)
	// Run connect to Mysql database
	if err := mySqlConn.Connect(); err != nil {
		log.Fatal(err)
	}

	// Create redis connection
	redisClientConfig := database.RedisConfig{
		Addr:     os.Getenv("RedisHost"),
		Password: os.Getenv("RedisPass"),
		Database: 0,
	}
	// Select redis connect to Redis database with config
	myRedisConn = database.NewRedisConnection(&redisClientConfig)

	// Run connect to Redis database
	if err := myRedisConn.Connect(); err != nil {
		log.Fatal(err)
	}
}

func setup() {

	//Init service

}
func main() {
	// Init repository
	authRepo := authenticationRepo.NewAuthenticationRepo(myRedisConn)
	// Init auth repo config
	PassHandler := &services.PasswordHandler{SaltSize: 16}
	auth := &services.JwtAuthService{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
		Issuer:    os.Getenv("JWT_ISSUER"),
		Repo:      authRepo,
	}

	userRepo := userRepo.NewUserRepo(mySqlConn)
	relationshipsRepo := relationshipsRepo.NewRelationshipsRepo(mySqlConn)
	newsfeedRepo := newsfeedRepo.NewNewsfeedRepo(mySqlConn)

	// Init service

	userServices := services.NewUserService(userRepo, PassHandler, auth)
	relationshipsService := services.NewRelationshipsService(relationshipsRepo)
	newsfeedService := services.NewNewsFeedService(newsfeedRepo)

	// Init middleware service
	middleware.AuthMdw = middleware.NewAuthorMdw(auth)

	//Init http server
	server := httpServer.NewServer()

	//Init API collections
	apiv1.NewUserAPI(server.Engine, userServices)
	apiv1.NewRelationshipsAPI(server.Engine, relationshipsService)
	apiv1.NewNewsFeedAPI(server.Engine, newsfeedService)
	//Start http server
	server.Start("8080")
}
