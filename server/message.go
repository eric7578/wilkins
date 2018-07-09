package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/eric7578/wilkins/packet"
	"github.com/gin-gonic/gin"
)

func (s *Server) listChannels(c *gin.Context) {
	claims, _ := c.Get("claims")
	sessionID := claims.(*sessionClaims).Session.ID

	channels, err := s.channels.FindBySession(sessionID)
	if err != nil {
		s.abortWithError(c, err)
		return
	}

	infos := make([]*packet.ChannelInfo, len(channels))
	sessions := make(map[string]*packet.Session)
	for channelIdx, channel := range channels {
		info := packet.NewChannelInfo(channel)
		infos[channelIdx] = info
		for sessionIdx, sessionID := range channel.Sessions {
			if _, exist := sessions[sessionID]; !exist {
				session, _ := s.sessions.Get(sessionID)
				if session != nil {
					sessions[sessionID] = session
				}
			}
			info.Sessions[sessionIdx] = sessions[sessionID]
		}
	}

	c.JSON(http.StatusOK, infos)
}

func (s *Server) getChannel(c *gin.Context) {
	claims, _ := c.Get("claims")
	sessionID := claims.(*sessionClaims).Session.ID

	channels, err := s.channels.FindBySession(sessionID)
	if err != nil {
		s.abortWithError(c, err)
		return
	}

	infos := make([]*packet.ChannelInfo, len(channels))
	sessions := make(map[string]*packet.Session)
	for channelIdx, channel := range channels {
		info := packet.NewChannelInfo(channel)
		infos[channelIdx] = info
		for sessionIdx, sessionID := range channel.Sessions {
			if _, exist := sessions[sessionID]; !exist {
				session, _ := s.sessions.Get(sessionID)
				if session != nil {
					sessions[sessionID] = session
				}
			}
			info.Sessions[sessionIdx] = sessions[sessionID]
		}
	}

	c.JSON(http.StatusOK, infos)
}

func (s *Server) createChannel(c *gin.Context) {
	claims, _ := c.Get("claims")
	sessionID := claims.(*sessionClaims).Session.ID

	if os.Getenv("ORPHAN_PROTECTION") != "0" && s.channels.HasOrphan(sessionID) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ch := s.channels.Create(sessionID)
	c.JSON(http.StatusOK, ch)
}

func (s *Server) postMessage(c *gin.Context) {
	var body packet.MessagesBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid body")
		return
	}
	if err := packet.IsValidMessageString(body.Messages); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	claims, _ := c.Get("claims")
	sessionID := claims.(*sessionClaims).Session.ID
	channelID := c.Param("channelID")
	if err := s.auth.CanSessionAccessChannel(sessionID, channelID); err != nil {
		s.abortWithError(c, err)
		return
	}

	if err := s.channels.PostTo(sessionID, channelID, body.Messages); err != nil {
		s.abortWithError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) joinChannel(c *gin.Context) {
	claims, _ := c.Get("claims")
	sessionID := claims.(*sessionClaims).Session.ID
	channelID := c.Param("channelID")

	if s.channels.InChannel(sessionID, channelID) {
		c.Status(http.StatusOK)
		return
	}

	err := s.channels.Join(sessionID, channelID)
	if err != nil {
		s.abortWithError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) readMessage(c *gin.Context) {
	claims, _ := c.Get("claims")
	sessionID := claims.(*sessionClaims).Session.ID
	channelID := c.Param("channelID")
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid offset")
		return
	}

	if err := s.auth.CanSessionAccessChannel(sessionID, channelID); err != nil {
		s.abortWithError(c, err)
		return
	}

	messages, err := s.channels.ReadMessages(sessionID, channelID, offset)
	if err != nil {
		s.abortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}
