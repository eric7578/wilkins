package server

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/eric7578/wilkins/packet"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type sessionClaims struct {
	jwt.StandardClaims
	Session *packet.Session
}

func generateToken(sess *packet.Session) (string, error) {
	claims := sessionClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Session: sess,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (sessionClaims, error) {
	claims := sessionClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// check alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return claims, err
	} else if !token.Valid {
		return claims, errors.New("Invalid token")
	}

	return claims, claims.Valid()
}

func (s *Server) checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.Request.Header.Get("Authorization")
		claims, err := validateToken(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, "Invalid token")
		} else {
			c.Set("claims", &claims)
			c.Next()
		}
	}
}
