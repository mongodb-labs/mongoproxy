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

func SetConfigSaveLocation(c *ConfigLocation) {
	configSaveLocation = c
}

func updateConfiguration(config bson.M) (*mgo.ChangeInfo, error) {

	if configSaveLocation == nil {
		return nil, fmt.Errorf("No configuration save location.")
	}
	c := configSaveLocation.Session.DB(configSaveLocation.Database).C(configSaveLocation.Collection)
	changeInfo, err := c.Update(bson.M{"modules.name": "bi"},
		bson.M{"$set": bson.M{"modules.$.config": config}},
	)

}
