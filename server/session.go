package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Server) createSession(c *gin.Context) {
	sess := s.sessions.Create()
	tokenStr, err := generateToken(sess)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "JWT creation failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func (s *Server) getSession(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(http.StatusOK, claims.(*sessionClaims).Session)
}
