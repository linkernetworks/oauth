package mongo

import "gopkg.in/mgo.v2/bson"

// Selector sql wrapper for mongo
type Selector struct {
	Selector   bson.M
	Collection string
	Sort       string
	Limit      int
	Skip       int
}
