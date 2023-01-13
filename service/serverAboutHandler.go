package service

import (
	"fmt"
	"github.com/je4/basel-collections/v2/directus"
	"net/http"
	"net/url"
	"strconv"
)

func (s *Server) aboutHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	detailValues := url.Values{}

	var tag, institution int64
	tagStr := req.URL.Query().Get("tag")
	if tagStr != "" {
		detailValues.Add("tag", tagStr)
		tag, err = strconv.ParseInt(tagStr, 10, 64)
	}
	institutionStr := req.URL.Query().Get("institution")
	if institutionStr != "" {
		detailValues.Add("institution", institutionStr)
		institution, err = strconv.ParseInt(institutionStr, 10, 64)
	}

	aboutPage, err := s.dir.GetPageByName("About")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get about about Page: %v", err)))
		return
	}
	sponsorPage, err := s.dir.GetPageByName("Sponsoren")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get about sponsor Page: %v", err)))
		return
	}
	tags, err := s.dir.GetTags()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get tags: %v", err)))
		return
	}
	institutions, err := s.dir.GetInstitutions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get tags: %v", err)))
		return
	}

	_locations, err := s.dir.GetLocations()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get locations: %v", err)))
		return
	}

	var locations = map[int64]*directus.Location{}
	for _, l := range _locations {
		locations[l.Id] = l
	}

	impressumLarge := &Impressum{
		Id: 0, Left: 1, Cols: 12, Top: 20, Rows: 3,
		Type:   "impressum",
		Scheme: IMPRESSUM,
		VAlign: "middle",
		Text:   "<a class=\"link\" href=\"impressum\">Impressum</a> | <a class=\"link\" href=\"datenschutz\">Datenschutz</a> | <a class=\"link\" href=\"about\">Information</a><br />(c) 2021 Basel Collections",
	}
	impressumSmall := &Impressum{
		Id: 0, Left: 1, Cols: 8, Top: 20, Rows: 3,
		Type:   "impressum",
		Scheme: IMPRESSUM,
		VAlign: "middle",
		Text:   "<a class=\"link\" href=\"impressum\">Impressum</a> | <a class=\"link\" href=\"datenschutz\">Datenschutz</a> | <a class=\"link\" href=\"about\">Information</a><br />(c) 2021 Basel Collections",
	}

	if s.templateReload {
		s.InitTemplates()
	}

	s.templateMutex.RLock()
	defer s.templateMutex.RUnlock()
	tpl := s.templates["about"]

	if err := tpl.Execute(w, struct {
		ImpressumLarge, ImpressumSmall *Impressum
		Tags                           []*directus.Tag
		Institutions                   []*directus.Institution
		Locations                      map[int64]*directus.Location
		Institution                    int64
		Tag                            int64
		DetailParam                    string
		LinkHome                       string
		LinkImpressum                  string
		LinkNews                       string
		LinkCollection                 string
		AboutContent                   string
		SponsorContent                 string
		BoxLarge                       Grid
		LinkAbout                      string
	}{
		ImpressumLarge: impressumLarge,
		ImpressumSmall: impressumSmall,
		Tags:           tags,
		Institutions:   institutions,
		Locations:      locations,
		Tag:            tag,
		Institution:    institution,
		DetailParam:    "?" + detailValues.Encode(),
		LinkHome:       "../",
		LinkImpressum:  "../impressum",
		LinkAbout:      "../about",
		LinkNews:       "../news",
		LinkCollection: "../detail",
		AboutContent:   aboutPage.Content,
		SponsorContent: sponsorPage.Content,
		BoxLarge:       Grid{Id: 0, Left: 1, Cols: 8, Top: 2, Rows: 2, Type: BoxImpressum, Scheme: SCHEMES[3], VAlign: bottom},
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s : %v", "root", err)))
		return
	}
}
