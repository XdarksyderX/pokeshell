package pokeshell

import (
	"time"
)

type CacheEntry struct {
	Data      []byte
	CreatedAt time.Time
}

type PokeCache struct {
	Cache map[string]CacheEntry
	Limit time.Duration
}

func (pC *PokeCache) Set(link string, data []byte) {
	pC.Cache[link] = CacheEntry{
		data,
		time.Now(),
	}
}

func (pC *PokeCache) Get(link string) ([]byte, error) {
	pC.CheckLoop()
	if cE, ok := pC.Cache[link]; ok {
		if time.Now().Sub(cE.CreatedAt) > pC.Limit {
			delete(pC.Cache, link)
			return nil, nil
		}
		return cE.Data, nil
	}

	fetched, err := fetchAPI(link)
	if err != nil {
		return nil, err
	}
	pC.Set(link, fetched)
	return fetched, nil
}

func (pC *PokeCache) CheckLoop() {
	for link, entry := range pC.Cache {
		if time.Now().Sub(entry.CreatedAt) > pC.Limit {
			delete(pC.Cache, link)
		}
	}
}

func NewPokeCache(limit time.Duration) *PokeCache {
	return &PokeCache{
		Cache: make(map[string]CacheEntry),
		Limit: limit,
	}
}
