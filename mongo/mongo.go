package mongo

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
)

type Config struct {
	Host string
	Port string

	Database string
	Username string
	Password string
}

// Shared mongo session
var config *Config = nil
var session *mgo.Session = nil

/** Create session. */
func ConnectMongo(c Config) {
	if len(c.Host) == 0 || len(c.Port) == 0 || len(c.Database) == 0 {
		return
	}

	if session == nil {
		var err error

		dialInfo := &mgo.DialInfo{
			Addrs:    []string{fmt.Sprintf("%s:%s", c.Host, c.Port)},
			Timeout:  5 * time.Second,
			Database: c.Database,
			Username: c.Username,
			Password: c.Password,
		}

		// Create a session which maintains a pool of socket connections
		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}

		// Keep config instance
		config = &c
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
