package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/bluele/gcache"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/je4/basel-collections/v2/directus"
	"github.com/je4/basel-collections/v2/files"
	dcert "github.com/je4/utils/v2/pkg/cert"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Server struct {
	service    string
	host, port string
	addrExt    string
	srv        *http.Server
	log        *logging.Logger
	accessLog  io.Writer
	templates  map[string]*template.Template
	cache      gcache.Cache
	dir        *directus.Directus
}

func NewServer(service, addr, addrExt string, dir *directus.Directus, log *logging.Logger, accessLog io.Writer) (*Server, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot split address %s", addr)
	}
	/*
		extUrl, err := url.Parse(addrExt)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse external address %s", addrExt)
		}
	*/

	srv := &Server{
		service:   service,
		host:      host,
		port:      port,
		addrExt:   addrExt,
		dir:       dir,
		log:       log,
		accessLog: accessLog,
		templates: map[string]*template.Template{},
		cache:     gcache.New(500).ARC().Build(),
	}

	return srv, srv.InitTemplates()
}

func (s *Server) InitTemplates() error {
	for key, val := range files.TemplateFiles {
		tpl, err := template.New("root.gohtml").Funcs(sprig.FuncMap()).ParseFS(files.TemplateFS, val)
		if err != nil {
			return errors.Wrapf(err, "cannot parse template %s - %s:", key, val)
		}
		s.templates[key] = tpl
	}
	return nil
}

func (s *Server) ListenAndServe(cert, key string) (err error) {
	router := mux.NewRouter()

	fsys, err := fs.Sub(files.StaticFS, "static")
	if err != nil {
		return errors.Wrap(err, "cannot get subtree of embedded static")
	}
	httpStaticServer := http.FileServer(http.FS(fsys))
	router.PathPrefix("/static").Handler(
		handlers.CompressHandler(func(prefix string, h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fullPrefix := fmt.Sprintf("/%s", prefix)
				p := strings.TrimPrefix(r.URL.Path, fullPrefix)
				rp := strings.TrimPrefix(r.URL.RawPath, fullPrefix)
				if len(p) <= len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
					r2 := new(http.Request)
					*r2 = *r
					r2.URL = new(url.URL)
					*r2.URL = *r.URL
					r2.URL.Path = p
					r2.URL.RawPath = rp
					w.Header().Set("Cache-Control", "max-age=3600")
					h.ServeHTTP(w, r2)
				} else {
					http.NotFound(w, r)
				}
			})

		}("static", httpStaticServer),
		// http.StripPrefix("/static", httpStaticServer)
		),
	).Methods("GET")

	router.HandleFunc("/", s.rootHandler).Methods("GET")

	loggedRouter := handlers.CombinedLoggingHandler(s.accessLog, handlers.ProxyHeaders(router))
	addr := net.JoinHostPort(s.host, s.port)
	s.srv = &http.Server{
		Handler: loggedRouter,
		Addr:    addr,
	}

	if cert == "auto" || key == "auto" {
		s.log.Info("generating new certificate")
		cert, err := dcert.DefaultCertificate()
		if err != nil {
			return errors.Wrap(err, "cannot generate default certificate")
		}
		s.srv.TLSConfig = &tls.Config{Certificates: []tls.Certificate{*cert}}
		s.log.Infof("starting server at https://%s:%s - %s", s.host, s.port, s.addrExt)
		return s.srv.ListenAndServeTLS("", "")
	} else if cert != "" && key != "" {
		s.log.Infof("starting server at https://%s:%s - %s", s.host, s.port, s.addrExt)
		return s.srv.ListenAndServeTLS(cert, key)
	} else {
		s.log.Infof("starting server at http://%s:%s - %s", s.host, s.port, s.addrExt)
		return s.srv.ListenAndServe()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
