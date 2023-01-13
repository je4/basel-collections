package service

import (
	"fmt"
	"github.com/je4/basel-collections/v2/directus"
	"net/http"
	"net/url"
	"strconv"
)

type Impressum struct {
	Id, Left, Cols, Top, Rows int64
	Type                      string
	Scheme                    map[string]string
	VAlign                    string
	Text                      string
}

func (s *Server) collectionsHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var selected string
	detailValues := url.Values{}

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

	var tag, institution int64
	tagStr := req.URL.Query().Get("tag")
	if tagStr != "" {
		detailValues.Add("tag", tagStr)
		tag, err = strconv.ParseInt(tagStr, 10, 64)
		for _, t := range tags {
			if t.Id == tag {
				selected = t.Tag
				break
			}
		}
	}
	institutionStr := req.URL.Query().Get("institution")
	if institutionStr != "" {
		detailValues.Add("institution", institutionStr)
		institution, err = strconv.ParseInt(institutionStr, 10, 64)
		for _, i := range institutions {
			if i.Id == institution {
				selected = i.Name
				break
			}
		}
	}

	var colls []*directus.Collection
	if institution > 0 {
		colls, err = s.dir.GetCollectionsByInstitution(institution)
	} else {
		if tag > 0 {
			colls, err = s.dir.GetCollectionsByTags([]int64{tag})
		} else {
			colls, err = s.dir.GetCollections()
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get collections: %v", err)))
		return
	}
	var contents = []Content{}
	for _, c := range colls {
		contents = append(contents, c)
	}

	gridLarge, lastRowLarge := buildGrid(GRIDLARGE, contents)
	gridSmall, lastRowSmall := buildGrid(GRIDSMALL, contents)

	impressumLarge := &Impressum{
		Id: 0, Left: 1, Cols: 12, Top: lastRowLarge, Rows: 3,
		Type:   "impressum",
		Scheme: IMPRESSUM,
		VAlign: "middle",
		Text:   "<a class=\"link\" href=\"impressum\">Impressum</a> | <a class=\"link\" href=\"datenschutz\">Datenschutz</a> | <a class=\"link\" href=\"about\">Information</a><br />(c) 2021 Basel Collections",
	}
	impressumSmall := &Impressum{
		Id: 0, Left: 1, Cols: 8, Top: lastRowSmall, Rows: 3,
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
	tpl := s.templates["collections"]

	if err := tpl.Execute(w, struct {
		GridLarge, GridSmall           []Grid
		ImpressumLarge, ImpressumSmall *Impressum
		Tags                           []*directus.Tag
		Institutions                   []*directus.Institution
		Collections                    []*directus.Collection
		Institution                    int64
		Tag                            int64
		Selected                       string
		DetailParam                    string
		LinkHome                       string
		LinkNews                       string
		LinkImpressum                  string
		LinkAbout                      string
		LinkCollection                 string
	}{
		GridLarge:      gridLarge,
		GridSmall:      gridSmall,
		ImpressumLarge: impressumLarge,
		ImpressumSmall: impressumSmall,
		Collections:    colls,
		Tags:           tags,
		Institutions:   institutions,
		Tag:            tag,
		Institution:    institution,
		Selected:       selected,
		DetailParam:    "?" + detailValues.Encode(),
		LinkHome:       "",
		LinkImpressum:  "impressum",
		LinkAbout:      "about",
		LinkNews:       "news",
		LinkCollection: "detail",
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s : %v", "root", err)))
		return
	}
}
