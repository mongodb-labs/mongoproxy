// Package bsonutil provides a utility function to retrieve values from bson documents.
package bsonutil

import (
	"gopkg.in/mgo.v2/bson"
)

// FindValueByKey returns the value of keyName in document. If keyName is not found
// in the top-level of the document, ErrNoSuchField is returned as the error.
func FindValueByKey(keyName string, document bson.D) interface{} {
	for _, key := range document {
		if key.Name == keyName {
			return key.Value
		}
	}
	return nil
}
