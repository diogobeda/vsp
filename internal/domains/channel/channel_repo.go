package channel

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Channel struct {
	ID        bson.ObjectId `bson:"_id"`
	Title     string
	URLHandle string
	CreatedAt string
	UpdatedAt string
}

type ChannelRepo struct {
	collection *mgo.Collection
}

func (repo *ChannelRepo) GetChannels() ([]Channel, error) {
	var result []Channel
	query := bson.M{}

	err := repo.collection.Find(query).All(&result)

	return result, err
}

func (repo *ChannelRepo) CreateChannel(channel Channel) error {
	channel.CreatedAt = time.Now().String()
	channel.UpdatedAt = time.Now().String()
	err := repo.collection.Insert(channel)

	return err
}

func (repo *ChannelRepo) GetChannelByHandle(handle string) (Channel, error) {
	var result Channel
	err := repo.collection.Find(bson.M{"urlhandle": handle}).One(&result)

	return result, err
}

func (repo *ChannelRepo) RemoveChannel(id string) error {
	objectid := bson.ObjectIdHex(id)
	err := repo.collection.RemoveId(objectid)
	return err
}

func (repo *ChannelRepo) UpdateChannel(channel Channel) error {
	channel.UpdatedAt = time.Now().String()
	err := repo.collection.UpdateId(channel.ID, channel)
	return err
}

func CreateChannelRepo(db *mgo.Database) *ChannelRepo {
	repo := &ChannelRepo{
		collection: db.C("channels"),
	}

	return repo
}
