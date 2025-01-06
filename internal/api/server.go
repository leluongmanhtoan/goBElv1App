package httpServer

import (
	"net"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	engine := gin.New()
	engine.Static("/uploads", "./uploads")
	engine.Use(gin.Recovery())
	engine.Use(CORSMiddleware())
	server := &Server{Engine: engine}
	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

func (server *Server) Start(port string) {
	listener, err := net.Listen("tcp4", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	if err := server.Engine.RunListener(listener); err != nil {
		log.Printf("service start failed: %v", err)
	}

	/*if err := server.Engine.Run("0.0.0.0:" + port); err != nil {
		log.WithError(err).Error("service start failed")
	}*/
}
