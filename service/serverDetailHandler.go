package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/je4/basel-collections/v2/directus"
	"net/http"
	"net/url"
	"strconv"
)

func (s *Server) detailHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	detailValues := url.Values{}

	vars := mux.Vars(req)

	var tag, institution, collection int64

	collectionStr := vars["collection"]
	if collectionStr != "" {
		collection, _ = strconv.ParseInt(collectionStr, 10, 64)
	}
	theCollection, err := s.dir.GetInstitution(collection)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("cannot get collection #%s: %v", collectionStr, err)))
		return
	}

	tagStr := req.URL.Query().Get("tag")
	if tagStr != "" {
		detailValues.Add("tag", tagStr)
		tag, _ = strconv.ParseInt(tagStr, 10, 64)
	}
	institutionStr := req.URL.Query().Get("institution")
	if institutionStr != "" {
		detailValues.Add("institution", institutionStr)
		institution, _ = strconv.ParseInt(institutionStr, 10, 64)
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

	gridLarge, lastRowLarge := buildGrid(GRIDLARGE, colls)
	gridSmall, lastRowSmall := buildGrid(GRIDLARGE, colls)

	theInstitution, err := theCollection.GetInstitution()

	impressumLarge := &Impressum{
		Id: 0, Left: 1, Cols: 12, Top: lastRowLarge, Rows: 3,
		Type:   "impressum",
		Scheme: IMPRESSUM,
		VAlign: "middle",
		Text:   "Impressum | Datenschutz | Informationen<br />(c) 2021 Basel Collections",
	}
	impressumSmall := &Impressum{
		Id: 0, Left: 1, Cols: 8, Top: lastRowSmall, Rows: 3,
		Type:   "impressum",
		Scheme: IMPRESSUM,
		VAlign: "middle",
		Text:   "Impressum | Datenschutz | Informationen<br />(c) 2021 Basel Collections",
	}

	tpl := s.templates["root"]
	if err := tpl.Execute(w, struct {
		GridLarge, GridSmall           []Grid
		ImpressumLarge, ImpressumSmall *Impressum
		Tags                           []*directus.Tag
		Institutions                   []*directus.Institution
		Collections                    []*directus.Collection
		Institution                    int64
		Tag                            int64
		DetailParam                    string
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
		DetailParam:    "?" + detailValues.Encode(),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s : %v", "root", err)))
		return
	}
}
