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

type Page struct {
	dir     *Directus `json:"-"`
	Id      int64     `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Content string    `json:"content,omitempty"`
}

type PageResult struct {
	Data   []*Page  `json:"data,omitempty"`
	Errors []*Error `json:"errors"`
}

func (n *Page) GetId() int64    { return n.Id }
func (n *Page) GetName() string { return n.Name }
func (d *Directus) GetPages() ([]*Page, error) {
	if err := d.loadPages(); err != nil {
		return nil, errors.Wrap(err, "cannot load news")
	}
	return d.pages, nil
}

func (d *Directus) GetPageByName(name string) (*Page, error) {
	pages, err := d.GetPages()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load news")
	}
	for _, p := range pages {
		if p.GetName() == name {
			return p, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("page %s not found", name))
}

func (d *Directus) loadPages() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.pages == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		if d.pages != nil {
			d.clearCache()
		}
		urlStr := "%s/items/page"
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
		page := PageResult{}
		if err := json.Unmarshal(data, &page); err != nil {
			d.news = nil
			return errors.Wrapf(err, "cannot parse result: %s", string(data))
		}
		if len(page.Errors) > 0 {
			d.news = nil
			error := page.Errors[0]
			return errors.New(fmt.Sprintf("%s", error.Message))
		}
		d.pages = []*Page{}
		for _, page := range page.Data {
			page.dir = d
			d.pages = append(d.pages, page)
		}
		d.lastAccess = time.Now()
	}
	return nil
}
