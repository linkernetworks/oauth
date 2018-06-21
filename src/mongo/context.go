package mongo

import mgo "gopkg.in/mgo.v2"

//Context :
type Context struct {
	client   *MongoClient
	session  *mgo.Session
	database *mgo.Database
}

func (c *Context) Close() {
	c.session.Close()
}

// InsertOne insert one doc to mongo
func (c *Context) InsertOne(sql Selector, doc interface{}) error {
	col := c.database.C(sql.Collection)
	err := col.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

// QueryOne query one doc
func (c *Context) QueryOne(sql Selector, doc interface{}) error {
	col := c.database.C(sql.Collection)
	err := col.Find(sql.Selector).One(doc)
	return err
}

// CheckExist check if doc exist in mongo
func (c *Context) CheckExist(sql Selector, doc interface{}) (bool, error) {
	err := c.QueryOne(sql, doc)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// UpdateOne update doc by Selector
func (c *Context) UpdateOne(sql Selector, doc interface{}) error {
	col := c.database.C(sql.Collection)
	return col.Update(sql.Selector, doc)
}

// Upsert Upsert doc by Selector
func (c *Context) Upsert(sql Selector, doc interface{}) error {
	col := c.database.C(sql.Collection)
	_, err := col.Upsert(sql.Selector, doc)
	return err
}

// RemoveOne remove doc by Selector
func (c *Context) RemoveOne(sql Selector) error {

	col := c.database.C(sql.Collection)
	return col.Remove(sql.Selector)
}

// DropDatabase drop mongo database
func (c *Context) DropDatabase(database string) error {
	return c.database.DropDatabase()
}
