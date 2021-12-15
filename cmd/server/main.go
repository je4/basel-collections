package main

import (
	"context"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/je4/basel-collections/v2/directus"
	"github.com/je4/basel-collections/v2/service"
	lm "github.com/je4/utils/v2/pkg/logger"
	"github.com/je4/utils/v2/pkg/ssh"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfgFile := flag.String("cfg", "/etc/basel-collections.toml", "locations of config file")
	flag.Parse()
	config := LoadConfig(*cfgFile)

	// create logger instance
	logger, lf := lm.CreateLogger("Basel-Collections", config.Logfile, nil, config.Loglevel, config.Logformat)
	defer lf.Close()

	var err error

	var tunnels []*ssh.SSHtunnel
	for name, tunnel := range config.Tunnel {
		logger.Infof("starting tunnel %s", name)

		forwards := make(map[string]*ssh.SourceDestination)
		for fwName, fw := range tunnel.Forward {
			forwards[fwName] = &ssh.SourceDestination{
				Local: &ssh.Endpoint{
					Host: fw.Local.Host,
					Port: fw.Local.Port,
				},
				Remote: &ssh.Endpoint{
					Host: fw.Remote.Host,
					Port: fw.Remote.Port,
				},
			}
		}

		t, err := ssh.NewSSHTunnel(
			tunnel.User,
			tunnel.PrivateKey,
			&ssh.Endpoint{
				Host: tunnel.Endpoint.Host,
				Port: tunnel.Endpoint.Port,
			},
			forwards,
			logger,
		)
		if err != nil {
			logger.Errorf("cannot create tunnel %v@%v:%v - %v", tunnel.User, tunnel.Endpoint.Host, tunnel.Endpoint.Port, err)
			return
		}
		if err := t.Start(); err != nil {
			logger.Errorf("cannot create configfile %v - %v", t.String(), err)
			return
		}
		tunnels = append(tunnels, t)
	}
	defer func() {
		for _, t := range tunnels {
			logger.Infof("closing ssh tunnel")
			t.Close()
		}
	}()
	// if tunnels are made, wait until connection is established
	if len(config.Tunnel) > 0 {
		time.Sleep(1 * time.Second)
	}

	/*
		// get database connection handle
		db, err := sql.Open(config.DB.ServerType, config.DB.DSN)
		if err != nil {
			log.Fatalf("error opening database: %v", err)
		}
		// close on shutdown
		defer db.Close()

		// Open doesn't open a connection. Validate DSN data:
		err = db.Ping()
		if err != nil {
			log.Fatalf("error pinging database: %v", err)
		}
	*/

	var accessLog io.Writer
	var f *os.File
	if config.AccessLog == "" {
		accessLog = os.Stdout
	} else {
		f, err = os.OpenFile(config.AccessLog, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logger.Panicf("cannot open file %s: %v", config.AccessLog, err)
			return
		}
		defer f.Close()
		accessLog = f
	}

	dir := directus.NewDirectus(config.Directus.BaseUrl, config.Directus.Token, config.Directus.CacheTime.Duration)
	/*
		colls, err := dir.GetCollections()
		if err != nil {
			logger.Fatalf("cannot get collections: %v", err)
		}
		logger.Infof("%v", colls)
		for _, coll := range colls {
			tags, err := coll.GetTags()
			if err != nil {
				logger.Fatalf("cannot get tags of collections %s", coll.Title)
			}
			logger.Infof("coll: %s, tags: %v", coll.Title, tags)
		}
	*/
	srv, err := service.NewServer(config.ServiceName, config.Addr, config.AddrExt, dir, logger, accessLog)
	if err != nil {
		logger.Panicf("cannot initialize server: %v", err)
	}
	go func() {
		if err := srv.ListenAndServe(config.CertPEM, config.KeyPEM); err != nil {
			log.Fatalf("server died: %v", err)
		}
	}()

	end := make(chan bool, 1)

	// process waiting for interrupt signal (TERM or KILL)
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)

		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGKILL)

		<-sigint

		// We received an interrupt signal, shut down.
		logger.Infof("shutdown requested")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.Shutdown(ctx)

		end <- true
	}()

	<-end
	logger.Info("server stopped")

}
