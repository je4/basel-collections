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
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Server struct {
	service        string
	host, port     string
	addrExt        string
	staticFiles    string
	templateFiles  string
	templateReload bool
	templateMutex  sync.RWMutex
	srv            *http.Server
	log            *logging.Logger
	accessLog      io.Writer
	templates      map[string]*template.Template
	cache          gcache.Cache
	dir            *directus.Directus
}

func NewServer(service, addr, addrExt, staticFiles, templateFiles string, templateReload bool, dir *directus.Directus, log *logging.Logger, accessLog io.Writer) (*Server, error) {
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
		service:        service,
		host:           host,
		port:           port,
		addrExt:        addrExt,
		staticFiles:    staticFiles,
		templateFiles:  templateFiles,
		templateReload: templateReload,
		dir:            dir,
		templateMutex:  sync.RWMutex{},
		log:            log,
		accessLog:      accessLog,
		templates:      map[string]*template.Template{},
		cache:          gcache.New(500).ARC().Build(),
	}

	return srv, srv.InitTemplates()
}

func (s *Server) InitTemplates() error {
	funcs := sprig.FuncMap()
	funcs["raw"] = func(s string) template.HTML { return template.HTML(s) }

	s.templateMutex.Lock()
	defer s.templateMutex.Unlock()

	var templateFS fs.FS
	var baseDir string
	if s.templateFiles != "" {
		templateFS = os.DirFS(s.templateFiles)
	} else {
		templateFS = files.TemplateFS
		baseDir = "template"
	}

	header := filepath.ToSlash(filepath.Join(baseDir, "header_content.inc.gohtml"))
	file := filepath.ToSlash(filepath.Join(baseDir, "collections.gohtml"))
	tpl, err := template.New("collections.gohtml").Funcs(funcs).ParseFS(templateFS, header, file)
	if err != nil {
		return errors.Wrapf(err, "cannot parse template %s - %s:", "collections", file)
	}
	s.templates["collections"] = tpl

	file = filepath.ToSlash(filepath.Join(baseDir, "news.gohtml"))
	tpl, err = template.New("news.gohtml").Funcs(funcs).ParseFS(templateFS, header, file)
	if err != nil {
		return errors.Wrapf(err, "cannot parse template %s - %s:", "news", file)
	}
	s.templates["news"] = tpl

	file = filepath.ToSlash(filepath.Join(baseDir, "impressum.gohtml"))
	tpl, err = template.New("impressum.gohtml").Funcs(funcs).ParseFS(templateFS, header, file)
	if err != nil {
		return errors.Wrapf(err, "cannot parse template %s - %s:", "impressum", file)
	}
	s.templates["impressum"] = tpl

	file = filepath.ToSlash(filepath.Join(baseDir, "kontakt.gohtml"))
	tpl, err = template.New("kontakt.gohtml").Funcs(funcs).ParseFS(templateFS, header, file)
	if err != nil {
		return errors.Wrapf(err, "cannot parse template %s - %s:", "kontakt", file)
	}
	s.templates["kontakt"] = tpl

	file = filepath.ToSlash(filepath.Join(baseDir, "information.gohtml"))
	tpl, err = template.New("information.gohtml").Funcs(funcs).ParseFS(templateFS, header, file)
	if err != nil {
		return errors.Wrapf(err, "cannot parse template %s - %s:", "information", file)
	}
	s.templates["information"] = tpl

	file = filepath.ToSlash(filepath.Join(baseDir, "datenschutz.gohtml"))
	tpl, err = template.New("datenschutz.gohtml").Funcs(funcs).ParseFS(templateFS, header, file)
	if err != nil {
		return errors.Wrapf(err, "cannot parse template %s - %s:", "datenschutz", file)
	}
	s.templates["datenschutz"] = tpl

	file = filepath.ToSlash(filepath.Join(baseDir, "collection.gohtml"))
	tpl, err = template.New("collection.gohtml").Funcs(funcs).ParseFS(templateFS, header, file)
	if err != nil {
		return errors.Wrapf(err, "cannot parse template %s - %s:", "detail", file)
	}
	s.templates["collection"] = tpl
	/*
		if s.templateFiles != "" {
			header := path.Join(s.templateFiles, "header_content.inc.gohtml")
			file := path.Join(s.templateFiles, "collections.gohtml")
			tpl, err := template.New("collections.gohtml").Funcs(funcs).ParseFiles(header, file)
			if err != nil {
				return errors.Wrapf(err, "cannot parse template %s - %s:", "collections", file)
			}
			s.templates["collections"] = tpl
			file = path.Join(s.templateFiles, "news.gohtml")
			tpl, err = template.New("news.gohtml").Funcs(funcs).ParseFiles(header, file)
			if err != nil {
				return errors.Wrapf(err, "cannot parse template %s - %s:", "news", file)
			}
			s.templates["news"] = tpl
			file = path.Join(s.templateFiles, "collection.gohtml")
			tpl, err = template.New("collection.gohtml").Funcs(funcs).ParseFiles(header, file)
			if err != nil {
				return errors.Wrapf(err, "cannot parse template %s - %s:", "detail", file)
			}
			s.templates["collection"] = tpl

		} else {
			tpl, err := template.New("collections.gohtml").Funcs(funcs).ParseFiles(files.CollectionsTemplate, files.HeaderContentIncTemplate)
			if err != nil {
				return errors.Wrapf(err, "cannot parse template %s - %s:", "collections", files.CollectionsTemplate)
			}
			s.templates["collections"] = tpl
			tpl, err = template.New("news.gohtml").Funcs(funcs).ParseFiles(files.NewsTemplate, files.HeaderContentIncTemplate)
			if err != nil {
				return errors.Wrapf(err, "cannot parse template %s - %s:", "news", files.NewsTemplate)
			}
			s.templates["news"] = tpl
			tpl, err = template.New("collection.gohtml").Funcs(funcs).ParseFiles(files.CollectionTemplate, files.HeaderContentIncTemplate)
			if err != nil {
				return errors.Wrapf(err, "cannot parse template %s - %s:", "detail", files.CollectionTemplate)
			}
			s.templates["collection"] = tpl
		}

	*/
	return nil
}

func (s *Server) ListenAndServe(cert, key string) error {
	router := mux.NewRouter()
	var fsys fs.FS
	var err error
	if s.staticFiles != "" {
		fsys = os.DirFS(s.staticFiles)
	} else {
		fsys, err = fs.Sub(files.StaticFS, "static")
		if err != nil {
			return errors.Wrap(err, "cannot get subtree of embedded static")
		}
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

	router.HandleFunc("/", s.collectionsHandler).Methods("GET")
	router.HandleFunc("/news", s.newsHandler).Methods("GET")
	router.HandleFunc("/impressum", s.impressumHandler).Methods("GET")
	router.HandleFunc("/information", s.informationHandler).Methods("GET")
	router.HandleFunc("/datenschutz", s.datenschutzHandler).Methods("GET")
	router.HandleFunc("/kontakt", s.kontaktHandler).Methods("GET")
	router.HandleFunc("/detail/{collection}", s.collectionHandler).Methods("GET")

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
