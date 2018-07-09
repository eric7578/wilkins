package storage

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	host, password, db := LoadRedisEnv()
	cli := InitClient(host, password, db)

	exit := m.Run()

	cli.FlushDB()

	os.Exit(exit)
}
