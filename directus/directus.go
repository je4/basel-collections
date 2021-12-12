package directus

import (
	"sync"
	"time"
)

type Directus struct {
	baseurl          string
	token            string
	collections      map[int64]Collection
	collectionsMutex sync.Mutex
	lastAccess       time.Time
	cacheTime        time.Duration
}

type Error struct {
	Message    string            `json:"message"`
	Extensions map[string]string `json:"extensions"`
}

type ErrorResult struct {
	Errors []Error `json:"errors"`
}

func NewDirectus(baseurl, token string, cacheTime time.Duration) *Directus {
	d := &Directus{
		baseurl:          baseurl,
		token:            token,
		collectionsMutex: sync.Mutex{},
		collections:      nil,
		lastAccess:       time.Time{},
		cacheTime:        cacheTime,
	}
	return d
}
