package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type CollectionTag struct {
	Id            int64 `json:"id"`
	CollectionsId int64 `json:"collections_id"`
	TagsId        int64 `json:"tags_id"`
}

type Collection struct {
	dir         *Directus       `json:"-"`
	Id          int64           `json:"id,omitempty"`
	Status      string          `json:"status,omitempty"`
	Sort        int64           `json:"sort,omitempty"`
	DateCreated string          `json:"date_created,omitempty"`
	DateUpdated string          `json:"date_updated,omitempty"`
	Title       string          `json:"title,omitempty"`
	Image       string          `json:"image,omitempty"`
	Description string          `json:"description,omitempty"`
	Institution int64           `json:"institution,omitempty"`
	Url         string          `json:"url,omitempty"`
	Hinweis     string          `json:"hinweis,omitempty"`
	Tags        []CollectionTag `json:"tags,omitempty"`
}

type sortCollectionsBySort []*Collection

func (sc sortCollectionsBySort) Len() int           { return len(sc) }
func (sc sortCollectionsBySort) Swap(i, j int)      { sc[i], sc[j] = sc[j], sc[i] }
func (sc sortCollectionsBySort) Less(i, j int) bool { return sc[i].Sort < sc[j].Sort }

type CollectionsResult struct {
	Data   []*Collection `json:"data,omitempty"`
	Errors []*Error      `json:"errors"`
}

func (c *Collection) GetTags() ([]*Tag, error) {
	list := []int64{}
	for _, ct := range c.Tags {
		if ct.TagsId == 0 {
			continue
		}
		list = append(list, ct.TagsId)
	}
	return c.dir.GetTagList(list)
}

func (c *Collection) GetInstitution() (*Institution, error) {
	return c.dir.GetInstitution(c.Institution)
}

func (d *Directus) GetCollections() ([]*Collection, error) {
	if err := d.loadCollections(); err != nil {
		return nil, err
	}
	return d.collections, nil
}

func (d *Directus) GetCollectionsByInstitution(institution int64) ([]*Collection, error) {
	if err := d.loadCollections(); err != nil {
		return nil, err
	}
	var collections = []*Collection{}
	for _, coll := range d.collections {
		if coll.Institution == institution {
			collections = append(collections, coll)
		}
	}
	return collections, nil
}

func (d *Directus) GetCollectionsByTags(tags []int64) ([]*Collection, error) {
	if err := d.loadCollections(); err != nil {
		return nil, err
	}
	var collections = []*Collection{}
	for _, coll := range d.collections {
		var tagCount int
		for _, t := range tags {
			for _, t2 := range coll.Tags {
				if t == t2.Id {
					tagCount++
				}
			}
		}
		if tagCount == len(tags) {
			collections = append(collections, coll)
		}
	}
	return collections, nil
}

func (d *Directus) GetCollection(id int64) (*Collection, error) {
	if err := d.loadCollections(); err != nil {
		return nil, err
	}
	for _, coll := range d.collections {
		if coll.Id == id {
			return coll, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("collection #%v not found", id))
}

func (d *Directus) loadCollections() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.collections == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		if d.collections != nil {
			d.clearCache()
		}
		urlStr := fmt.Sprintf("%s/items/collections?filter[status][_eq]=published&limit=-1&fields=*,tags.*", d.baseurl)
		req, err := http.NewRequest("GET", urlStr, bytes.NewReader(nil))
		if err != nil {
			return errors.Wrapf(err, "cannot create request %s", urlStr)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", d.token))
		client := http.Client{
			Transport:     http.DefaultTransport,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "error executing %s", urlStr)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrapf(err, "cannot read result of %s", urlStr)
		}
		collections := CollectionsResult{}
		if err := json.Unmarshal(data, &collections); err != nil {
			d.collections = nil
			return errors.Wrapf(err, "cannot parse result: %s", string(data))
		}
		if len(collections.Errors) > 0 {
			d.collections = nil
			error := collections.Errors[0]
			return errors.New(fmt.Sprintf("%s", error.Message))
		}
		d.collections = []*Collection{}
		for _, coll := range collections.Data {
			coll.dir = d
			d.collections = append(d.collections, coll)
		}
		sort.Sort(sortCollectionsBySort(d.collections))
		d.lastAccess = time.Now()
	}
	return nil
}
