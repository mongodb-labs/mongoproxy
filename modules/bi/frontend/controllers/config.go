package controllers

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ConfigLocation struct {
	Session    *mgo.Session
	Database   string
	Collection string
}

var configSaveLocation *ConfigLocation

// SetConfigSaveLocation takes a driver session, a database and a collection to determine where the
// configuration of the BI module is stored.
func SetConfigSaveLocation(c *ConfigLocation) {
	configSaveLocation = c
}

// helper function to update the BI module's configuration. Assumes that the frontend
// was started with a configuration in a mongod instance, not a file.
func updateConfiguration(config bson.M) error {

	sessionCopy := configSaveLocation.Session.Copy()
	defer sessionCopy.Close()
	if configSaveLocation == nil {
		return fmt.Errorf("No configuration save location.")
	}
	c := sessionCopy.DB(configSaveLocation.Database).C(configSaveLocation.Collection)
	return c.Update(bson.M{"modules.name": "bi"},
		bson.M{"$set": bson.M{"modules.$.config": config}},
	)

}
