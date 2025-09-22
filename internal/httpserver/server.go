
package httpserver

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/MohamedBabker/project1-jwks/internal/keystore"
	"github.com/MohamedBabker/project1-jwks/jwks"
)

type Server struct {
	ks  *keystore.Store
	mux *http.ServeMux
}

func New(ks *keystore.Store) *Server {
	s := &Server{ks: ks, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	s.mux.HandleFunc("/.well-known/jwks.json", s.handleJWKS())
	s.mux.HandleFunc("/auth", s.handleAuth())
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.logRequests(s.mux))
}

func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleJWKS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		pubKeys := s.ks.UnexpiredPublicKeys(time.Now())
		keys := make([]jwks.JWK, 0, len(pubKeys))
		for _, k := range pubKeys {
			pub := k.Key.Public().(*rsa.PublicKey)
			keys = append(keys, jwks.FromRSAPublicKey(pub, k.KID))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks.JWKS{Keys: keys})
	}
}

func (s *Server) handleAuth() http.HandlerFunc {
	type resp struct{ Token string `json:"token"` }
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		now := time.Now()
		useExpired := r.URL.Query().Has("expired")

		var kp keystore.KeyPair
		if useExpired {
			kp = s.ks.ExpiredKey()
		} else {
			kp = s.ks.ActiveKey(now)
		}

		claims := jwt.MapClaims{
			"sub": "fake-user",
			"iat": now.Unix(),
		}
		if useExpired {
			claims["exp"] = now.Add(-15 * time.Minute).Unix()
		} else {
			claims["exp"] = now.Add(15 * time.Minute).Unix()
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = kp.KID

		signed, err := token.SignedString(kp.Key)
		if err != nil {
			http.Error(w, "failed to sign", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp{Token: signed})
	}
}

func (s *Server) Handler() http.Handler { return s.mux }
