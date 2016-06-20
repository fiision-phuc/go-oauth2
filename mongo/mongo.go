package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

// Shared mongo session
var (
	config  *Config
	session *mgo.Session
)

// ConnectMongo creates session.
func ConnectMongo() {
	config = LoadConfigs()

	if session == nil {
		dialInfo := &mgo.DialInfo{
			Addrs:    config.Addresses,
			Timeout:  5 * time.Second,
			Database: config.Database,
			Username: config.Username,
			Password: config.Password,
		}

		// Create a session which maintains a pool of socket connections
		var err error
		if session, err = mgo.DialWithInfo(dialInfo); err != nil {
			panic(err)
		}
	}
}

// GetEventualSession clones session with eventual mode.
func GetEventualSession() (*mgo.Session, *mgo.Database) {
	/* Condition validation */
	if session == nil {
		return nil, nil
	}

	// Apply mode (http://godoc.org/labix.org/v2/mgo#Session.SetMode)
	clone := session.Clone()
	clone.SetMode(mgo.Eventual, true)
	return clone, clone.DB(config.Database)
}

// GetMonotonicSession clones session with monotonic mode.
func GetMonotonicSession() (*mgo.Session, *mgo.Database) {
	/* Condition validation */
	if session == nil {
		return nil, nil
	}

	// Apply mode (http://godoc.org/labix.org/v2/mgo#Session.SetMode)
	clone := session.Clone()
	clone.SetMode(mgo.Monotonic, true)
	return clone, clone.DB(config.Database)
}

// GetStrongSession clones session with strong mode.
func GetStrongSession() (*mgo.Session, *mgo.Database) {
	/* Condition validation */
	if session == nil {
		return nil, nil
	}

	// Apply mode (http://godoc.org/labix.org/v2/mgo#Session.SetMode)
	clone := session.Clone()
	clone.SetMode(mgo.Strong, true)
	return clone, clone.DB(config.Database)
}
