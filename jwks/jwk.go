
package jwks

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"strings"
)

type JWK struct {
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func FromRSAPublicKey(pub *rsa.PublicKey, kid string) JWK {
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes())
	return JWK{
		Kty: "RSA",
		N:   strings.TrimRight(n, "="),
		E:   strings.TrimRight(e, "="),
		Alg: "RS256",
		Use: "sig",
		Kid: kid,
	}
}
