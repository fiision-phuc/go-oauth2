package mongo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/phuc0302/go-oauth2/utils"
)

// ConfigFile defines configuration file's name.
const ConfigFile = "mongodb.cfg"

////////////////////////////////////////////////////////////////////////////////////////////////////

// MongoConfig descripts a configuration  object  that  will  be  used  during application life time.
type MongoConfig struct {
	Addresses []string `json:"addresses"`

	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateConfigs generates a default configuration file.
func CreateMongoConfigs() {
	if utils.FileExisted(ConfigFile) {
		os.Remove(ConfigFile)
	}

	host := os.Getenv("OPENSHIFT_MONGODB_DB_HOST")
	port := os.Getenv("OPENSHIFT_MONGODB_DB_PORT")
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	if len(port) == 0 {
		port = "27017"
	}

	// Create default config
	config := MongoConfig{
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
func LoadMongoConfigs() *MongoConfig {
	if !utils.FileExisted(ConfigFile) {
		CreateMongoConfigs()
	}

	file, err := os.Open(ConfigFile)
	if err != nil {
		return nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	config := MongoConfig{}
	json.Unmarshal(bytes, &config)

	return &config
}
