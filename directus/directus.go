package directus

import (
	"sync"
	"time"
)

type Directus struct {
	baseurl      string
	token        string
	collections  []*Collection
	tags         []*Tag
	institutions []*Institution
	locations    []*Location
	news         []*News
	pages        []*Page
	mutex        sync.Mutex
	lastAccess   time.Time
	cacheTime    time.Duration
}

type Error struct {
	Message    string            `json:"message"`
	Extensions map[string]string `json:"extensions"`
}

func NewDirectus(baseurl, token string, cacheTime time.Duration) *Directus {
	d := &Directus{
		baseurl:     baseurl,
		token:       token,
		mutex:       sync.Mutex{},
		collections: nil,
		lastAccess:  time.Time{},
		cacheTime:   cacheTime,
	}
	return d
}

func (d *Directus) clearCache() {
	d.collections = nil
	d.tags = nil
	d.news = nil
	d.institutions = nil
	d.pages = nil
}
