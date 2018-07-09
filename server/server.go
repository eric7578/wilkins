package server

import (
	"net/http"
	"time"

	"github.com/eric7578/wilkins/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server reprsent server for Wilkins AP
type Server struct {
	engine   *gin.Engine
	auth     *storage.Auth
	sessions *storage.Sessions
	channels *storage.Channels
}

// NewServer get a new instance of Server
func NewServer() *Server {
	engine := gin.Default()
	s := Server{
		engine:   engine,
		auth:     storage.NewAuth(),
		sessions: storage.NewSessions(),
		channels: storage.NewChannels(),
	}

	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	engine.POST("/session/token", s.createSession)
	engine.Any("/ws", s.websocketUpgrader)
	engine.Use(s.checkToken())
	engine.GET("/session", s.getSession)
	engine.GET("/channel", s.listChannels)
	engine.POST("/channel", s.createChannel)
	engine.GET("/channel/:channelID", s.readMessage)
	engine.POST("/channel/:channelID", s.postMessage)
	engine.POST("/channel/:channelID/participants", s.joinChannel)

	return &s
}

// Run register api routes and websocket route
func (s *Server) Run() {
	s.engine.Run()
}

func (s *Server) abortWithError(c *gin.Context, err error) {
	switch {
	case storage.IsErrNotAllowed(err):
		c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
	case storage.IsErrNotFound(err):
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
	default:
		// TODO
		// seperate production/develop mode
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	}
}
