
package keystore

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/google/uuid"
)

type KeyPair struct {
	KID    string
	Key    *rsa.PrivateKey
	Expiry time.Time
}

type Store struct {
	active  KeyPair
	expired KeyPair
}

func NewDefaultStore() (*Store, error) {
	activeKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	expiredKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Store{
		active: KeyPair{
			KID:    uuid.NewString(),
			Key:    activeKey,
			Expiry: now.Add(1 * time.Hour),
		},
		expired: KeyPair{
			KID:    uuid.NewString(),
			Key:    expiredKey,
			Expiry: now.Add(-1 * time.Hour),
		},
	}, nil
}

func (s *Store) ActiveKey(now time.Time) KeyPair {
	if now.After(s.active.Expiry) {
		k, _ := rsa.GenerateKey(rand.Reader, 2048)
		s.active = KeyPair{
			KID:    uuid.NewString(),
			Key:    k,
			Expiry: now.Add(1 * time.Hour),
		}
	}
	return s.active
}

func (s *Store) ExpiredKey() KeyPair { return s.expired }

func (s *Store) UnexpiredPublicKeys(now time.Time) []KeyPair {
	keys := []KeyPair{}
	if now.Before(s.active.Expiry) {
		keys = append(keys, s.active)
	}
	return keys
}
