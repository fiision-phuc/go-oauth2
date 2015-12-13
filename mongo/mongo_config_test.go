package mongo

import (
	"os"
	"reflect"
	"testing"

	"github.com/phuc0302/go-oauth2/utils"
)

func Test_CreateMongoConfigs(t *testing.T) {
	CreateMongoConfigs()
	defer os.Remove(ConfigFile)

	if !utils.FileExisted(ConfigFile) {
		t.Errorf("Expected %s file had been created but found nil.", ConfigFile)
	}
}

func Test_LoadMongoConfigs(t *testing.T) {
	defer os.Remove(ConfigFile)
	config := LoadMongoConfigs()

	if config == nil {
		t.Errorf("Expected not nil when %s is not available.", ConfigFile)
	}

	addresses := []string{"127.0.0.1:27017"}
	if !reflect.DeepEqual(addresses, config.Addresses) {
		t.Errorf("Expected %s but found %s", addresses, config.Addresses)
	}
}
