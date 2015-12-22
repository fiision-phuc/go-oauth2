package mongo

import (
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type user struct {
	UserID   bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"username,omitempty"`
	Password string        `bson:"password,omitempty"`
	Roles    []string      `bson:"roles,omitempty"`
}

func Test_EntityWithID(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()

	// Reset database
	session, database := GetMonotonicSession()
	database.C("Test").DropCollection()
	session.Close()

	// Testing process
	err := EntityWithID("", bson.NewObjectId(), nil)
	if err == nil {
		t.Error("Expected error not nil but found nil.")
	}
	if err != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err)
	}

	err = EntityWithID("Test", bson.NewObjectId(), nil)
	if err != "Invalid entity object." {
		t.Errorf("Expected \"%s\"  but found \"%s\"", "Invalid entity object.", err)
	}

	u := &user{}
	err = EntityWithID("Test", bson.NewObjectId(), u)
	if err == nil {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "not found", err)
	}
}
