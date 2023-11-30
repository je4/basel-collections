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

type Location struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Street     string `json:"street"`
	Zip        string `json:"zip"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Additional string `json:"additional"`
	Position   string `json:"position"`
}

type LocationsResult struct {
	Data   []*Location `json:"data,omitempty"`
	Errors []*Error    `json:"errors"`
}

func (d *Directus) loadLocations() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.locations == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		if d.locations != nil {
			d.clearCache()
		}
		urlStr := fmt.Sprintf("%s/items/locations?limit=-1", d.baseurl)
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
		locations := LocationsResult{}
		if err := json.Unmarshal(data, &locations); err != nil {
			d.collections = nil
			return errors.Wrapf(err, "cannot parse result: %s", string(data))
		}
		if len(locations.Errors) > 0 {
			d.collections = nil
			error := locations.Errors[0]
			return errors.New(fmt.Sprintf("%s", error.Message))
		}
		d.locations = []*Location{}
		for _, inst := range locations.Data {
			d.locations = append(d.locations, inst)
		}
		d.lastAccess = time.Now()
	}
	return nil
}

func (d *Directus) GetLocations() ([]*Location, error) {
	if err := d.loadLocations(); err != nil {
		return nil, err
	}
	return d.locations, nil
}

func (d *Directus) GetLocation(id int64) (*Location, error) {
	if err := d.loadLocations(); err != nil {
		return nil, err
	}
	for _, loc := range d.locations {
		if loc.Id == id {
			return loc, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("tag #%v not found", id))
}
