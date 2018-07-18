package channel

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/diogobeda/vsp/internal/web"
	"github.com/globalsign/mgo"
)

type ChannelHandler struct {
	webHandler  web.WebHandler
	channelRepo *ChannelRepo
}

func (handler *ChannelHandler) GetChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := handler.channelRepo.GetChannels()

	if err != nil {
		handler.webHandler.Internal(w, "There was a problem fetching the channels")
	}

	json.NewEncoder(w).Encode(channels)
}

func (handler *ChannelHandler) CraetChannel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var channel Channel
	decodeErr := decoder.Decode(&channel)

	if decodeErr != nil {
		handler.webHandler.BadRequest(w, "There was a problem creating the channel")
		return
	}

	_, findByHandleErr := handler.channelRepo.GetChannelByHandle(channel.URLHandle)
	if findByHandleErr == nil {
		handler.webHandler.BadRequest(w, "A channel with that url already exists")
		return
	}

	creationErr := handler.channelRepo.CreateChannel(channel)
	if creationErr != nil {
		handler.webHandler.BadRequest(w, "There was a problem creating the channel")
		return
	}

	handler.webHandler.Created(w)
}

func (handler *ChannelHandler) GetChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channel, err := handler.channelRepo.GetChannelByHandle(vars["handle"])

	if err != nil {
		handler.webHandler.NotFound(w, "There isn't a channel with that url")
		return
	}

	handler.webHandler.Ok(w)
	json.NewEncoder(w).Encode(channel)
}

func (handler *ChannelHandler) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var channel Channel
	decodeErr := decoder.Decode(&channel)

	if decodeErr != nil {
		handler.webHandler.BadRequest(w, "There was a problem updating the channel")
		return
	}

	updateErr := handler.channelRepo.UpdateChannel(channel)

	if updateErr != nil {
		handler.webHandler.BadRequest(w, "There was a problem updating the channel")
		return
	}

	handler.webHandler.Ok(w)
}

func (handler *ChannelHandler) RemoveChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := handler.channelRepo.RemoveChannel(vars["id"])

	if err != nil {
		handler.webHandler.BadRequest(w, "There was a problem removing the channel")
		return
	}

	handler.webHandler.Ok(w)
}

func CreateChannelHandler(db *mgo.Database) *ChannelHandler {
	channelRepo := CreateChannelRepo(db)
	handler := &ChannelHandler{
		channelRepo: channelRepo,
	}

	return handler
}
