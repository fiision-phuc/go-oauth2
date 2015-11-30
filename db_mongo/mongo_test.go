package mongo

import (
	"testing"

	"gopkg.in/mgo.v2"
)

func TestCreateSession(t *testing.T) {
	ConnectMongo(Config{Host: "", Port: "", Database: "", Username: "", Password: ""})
	ConnectMongo(Config{Host: "localhost", Port: "27017", Database: "concept"})
	ConnectMongo(Config{Host: "localhost", Port: "27017", Database: "concept", Username: "", Password: ""})
}

func TestGetEventualSession(t *testing.T) {
	session, database := GetEventualSession()

	if session.Mode() != mgo.Eventual {
		t.Error("Not Eventual Session.")
	}

	if database.Name != "concept" {
		t.Error("Invalid database name.")
	}
}
