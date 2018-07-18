package db

import (
	"github.com/globalsign/mgo"
)

func CreateDbConnection(dbSession *mgo.Session) *mgo.Database {
	return dbSession.DB("youtube")
}
