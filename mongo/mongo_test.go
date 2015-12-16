package mongo

import (
	"os"
	"testing"

	"gopkg.in/mgo.v2"
)

func Test_CreateSession(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()
}

func Test_GetEventualSession(t *testing.T) {
	session, database := GetEventualSession()

	if session.Mode() != mgo.Eventual {
		t.Error("Not Eventual Session.")
	}

	if database.Name != "mongo" {
		t.Error("Invalid database name.")
	}
}
