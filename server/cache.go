package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Mopsgamer/space-soup/server/environment"
	"github.com/Mopsgamer/space-soup/server/soup"
)

type SoupCache struct {
	ExpiresAt time.Time

	PlotImageBytes []byte
	TestList       []soup.MovementTest
}

func (cache *SoupCache) IsExpired() bool {
	return cache.ExpiresAt.Before(time.Now())
}

// Increases ExpiresAt time.
func (cache *SoupCache) Live() {
	cache.ExpiresAt = time.Now().Add(environment.ImageCacheDuration)
}

type FileHashCacheMap map[string]*SoupCache

func (m FileHashCacheMap) Free() {
	for hash, cache := range m {
		if cache.ExpiresAt.After(time.Now()) {
			continue
		}

		delete(m, hash)
	}
}

func (m FileHashCacheMap) Add(key string, cache SoupCache) SoupCache {
	cache.Live()
	m.Free()
	m[key] = &cache
	return cache
}

// Create identificator for file content.
func HashString(data []byte) string {
	hashBytes := sha256.Sum256(data)
	return hex.EncodeToString(hashBytes[:])
}
