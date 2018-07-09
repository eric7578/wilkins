package main

import (
	"os"

	"github.com/eric7578/wilkins/server"
	"github.com/eric7578/wilkins/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	host, password, db := storage.LoadRedisEnv()
	storage.InitClient(host, password, db)

	s := server.NewServer()
	s.Run()
}
