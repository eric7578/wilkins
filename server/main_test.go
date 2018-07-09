package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eric7578/wilkins/storage"
)

var testServer *Server
var token string

func TestMain(m *testing.M) {
	host, password, db := storage.LoadRedisEnv()
	cli := storage.InitClient(host, password, db)
	testServer = NewServer()
	token = getTestToken()

	exit := m.Run()

	cli.FlushDB()

	os.Exit(exit)
}

func getTestToken() string {
	req, _ := http.NewRequest("POST", "/session/token", nil)
	w := httptest.NewRecorder()

	testServer.engine.ServeHTTP(w, req)
	bytes, _ := ioutil.ReadAll(w.Body)

	type RespBody struct {
		Token string `json:"token"`
	}
	var respBody = new(RespBody)
	json.Unmarshal(bytes, respBody)

	return respBody.Token
}

func getTestRequest(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", token)
	testServer.engine.ServeHTTP(w, req)
	return w
}
