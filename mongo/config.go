package mongo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/phuc0302/go-oauth2/util"
)

// ConfigFile defines configuration file's name.
const ConfigFile = "mongodb.cfg"

////////////////////////////////////////////////////////////////////////////////////////////////////

// Config descripts a configuration  object  that  will  be  used  during application life time.
type Config struct {
	Addresses []string `json:"addresses"`
	Database  string   `json:"database"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

// CreateConfigs generates a default configuration file.
func CreateConfigs() {
	/* Condition validation */
	if util.FileExisted(ConfigFile) {
		os.Remove(ConfigFile)
	}

	// Create default config
	host := "127.0.0.1"
	port := "27017"
	config := Config{
		Addresses: []string{fmt.Sprintf("%s:%s", host, port)},
		Database:  "mongo",
		Username:  "",
		Password:  "",
	}

	// Create new file
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	file, _ := os.Create(ConfigFile)
	file.Write(configJSON)
	file.Close()
}

// LoadConfigs retrieves previous configuration from file.
func LoadConfigs() *Config {
	/* Condition validation */
	if !util.FileExisted(ConfigFile) {
		CreateConfigs()
	}

	file, err := os.Open(ConfigFile)
	if err != nil {
		return nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	config := Config{}
	json.Unmarshal(bytes, &config)

	return &config
}
