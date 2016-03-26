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

func Test_AllEntities(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()

	// Reset database
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C("Test")
	defer collection.DropCollection()

	collection.Insert(
		&user{
			UserID:   bson.NewObjectId(),
			Username: "test1",
			Password: "test1",
			Roles:    []string{"r_user"},
		},
		&user{
			UserID:   bson.NewObjectId(),
			Username: "test2",
			Password: "test2",
			Roles:    []string{"r_user"},
		},
	)

	var list []user
	err := AllEntities("", &list)
	if err.Error() != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err.Error())
	}

	err = AllEntities("Test", nil)
	if err.Error() != "Invalid list object." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid list object.", err.Error())
	}

	err = AllEntities("Test", &list)
	if len(list) != 2 {
		t.Errorf("Expected %d but found %d.", 2, len(list))
	}
	if list[0].Username != "test1" {
		t.Errorf("Expected \"%s\" but found \"%s\".", "test1", list[0].Username)
	}
}

func Test_AllEntitiesWithCriteria(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()

	// Reset database
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C("Test")
	defer collection.DropCollection()

	collection.Insert(
		&user{
			UserID:   bson.NewObjectId(),
			Username: "test1",
			Password: "test1",
			Roles:    []string{"r_user"},
		},
		&user{
			UserID:   bson.NewObjectId(),
			Username: "test1",
			Password: "test1",
			Roles:    []string{"r_user"},
		},
		&user{
			UserID:   bson.NewObjectId(),
			Username: "test2",
			Password: "test2",
			Roles:    []string{"r_user"},
		},
	)

	var list []user
	err := AllEntitiesWithCriteria("", nil, nil)
	if err.Error() != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err.Error())
	}

	err = AllEntitiesWithCriteria("Test", nil, nil)
	if err.Error() != "Invalid criterion object." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid criterion object.", err.Error())
	}

	err = AllEntitiesWithCriteria("Test", map[string]interface{}{"username": "test1"}, nil)
	if err.Error() != "Invalid list object." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid list object.", err.Error())
	}

	err = AllEntitiesWithCriteria("Test", map[string]interface{}{"username": "test1"}, &list)
	if len(list) != 2 {
		t.Errorf("Expected %d but found %d.", 2, len(list))
	}
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
	if err.Error() != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err.Error())
	}

	err = EntityWithID("Test", bson.NewObjectId(), nil)
	if err.Error() != "Invalid entity object." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid entity object.", err.Error())
	}

	u := &user{}
	err = EntityWithID("Test", bson.NewObjectId(), u)
	if err == nil {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "not found", err.Error())
	}
}

func Test_EntityWithCriteria(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()

	// Reset database
	session, database := GetMonotonicSession()
	database.C("Test").DropCollection()
	session.Close()

	// Testing process
	err := EntityWithCriteria("", nil, nil)
	if err.Error() != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err.Error())
	}

	err = EntityWithCriteria("Test", nil, nil)
	if err.Error() != "Invalid criterion object." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid criterion object.", err.Error())
	}

	err = EntityWithCriteria("Test", bson.M{"_id": bson.NewObjectId()}, nil)
	if err.Error() != "Invalid entity object." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid entity object.", err.Error())
	}

	u := &user{}
	id := bson.NewObjectId()
	err = EntityWithID("Test", id, u)
	if err == nil {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "not found", err.Error())
	}

	recordUser := &user{}
	err = EntityWithCriteria("Test", bson.M{"_id": id}, recordUser)
	if recordUser.UserID != id {
		t.Errorf("Expected \"%s\" but found \"%s\".", id.Hex(), recordUser.UserID.Hex())
	}
}

func Test_SaveEntity(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()

	// Reset database
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C("Test")
	collection.DropCollection()

	// Testing process
	err := SaveEntity("", bson.NewObjectId(), nil)
	if err.Error() != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err.Error())
	}

	// [Case insert] we can skip the entityID
	err = SaveEntity("Test", bson.NewObjectId(), nil)
	if err.Error() != "Invalid entity object." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid entity object.", err.Error())
	}

	u := &user{
		UserID:   bson.NewObjectId(),
		Username: "admin",
		Password: "admin",
		Roles:    []string{"r_user", "r_admin"},
	}
	err = SaveEntity("Test", u.UserID, u)
	if err != nil {
		t.Errorf("Expected no error but found \"%s\".", err.Error())
	}

	// [Case update] we cannot skip the entityID
	u.Username = "admin1"
	u.Password = "admin1"
	err = SaveEntity("Test", u.UserID, u)
	if err != nil {
		t.Errorf("Expected no error but found \"%s\".", err)
	}

	var results []user
	err = collection.Find(nil).All(&results)
	if err != nil {
		t.Errorf("Expected no error but found \"%s\".", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 but found %d.", len(results))
	}
}

func Test_DeleteEntity(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()

	// Reset database
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C("Test")
	collection.DropCollection()

	u := &user{
		UserID:   bson.NewObjectId(),
		Username: "admin",
		Password: "admin",
		Roles:    []string{"r_user", "r_admin"},
	}
	SaveEntity("Test", u.UserID, u)

	// Testing process
	err := DeleteEntity("", bson.NewObjectId())
	if err.Error() != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err.Error())
	}

	err = DeleteEntity("Test", bson.NewObjectId())
	if err.Error() != "not found" {
		t.Errorf("Expected \"%s\" but found \"%s\".", "not found", err.Error())
	}

	err = DeleteEntity("Test", u.UserID)
	if err != nil {
		t.Errorf("Expected no error but found none \"%s\".", err.Error())
	}

	var results []user
	err = collection.Find(nil).All(&results)

	if len(results) != 0 {
		t.Errorf("Expected 0 but found %d.", len(results))
	}
}

func Test_DeleteEntityWithCriteria(t *testing.T) {
	defer os.Remove(ConfigFile)
	ConnectMongo()

	// Reset database
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C("Test")
	collection.DropCollection()

	u := &user{
		UserID:   bson.NewObjectId(),
		Username: "admin",
		Password: "admin",
		Roles:    []string{"r_user", "r_admin"},
	}
	SaveEntity("Test", u.UserID, u)

	// Testing process
	err := DeleteEntityWithCriteria("", map[string]interface{}{"username": "admin"})
	if err.Error() != "Invalid table name." {
		t.Errorf("Expected \"%s\"  but found \"%s\".", "Invalid table name.", err.Error())
	}

	err = DeleteEntityWithCriteria("Test", map[string]interface{}{"username": "admin"})
	if err != nil {
		t.Errorf("Expected no error but found none \"%s\".", err.Error())
	}

	var results []user
	err = collection.Find(nil).All(&results)

	if len(results) != 0 {
		t.Errorf("Expected 0 but found %d.", len(results))
	}
}
