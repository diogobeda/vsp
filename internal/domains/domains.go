package domains

import (
	"github.com/diogobeda/vsp/internal/domains/channel"
	"github.com/diogobeda/vsp/internal/domains/video"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

func InitChannelDomain(router *mux.Router, db *mgo.Database) {
	channelHandler := channel.CreateChannelHandler(db)
	router.HandleFunc("/channels", channelHandler.GetChannels).Methods("GET")
	router.HandleFunc("/channels", channelHandler.CraetChannel).Methods("POST")
	router.HandleFunc("/channels", channelHandler.UpdateChannel).Methods("PUT")
	router.HandleFunc("/channels/{id}", channelHandler.RemoveChannel).Methods("DELETE")
	router.HandleFunc("/channels/{handle}", channelHandler.GetChannel).Methods("GET")
}

func InitVideoDomain(router *mux.Router, db *mgo.Database, logger loggo.Logger) {
	videoHandler := video.CreaateVideoHandler(db, logger)
	router.HandleFunc("/videos", videoHandler.CreateVideo).Methods("POST")
	router.HandleFunc("/videos/upload", videoHandler.UploadVideo).Methods("POST")
}
