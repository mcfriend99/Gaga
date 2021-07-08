package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
)

// ServerConfig configuration struct
type ServerConfig struct {
	Port               int    `json:"port"`
	ListenOn           string `json:"listen_on"`
	Secure             bool   `json:"secure,omitempty"`
	TLSCertificateFile string `json:"tls_certificate_file,omitempty"`
	TLSKeyFile         string `json:"tls_key_file,omitempty"`
}

// DatabaseConfig configuration struct
type DatabaseConfig struct {
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// LogConfig configuration struct
type LogConfig struct {
	// Engine is the device to log into.
	//
	// Options include:
	//  file, console, both
	Engine     string `json:"engine"`
	Path       string `json:"path,omitempty"`
	ShowSource bool   `json:"show_source,omitempty"`

	// Level is the minimum log level to show.
	//
	// Options include:
	//  info, error, warn, trace, fatal, panic.
	Level string `json:"level,omitempty"`
}

// SEOConfig configuration struct
type SEOConfig struct {
	// Compress indicates whether to compress output.
	// The output is only compressed when the requester (most likely
	// a browser supports it through the Accept-Encoding header.
	//
	//  gzip is the default if supported, if not,
	//  deflate will be attempted if supported.
	Compress bool `json:"compress,omitempty"`

	// CompressionThreshold is the minimum response size before
	// response compression sets in.
	//
	// For very light data (few bytes), compression may increase the
	// size of the response.
	CompressionThreshold int `json:"compression_threshold,omitempty"`
}

// Config is the main configuration struct
type Config struct {
	Server   ServerConfig `json:"server"`
	Database interface{}  `json:"database,omitempty"`
	Log      LogConfig    `json:"log,omitempty"`
	SEO      SEOConfig    `json:"seo,omitempty"`
	Custom   interface{}  `json:"custom,omitempty"`
}

// LoadConfig loads Gaga configurations from the specified file.
// It automatically allows environment overrides for dev and production
// as well as for different users on a device.
func LoadConfig() *Config {
	name := "config"
	config := Config{}

	var priority []string
	if u, err := user.Current(); err == nil {
		priority = []string{u.Username, "dev", ""}
	} else {
		priority = []string{"dev", ""}
	}

	matchFound := false
	for _, p := range priority {
		var n string
		if p != "" {
			n = name + "." + p + ".json"
		} else {
			n = name + ".json"
		}

		if bytes, err := ioutil.ReadFile(n); err == nil {
			matchFound = true
			if err = json.Unmarshal(bytes, &config); err != nil {
				log.Fatalln("Could not decode configuration file.")
			}
		}
	}

	if !matchFound {
		log.Fatalln("No suitable configuration  file found. " +
			"Try renaming config.sample.json to config.json")
	}

	return &config
}
