package storage

import (
	"strings"
)

func isErrMessageContains(err error, msgs ...string) bool {
	if err == nil {
		return false
	}
	for _, msg := range msgs {
		if msg == "" {
			return false
		}
		if !strings.Contains(err.Error(), msg) {
			return false
		}
	}
	return true
}
