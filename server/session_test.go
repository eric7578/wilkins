package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_createSession(t *testing.T) {
	assert.NotEmpty(t, token)

	t.Run("GET /session should return status ok when token is valid", func(t *testing.T) {
		w := getTestRequest("GET", "/session", nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET /session should return status fobidden when token is invalid", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/session", nil)
		req.Header.Add("Authorization", "foo")
		w := httptest.NewRecorder()

		testServer.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, "\"Invalid token\"", w.Body.String())
	})
}
