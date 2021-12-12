package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Collection struct {
	Id          int64    `json:"id,omitempty"`
	Status      string   `json:"status,omitempty"`
	Sort        int64    `json:"sort,omitempty"`
	DateCreated string   `json:"date_created,omitempty"`
	DateUpdated string   `json:"date_updated,omitempty"`
	Title       string   `json:"title,omitempty"`
	Image       string   `json:"image,omitempty"`
	Description string   `json:"description,omitempty"`
	Institution int64    `json:"institution,omitempty"`
	Url         string   `json:"url,omitempty"`
	Hinweis     string   `json:"hinweis,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (d *Directus) GetCollections() (map[int64]Collection, error) {
	d.collectionsMutex.Lock()
	defer d.collectionsMutex.Unlock()

	if d.collections == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		d.collections = nil
		urlStr := fmt.Sprintf("%s/items/collections?limit=-1", d.baseurl)
		req, err := http.NewRequest("GET", urlStr, bytes.NewReader(nil))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot create request %s", urlStr)
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
			return nil, errors.Wrapf(err, "error executing %s", urlStr)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read result of %s", urlStr)
		}
		collections := []Collection{}
		if err := json.Unmarshal(data, &collections); err != nil {
			d.collections = nil
			errorResult := ErrorResult{}
			if err := json.Unmarshal(data, &errorResult); err != nil {
				return nil, errors.Wrapf(err, "cannot parse result: %s", string(data))
			}
			if len(errorResult.Errors) > 0 {
				error := errorResult.Errors[0]
				return nil, errors.New(fmt.Sprintf("%s", error.Message))
			} else {
				return nil, errors.New(fmt.Sprintf("unknown error in result: %s", string(data)))
			}
		}
		d.collections = map[int64]Collection{}
		for _, coll := range collections {
			d.collections[coll.Id] = coll
		}
		d.lastAccess = time.Now()
	}
	return d.collections, nil
}
