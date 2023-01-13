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

type NewsType struct {
	Id   int64  `json:"id"`
	Key  string `json:"key"`
	Type string `json:"type"`
}

type News struct {
	dir         *Directus `json:"-"`
	start       time.Time
	end         time.Time
	publishDate time.Time
	Id          int64    `json:"id,omitempty"`
	Status      string   `json:"status,omitempty"`
	Sort        int64    `json:"sort,omitempty"`
	DateCreated string   `json:"date_created,omitempty"`
	DateUpdated string   `json:"date_updated,omitempty"`
	Title       string   `json:"title,omitempty"`
	Image       string   `json:"image,omitempty"`
	Text        string   `json:"text,omitempty"`
	Start       string   `json:"start,omitempty"`
	End         string   `json:"end,omitempty"`
	PublishDate string   `json:"publish_date,omitempty"`
	Url         string   `json:"url,omitempty"`
	Institution int64    `json:"institution,omitempty"`
	Type        NewsType `json:"type,omitempty"`
}

func (n *News) GetTagIDs() []int64 {
	return []int64{}
}

func (n *News) GetInstitutionID() int64 {
	return 0
}

func absDuration(d time.Duration) time.Duration {
	if d > 0 {
		return d
	} else {
		return -d
	}
}
func minDuration(d1, d2 time.Duration) time.Duration {
	if d1 < d2 {
		return d1
	} else {
		return d2
	}
}

type sortNews []*News

func (sc sortNews) Len() int      { return len(sc) }
func (sc sortNews) Swap(i, j int) { sc[i], sc[j] = sc[j], sc[i] }
func (sc sortNews) Less(i, j int) bool {
	now := time.Now()
	iNews := sc[i]
	jNews := sc[j]

	iDiff := absDuration(now.Sub(iNews.publishDate))
	iDiff = minDuration(iDiff, absDuration(now.Sub(iNews.start)))
	iDiff = minDuration(iDiff, absDuration(now.Sub(iNews.end)))

	jDiff := absDuration(now.Sub(jNews.publishDate))
	jDiff = minDuration(jDiff, absDuration(now.Sub(jNews.start)))
	jDiff = minDuration(jDiff, absDuration(now.Sub(jNews.end)))

	return iDiff < jDiff
}

type NewsResult struct {
	Data   []*News  `json:"data,omitempty"`
	Errors []*Error `json:"errors"`
}

func (n *News) GetId() int64 { return n.Id }
func (n *News) GetDate() string {
	result := n.start.Format("02.01.2006")
	if n.Start != n.End {
		result += " bis " + n.end.Format("02.01.2006")
	}
	return result
}
func (n *News) GetTitle() string { return n.Title }
func (n *News) GetInst() string {
	inst, err := n.GetInstitution()
	if err != nil {
		return err.Error()
	}
	return inst.Name
}
func (n *News) GetInstitution() (*Institution, error) {
	return n.dir.GetInstitution(n.Institution)
}
func (n *News) GetUrl() string {
	return n.Url
}

func (d *Directus) GetNews() ([]*News, error) {
	if err := d.loadNews(); err != nil {
		return nil, errors.Wrap(err, "cannot load news")
	}
	return d.news, nil
}

func (d *Directus) GetNewsByInstitution(institution int64) ([]*News, error) {
	if err := d.loadNews(); err != nil {
		return nil, err
	}
	var news = []*News{}
	for _, n := range d.news {
		if n.Institution == institution {
			news = append(news, n)
		}
	}
	return news, nil
}

func (d *Directus) loadNews() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.news == nil || time.Now().Add(-d.cacheTime).After(d.lastAccess) {
		if d.news != nil {
			d.clearCache()
		}
		//		urlStr := fmt.Sprintf("%s/items/news?filter[status][_eq]=published&filter[publish_date][_lte]=%s&filter[end][_gte]=%s&limit=-1&fields=*,type.*", d.baseurl, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
		urlStr := fmt.Sprintf("%s/items/news?filter[status][_eq]=published&filter[publish_date][_lte]=%s&limit=-1&fields=*,type.*", d.baseurl, time.Now().Format("2006-01-02"))
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
		news := NewsResult{}
		if err := json.Unmarshal(data, &news); err != nil {
			d.news = nil
			return errors.Wrapf(err, "cannot parse result: %s", string(data))
		}
		if len(news.Errors) > 0 {
			d.news = nil
			error := news.Errors[0]
			return errors.New(fmt.Sprintf("%s", error.Message))
		}
		d.news = []*News{}
		for _, news := range news.Data {
			news.dir = d
			if news.start, err = time.Parse("2006-01-02", news.Start); err != nil {
				return errors.Wrapf(err, "cannot parse start time - %s", news.Start)
			}
			if news.end, err = time.Parse("2006-01-02", news.End); err != nil {
				return errors.Wrapf(err, "cannot parse end time - %s", news.End)
			}
			if news.publishDate, err = time.Parse("2006-01-02", news.PublishDate); err != nil {
				return errors.Wrapf(err, "cannot parse publishDate time - %s", news.PublishDate)
			}
			d.news = append(d.news, news)
		}
		sort.Sort(sortNews(d.news))
		d.lastAccess = time.Now()
	}
	return nil
}
