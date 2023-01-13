package service

import (
	"fmt"
	"github.com/je4/basel-collections/v2/directus"
	"net/http"
	"net/url"
	"strconv"
)

func (s *Server) newsHandler(w http.ResponseWriter, req *http.Request) {
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

	var news []*directus.News
	news, err = s.dir.GetNews()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get collections: %v", err)))
		return
	}
	var contents = []Content{}
	for _, n := range news {
		contents = append(contents, n)
	}
	tags, err := s.dir.GetTags()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get tags: %v", err)))
		return
	}
	_institutions, err := s.dir.GetInstitutions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get tags: %v", err)))
		return
	}

	var institutions = map[int64]*directus.Institution{}
	for _, i := range _institutions {
		institutions[i.Id] = i
	}

	gridLarge, lastRowLarge := buildGrid(NEWSLARGE, contents)
	gridSmall, lastRowSmall := buildGrid(NEWSLARGE, contents)

	var theBoxLarge Grid
	for _, box := range gridLarge {
		if box.Type == "news" {
			theBoxLarge = box
			break
		}
	}
	var theBoxSmall Grid
	for _, box := range gridSmall {
		if box.Type == "news" {
			theBoxSmall = box
			break
		}
	}

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
	tpl := s.templates["news"]

	if err := tpl.Execute(w, struct {
		GridLarge, GridSmall           []Grid
		ImpressumLarge, ImpressumSmall *Impressum
		Tags                           []*directus.Tag
		Institutions                   map[int64]*directus.Institution
		News                           []*directus.News
		Institution                    int64
		Tag                            int64
		DetailParam                    string
		LinkHome                       string
		LinkImpressum                  string
		LinkNews                       string
		LinkCollection                 string
		BoxLarge                       Grid
		BoxSmall                       Grid
		LinkAbout                      string
	}{
		GridLarge:      gridLarge,
		GridSmall:      gridSmall,
		BoxLarge:       theBoxLarge,
		BoxSmall:       theBoxSmall,
		ImpressumLarge: impressumLarge,
		ImpressumSmall: impressumSmall,
		News:           news,
		Tags:           tags,
		Institutions:   institutions,
		Tag:            tag,
		Institution:    institution,
		DetailParam:    "?" + detailValues.Encode(),
		LinkHome:       "../",
		LinkImpressum:  "../impressum",
		LinkAbout:      "../about",
		LinkNews:       "",
		LinkCollection: "../detail",
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s : %v", "root", err)))
		return
	}
}
