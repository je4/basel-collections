package service

import (
	"fmt"
	"net/http"
)

func (s *Server) rootHandler(w http.ResponseWriter, req *http.Request) {
	tpl := s.templates["root"]
	if err := tpl.Execute(w, struct {
		Error   string
		Message string
	}{Error: "doiError", Message: "message"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s : %v", "root", err)))
		return
	}
}
