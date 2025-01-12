package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0hlov3/FedSplitDomainChecker/internal/logger"
	"github.com/stretchr/testify/assert"
)

// Initialize the logger for tests
func TestMain(m *testing.M) {
	logger.Init(true) // Enable verbose logging for tests
	m.Run()
}

func TestMakeRequest(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusMovedPermanently)
		} else if r.URL.Path == "/final" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	resp, err := makeRequest(mockServer.URL + "/redirect")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)

	resp, err = makeRequest(mockServer.URL + "/final")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = makeRequest(mockServer.URL + "/non-existent")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
