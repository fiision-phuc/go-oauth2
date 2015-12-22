package mongo

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// EntityWithID finds entity with ID.
func EntityWithID(tableName string, entityID bson.ObjectId, entity interface{}) error {
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C(tableName)
	err := collection.FindId(entityID).One(entity)

	return err
}

// EntityWithCriteria finds entity with creteria.
func EntityWithCriteria(tableName string, criterion map[string]interface{}, entity interface{}) error {
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C(tableName)
	err := collection.Find(criterion).One(entity)

	return err
}

// SaveEntity inserts/updates an entity.
func SaveEntity(tableName string, entityID bson.ObjectId, entity interface{}) error {
	session, database := GetMonotonicSession()
	defer session.Close()

	session.SetSafe(&mgo.Safe{})
	collection := database.C(tableName)

	info, err := collection.UpsertId(entityID, entity)
	fmt.Println(info)
	return err
}

// DeleteEntity deletes a record from collection.
func DeleteEntity(tableName string, entityID bson.ObjectId) {
	session, database := GetMonotonicSession()
	defer session.Close()

	collection := database.C(tableName)
	collection.RemoveId(entityID)
}
