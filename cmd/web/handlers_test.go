package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heisenberg8055/gosts/internal/assert"
)

func TestPingHandler(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(healthCheck)

	handler.ServeHTTP(rr, r)

	rs := rr.Result()

	t.Run("Status Code", func(t *testing.T) { assert.Equal(t, rs.StatusCode, http.StatusOK) })

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	t.Run("Response Body", func(t *testing.T) {
		assert.Equal(t, string(body), "OK")
	})

}

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	defer ts.Close()

	code, _, body := ts.get(t, "/healthz")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
