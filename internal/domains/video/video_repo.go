package video

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Media struct {
	ID        bson.ObjectId `bson:"_id"`
	URL       string
	Quality   string
	CreatedAt string
}

type Video struct {
	ID          bson.ObjectId `bson:"_id"`
	Title       string
	Description string
	Media       Media
	CreatedAt   string
}

type VideoRepo struct {
	collection *mgo.Collection
}

func (repo *VideoRepo) CraeteVideo(video Video) error {
	video.CreatedAt = time.Now().String()
	err := repo.collection.Insert(video)
	return err
}

func (repo *VideoRepo) RemoveVideo(id string) error {
	objectid := bson.ObjectIdHex(id)
	err := repo.collection.RemoveId(objectid)
	return err
}

func CreateVideoRepo(db *mgo.Database) *VideoRepo {
	repo := &VideoRepo{
		collection: db.C("videos"),
	}

	return repo
}
