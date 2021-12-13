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

type Institution struct {
	d           *Directus `json:"-"`
	Id          int64     `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Location    int64     `json:"location"`
	Contact     string    `json:"contact"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Description string    `json:"description"`
	Sort        int64     `json:"sort"`
}

type InstitutionsResult struct {
	Data   []*Institution `json:"data,omitempty"`
	Errors []*Error       `json:"errors"`
}

type sortInstitutionsBySort []*Institution

func (si sortInstitutionsBySort) Len() int           { return len(si) }
func (si sortInstitutionsBySort) Swap(i, j int)      { si[i], si[j] = si[j], si[i] }
func (si sortInstitutionsBySort) Less(i, j int) bool { return si[i].Sort < si[j].Sort }

func (d *Directus) loadInstitutions() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.institutions == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		if d.institutions != nil {
			d.clearCache()
		}
		urlStr := fmt.Sprintf("%s/items/institutions?limit=-1", d.baseurl)
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
		institutions := InstitutionsResult{}
		if err := json.Unmarshal(data, &institutions); err != nil {
			d.collections = nil
			return errors.Wrapf(err, "cannot parse result: %s", string(data))
		}
		if len(institutions.Errors) > 0 {
			d.collections = nil
			error := institutions.Errors[0]
			return errors.New(fmt.Sprintf("%s", error.Message))
		}
		d.institutions = []*Institution{}
		for _, inst := range institutions.Data {
			inst.d = d
			d.institutions = append(d.institutions, inst)
		}
		sort.Sort(sortInstitutionsBySort(d.institutions))
		d.lastAccess = time.Now()

	}
	return nil
}

func (d *Directus) GetInstitutions() ([]*Institution, error) {
	if err := d.loadInstitutions(); err != nil {
		return nil, err
	}
	return d.institutions, nil
}

func (d *Directus) GetInstitution(id int64) (*Institution, error) {
	if err := d.loadInstitutions(); err != nil {
		return nil, err
	}
	for _, inst := range d.institutions {
		if inst.Id == id {
			return inst, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("tag #%v not found", id))
}

func (i *Institution) GetLocation() (*Location, error) {
	return i.d.GetLocation(i.Location)
}
