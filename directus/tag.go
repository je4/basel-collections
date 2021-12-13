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

type Tag struct {
	Id  int64  `json:"id"`
	Tag string `json:"tag"`
}

type TagsResult struct {
	Data   []*Tag   `json:"data,omitempty"`
	Errors []*Error `json:"errors"`
}

func (d *Directus) loadTags() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.tags == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		if d.tags != nil {
			d.clearCache()
		}
		urlStr := fmt.Sprintf("%s/items/tags?limit=-1", d.baseurl)
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
		tags := TagsResult{}
		if err := json.Unmarshal(data, &tags); err != nil {
			d.collections = nil
			return errors.Wrapf(err, "cannot parse result: %s", string(data))
		}
		if len(tags.Errors) > 0 {
			d.collections = nil
			error := tags.Errors[0]
			return errors.New(fmt.Sprintf("%s", error.Message))
		}
		d.tags = []*Tag{}
		for _, tag := range tags.Data {
			d.tags = append(d.tags, tag)
		}
		d.lastAccess = time.Now()
	}
	return nil
}

func (d *Directus) GetTags() ([]*Tag, error) {
	if err := d.loadTags(); err != nil {
		return nil, err
	}
	return d.tags, nil
}

func (d *Directus) GetTag(id int64) (*Tag, error) {
	if err := d.loadTags(); err != nil {
		return nil, err
	}
	for _, tag := range d.tags {
		if tag.Id == id {
			return tag, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("tag #%v not found", id))
}

func (d *Directus) GetTagList(ids []int64) ([]*Tag, error) {
	if err := d.loadTags(); err != nil {
		return nil, err
	}
	result := []*Tag{}
	for _, id := range ids {
		found := false
		for _, tag := range d.tags {
			if tag.Id == id {
				found = true
				result = append(result, tag)
			}
		}
		if !found {
			return nil, errors.New(fmt.Sprintf("tag #%v not found", id))
		}
	}
	return result, nil
}
