package db

import (
	"os"

	"github.com/globalsign/mgo"
)

func CreateDbConnection(dbSession *mgo.Session) *mgo.Database {
	return dbSession.DB(os.Getenv("DB_NAME"))
}
