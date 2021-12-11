package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"strings"
	"time"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type CfgDatabase struct {
	ServerType string
	DSN        string
	ConnMax    int `toml:"connection_max"`
	Schema     string
}

type Endpoint struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Forward struct {
	Local  Endpoint `toml:"local"`
	Remote Endpoint `toml:"remote"`
}

type SSHTunnel struct {
	User       string             `toml:"user"`
	PrivateKey string             `toml:"privatekey"`
	Endpoint   Endpoint           `toml:"endpoint"`
	Forward    map[string]Forward `toml:"forward"`
}

type Config struct {
	ServiceName string `toml:"servicename"`
	Logfile     string `toml:"logfile"`
	Loglevel    string `toml:"loglevel"`
	Logformat   string `toml:"logformat"`
	AccessLog   string `toml:"accesslog"`
	Addr        string `toml:"addr"`
	AddrExt     string `toml:"addrext"`
	CertPEM     string `toml:"certpem"`
	KeyPEM      string `toml:"keypem"`
	//	JWTKey       string               `toml:"jwtkey"`
	//	JWTAlg       []string             `toml:"jwtalg"`
	DB     CfgDatabase          `toml:"database"`
	Tunnel map[string]SSHTunnel `toml:"tunnel"`
}

func LoadConfig(filepath string) Config {
	var conf Config
	conf.ServiceName = "BaselCollections"
	_, err := toml.DecodeFile(filepath, &conf)
	if err != nil {
		log.Fatalln("Error on loading config: ", err)
	}
	conf.AddrExt = strings.TrimRight(conf.AddrExt, "/")

	return conf
}
