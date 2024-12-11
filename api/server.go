package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())
	server := &Server{Engine: engine}
	return server
}

func (server *Server) Start(port string) {
	if err := server.Engine.Run(":" + port); err != nil {
		log.WithError(err).Error("service start failed")
	}
}
