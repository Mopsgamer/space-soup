package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Mopsgamer/space-soup/server/environment"
)

type TimeLimitedCache struct {
	ExpiresAt time.Time
	Bytes     []byte
}

func (cache *TimeLimitedCache) Hash() string {
	hashBytes := sha256.Sum256(cache.Bytes)
	return hex.EncodeToString(hashBytes[:])
}

func (cache *TimeLimitedCache) IsExpired() bool {
	return cache.ExpiresAt.Before(time.Now())
}

// Increases ExpiresAt time.
func (cache *TimeLimitedCache) Live() {
	cache.ExpiresAt = time.Now().Add(environment.ImageCacheDuration)
}

type VisualizationCache map[string]*TimeLimitedCache

func (m VisualizationCache) Free() {
	for hash, cache := range m {
		if cache.ExpiresAt.Before(time.Now()) {
			continue
		}

		delete(m, hash)
	}
}

func (m VisualizationCache) Add(bytes []byte) TimeLimitedCache {
	cache := TimeLimitedCache{Bytes: bytes}
	cache.Live()
	m.Free()
	m[cache.Hash()] = &cache
	return cache
}
