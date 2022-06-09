package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	err := checkArgs(nil)
	assert.Error(err)

	plugin.Url = "http://example.com"
	err = checkArgs(nil)
	assert.NoError(err)

}

func TestSendRequest(t *testing.T) {
	assert := assert.New(t)
	requestBody = strings.NewReader("This is fun")
	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		expectedBody := `This is fun`
		assert.Equal("POST", r.Method)
		assert.Contains(string(body), expectedBody)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"ok": true}`))
		require.NoError(t, err)
	}))

	plugin.Url = apiStub.URL
	plugin.Method = "POST"
	err := sendRequest(nil)
	assert.NoError(err)
}
