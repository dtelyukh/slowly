package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestSlowBadRequestNoData(t *testing.T) {
	logger, _ := test.NewNullLogger()

	handler := NewHandler(
		logger,
		&Config{MaxTimeout: 100},
	)

	req := httptest.NewRequest(
		"POST",
		"http://example.com/api/slow",
		nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"Invalid timeout value"}`, string(body))
}

func TestSlowBadRequestThrottle(t *testing.T) {
	logger, _ := test.NewNullLogger()

	handler := NewHandler(
		logger,
		&Config{MaxTimeout: 100},
	)

	req := httptest.NewRequest(
		"POST",
		"http://example.com/api/slow",
		strings.NewReader(`{"timeout":101}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"error":"timeout too long"}`, string(body))
}

func TestSlowSuccess(t *testing.T) {
	logger, _ := test.NewNullLogger()

	handler := NewHandler(
		logger,
		&Config{MaxTimeout: 200},
	)

	req := httptest.NewRequest(
		"POST",
		"http://example.com/api/slow",
		strings.NewReader(`{"timeout":100}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"status":"ok"}`, string(body))
}
