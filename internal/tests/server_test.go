package tests

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/MohamedBabker/project1-jwks/internal/httpserver"
	"github.com/MohamedBabker/project1-jwks/internal/keystore"
)

func newTestServer(t *testing.T) (*httpserver.Server, *keystore.Store) {
	t.Helper()
	ks, err := keystore.NewDefaultStore()
	if err != nil {
		t.Fatalf("keystore init: %v", err)
	}
	s := httpserver.New(ks)
	return s, ks
}

func TestHealthz(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
}

func TestJWKSOnlyUnexpired(t *testing.T) {
	s, ks := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	var payload struct {
		Keys []struct {
			Kid string `json:"kid"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(payload.Keys) != 1 {
		t.Fatalf("expected 1 unexpired key, got %d", len(payload.Keys))
	}
	if payload.Keys[0].Kid != ks.ActiveKey(time.Now()).KID {
		t.Fatalf("kid mismatch")
	}
}

func TestAuthReturnsJWTWithKID(t *testing.T) {
	s, ks := newTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	var payload struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	tokenStr := payload.Token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("parse jwt: %v", err)
	}
	if token.Header["kid"] != ks.ActiveKey(time.Now()).KID {
		t.Fatalf("expected kid %s got %v", ks.ActiveKey(time.Now()).KID, token.Header["kid"])
	}
}

func TestAuthExpiredFlow(t *testing.T) {
	s, ks := newTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/auth?expired=1", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	var payload struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	token, _, err := new(jwt.Parser).ParseUnverified(payload.Token, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("parse jwt: %v", err)
	}
	if token.Header["kid"] != ks.ExpiredKey().KID {
		t.Fatalf("expected expired kid %s got %v", ks.ExpiredKey().KID, token.Header["kid"])
	}
	claims := token.Claims.(jwt.MapClaims)
	expF, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("exp missing")
	}
	if time.Unix(int64(expF), 0).After(time.Now()) {
		t.Fatalf("expected expired exp claim")
	}
}

func TestJWKSModulusLooksRight(t *testing.T) {
	s, ks := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)

	var payload struct {
		Keys []struct {
			N string `json:"n"`
			E string `json:"e"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(payload.Keys) != 1 {
		t.Fatalf("expected 1 key")
	}
	activePub := ks.ActiveKey(time.Now()).Key.Public().(*rsa.PublicKey)
	nB64 := payload.Keys[0].N
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	if len(nBytes) != len(activePub.N.Bytes()) {
		t.Fatalf("modulus length mismatch")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 got %d", rr.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/.well-known/jwks.json", strings.NewReader("{}"))
	rr = httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 got %d", rr.Code)
	}
}
