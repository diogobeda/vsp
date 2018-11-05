package domains

import (
	"github.com/diogobeda/vsp/channel_api/internal/domains/channel"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
)

func InitChannelDomain(router *mux.Router, db *mgo.Database) {
	channelHandler := channel.CreateChannelHandler(db)
	router.HandleFunc("/channels", channelHandler.GetChannels).Methods("GET")
	router.HandleFunc("/channels", channelHandler.CraetChannel).Methods("POST")
	router.HandleFunc("/channels", channelHandler.UpdateChannel).Methods("PUT")
	router.HandleFunc("/channels/{id}", channelHandler.RemoveChannel).Methods("DELETE")
	router.HandleFunc("/channels/{handle}", channelHandler.GetChannel).Methods("GET")
}
