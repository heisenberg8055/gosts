package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heisenberg8055/gosts/internal/assert"
)

func TestSecureHeaders(t *testing.T) {

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	secureHeaders(http.HandlerFunc(healthCheck)).ServeHTTP(rr, r)

	rs := rr.Result()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "Content-Security-Policy",
			want: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			name: "Referrer-Policy",
			want: "origin-when-cross-origin",
		},
		{
			name: "X-Content-Type-Options",
			want: "nosniff",
		},
		{
			name: "X-Frame-Options",
			want: "deny",
		},
		{
			name: "X-XSS-Protection",
			want: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, rs.Header.Get(tt.name), tt.want)
		})
	}
	t.Run("Status Code", func(t *testing.T) {
		assert.Equal(t, rs.StatusCode, http.StatusOK)
	})
}
