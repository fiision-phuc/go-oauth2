package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/util"
)

// Config file's name.
const (
	debug   = "server.debug.cfg"
	release = "server.release.cfg"
)

// HTTP Methods.
const (
	Copy    = "copy"
	Delete  = "delete"
	Get     = "get"
	Head    = "head"
	Link    = "link"
	Options = "options"
	Patch   = "patch"
	Post    = "post"
	Purge   = "purge"
	Put     = "put"
	Unlink  = "unlink"
)

// Config describes a configuration object that will be used during application life time.
type Config struct {
	// Server
	Host    string `json:"host"`
	Port    int    `json:"port"`
	TLSPort int    `json:"tls_port"`

	// Header
	HeaderSize    int           `json:"header_size"`    // In KB
	MultipartSize int64         `json:"multipart_size"` // In MB
	ReadTimeout   time.Duration `json:"timeout_read"`   // In seconds
	WriteTimeout  time.Duration `json:"timeout_write"`  // In seconds

	// HTTP Method
	AllowMethods  []string          `json:"allow_methods"`
	RedirectPaths map[string]string `json:"redirect_paths"`
	StaticFolders map[string]string `json:"static_folders"`

	// Log
	LogLevel     string `json:"log_level"`
	SlackURL     string `json:"slack_url"`
	SlackIcon    string `json:"slack_icon"`
	SlackUser    string `json:"slack_user"`
	SlackChannel string `json:"slack_channel"`
}

// CreateConfig generates a default configuration file.
func CreateConfig(configFile string) {
	if util.FileExisted(configFile) {
		os.Remove(configFile)
	}

	// Create default config
	config := Config{
		Host:    "localhost",
		Port:    8080,
		TLSPort: 8443,

		HeaderSize:    5,
		MultipartSize: 1,
		ReadTimeout:   15,
		WriteTimeout:  15,

		AllowMethods: []string{Copy, Delete, Get, Head, Link, Options, Patch, Post, Purge, Put, Unlink},
		RedirectPaths: map[string]string{
			"401": "/login",
		},
		StaticFolders: map[string]string{
			"/assets":    "assets",
			"/resources": "resources",
		},

		LogLevel:     "debug",
		SlackURL:     "",
		SlackIcon:    ":ghost:",
		SlackUser:    "Server",
		SlackChannel: "#channel",
	}

	// Create new file
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	file, _ := os.Create(configFile)
	file.Write(configJSON)
	file.Close()
}

// LoadConfig retrieves previous configuration from file.
func LoadConfig(configFile string) Config {
	// Generate config file if neccessary
	if !util.FileExisted(configFile) {
		CreateConfig(configFile)
	}

	// Load config file
	config := Config{}
	file, _ := os.Open(configFile)
	bytes, _ := ioutil.ReadAll(file)

	if err := json.Unmarshal(bytes, &config); err == nil {
		// Convert duration to seconds
		config.HeaderSize <<= 10
		config.MultipartSize <<= 20
		config.ReadTimeout *= time.Second
		config.WriteTimeout *= time.Second

		// Define redirectPaths
		redirectPaths = make(map[int]string, len(config.RedirectPaths))
		for s, path := range config.RedirectPaths {
			if status, err := strconv.Atoi(s); err == nil {
				redirectPaths[status] = path
			}
		}

		// Define regular expressions
		methodsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.AllowMethods, "|")))
	}
	return config
}
