package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

// Shared mongo session
var config *MongoConfig = nil
var session *mgo.Session = nil

// ConnectMongo creates session.
func ConnectMongo() {
	config = LoadMongoConfigs()

	if session == nil {
		var err error

		dialInfo := &mgo.DialInfo{
			Addrs:    config.Addresses,
			Timeout:  5 * time.Second,
			Database: config.Database,
			Username: config.Username,
			Password: config.Password,
		}

		// Create a session which maintains a pool of socket connections
		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
	}
}

/** Clone session with eventual mode. */
func GetEventualSession() (*mgo.Session, *mgo.Database) {
	/* Condition validation */
	if session == nil {
		return nil, nil
	} else {
		// Apply mode (http://godoc.org/labix.org/v2/mgo#Session.SetMode)
		clone := session.Clone()
		clone.SetMode(mgo.Eventual, true)
		return clone, clone.DB(config.Database)
	}
}

/** Clone session with monotonic mode. */
func GetMonotonicSession() (*mgo.Session, *mgo.Database) {
	/* Condition validation */
	if session == nil {
		return nil, nil
	} else {
		// Apply mode (http://godoc.org/labix.org/v2/mgo#Session.SetMode)
		clone := session.Clone()
		clone.SetMode(mgo.Monotonic, true)
		return clone, clone.DB(config.Database)
	}
}

/** Clone session with strong mode. */
func GetStrongSession() (*mgo.Session, *mgo.Database) {
	/* Condition validation */
	if session == nil {
		return nil, nil
	} else {
		clone := session.Clone()
		clone.SetMode(mgo.Strong, true)
		return clone, clone.DB(config.Database)
	}
}
