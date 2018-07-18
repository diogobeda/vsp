package db

import (
	"log"
	"os"

	"github.com/globalsign/mgo"
)

func CreateSession() *mgo.Session {
	session, err := mgo.Dial(os.Getenv("MONGODB_URL"))

	if err != nil {
		log.Fatal(err)
	}

	return session
}
