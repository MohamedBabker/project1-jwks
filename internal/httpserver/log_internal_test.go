package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MohamedBabker/project1-jwks/internal/keystore"
)

// This covers the unexported logRequests middleware so it counts toward coverage.
func TestLogRequestsWrapper(t *testing.T) {
	ks, err := keystore.NewDefaultStore()
	if err != nil {
		t.Fatalf("keystore init: %v", err)
	}
	s := New(ks)

	called := false
	wrapped := s.logRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent) // 204
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	if !called {
		t.Fatalf("inner handler was not called")
	}
}
