package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

type CatalogueCollection struct {
	Id          int64  `json:"id"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Url         string `json:"url"`
}

type Catalogue struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	IdentifierField string `json:"identifier_field"`
	Collections     []struct {
		CollectionID struct {
			Id int64 `json:"id"`
		} `json:"collections_id"`
	} `json:"collections"`
}

type CataloguesResult struct {
	Data   []*Catalogue `json:"data,omitempty"`
	Errors []*Error     `json:"errors"`
}

func (d *Directus) loadCatalogues() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.catalogues == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		if d.catalogues != nil {
			d.clearCache()
		}
		urlStr := fmt.Sprintf("%s/items/catalogs?limit=-1&fields=*,collections.collections_id.id", d.baseurl)
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
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrapf(err, "cannot read result of %s", urlStr)
		}
		Catalogues := CataloguesResult{}
		if err := json.Unmarshal(data, &Catalogues); err != nil {
			d.collections = nil
			return errors.Wrapf(err, "cannot parse result: %s", string(data))
		}
		if len(Catalogues.Errors) > 0 {
			d.collections = nil
			error := Catalogues.Errors[0]
			return errors.New(fmt.Sprintf("%s", error.Message))
		}
		d.catalogues = []*Catalogue{}
		for _, inst := range Catalogues.Data {
			d.catalogues = append(d.catalogues, inst)
		}
		d.lastAccess = time.Now()
	}
	return nil
}

func (d *Directus) GetCatalogues() ([]*Catalogue, error) {
	if err := d.loadCatalogues(); err != nil {
		return nil, err
	}
	return d.catalogues, nil
}

func (d *Directus) GetCatalogue(id int64) (*Catalogue, error) {
	if err := d.loadCatalogues(); err != nil {
		return nil, err
	}
	for _, loc := range d.catalogues {
		if loc.Id == id {
			return loc, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("tag #%v not found", id))
}
