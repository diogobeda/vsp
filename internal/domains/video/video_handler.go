package video

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/nsqio/go-nsq"

	"github.com/diogobeda/vsp/internal/messaging"
	"github.com/diogobeda/vsp/internal/web"
	"github.com/globalsign/mgo"
	"github.com/juju/loggo"
)

type UploadableVideo struct {
	FileName string
	Size     int64
	Content  []byte
	Bucket   string
}

type VideoHandler struct {
	webHandler  web.WebHandler
	videoRepo   *VideoRepo
	nsqProducer *nsq.Producer
	logger      loggo.Logger
}

func (handler *VideoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var video Video

	decodeErr := json.NewDecoder(r.Body).Decode(&video)
	if decodeErr != nil {
		handler.webHandler.BadRequest(w, "There was an error creating the video")
		return
	}

	createErr := handler.videoRepo.CraeteVideo(video)
	if createErr != nil {
		handler.webHandler.BadRequest(w, "There was an error creating the video")
		return
	}

	handler.webHandler.Created(w)
}

func (handler *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	handler.logger.Infof("[VideoHandler] Reading file from request")
	uploadableVideo, readErr := handler.ReadVideoFile(r)
	if readErr != nil {
		handler.logger.Infof("[VideoHandler] Error reading file from request")
		handler.logger.Errorf(readErr.Error())
		handler.webHandler.BadRequest(w, "There was an error reading the video file")
		return
	}

	message, encodeErr := messaging.EncodeMessage(uploadableVideo)
	if encodeErr != nil {
		handler.logger.Infof("[VideoHandler] Error encoding nsq message")
		handler.logger.Errorf(encodeErr.Error())
		handler.webHandler.BadRequest(w, "There was an error reading the video file")
		return
	}

	handler.nsqProducer.Publish("upload", message)
	handler.webHandler.Ok(w)
}

func (handler *VideoHandler) FileToBytes(file multipart.File) []byte {
	var buf bytes.Buffer
	io.Copy(&buf, file)
	return buf.Bytes()
}

func (handler *VideoHandler) ReadVideoFile(r *http.Request) (*UploadableVideo, error) {
	file, header, readErr := r.FormFile("video")
	if readErr != nil {
		return &UploadableVideo{}, readErr
	}

	fileBytes := handler.FileToBytes(file)
	// mimeType := http.DetectContentType(fileBytes)
	// if mimeType

	defer file.Close()

	uv := &UploadableVideo{
		FileName: header.Filename,
		Size:     header.Size,
		Content:  fileBytes,
		Bucket:   "vsp-originals",
	}
	return uv, nil
}

func CreaateVideoHandler(db *mgo.Database, logger loggo.Logger) *VideoHandler {
	nsqConfig := nsq.NewConfig()
	nsqProducer, _ := nsq.NewProducer(os.Getenv("NSQD_URL"), nsqConfig)
	videoRepo := CreateVideoRepo(db)

	handler := &VideoHandler{
		videoRepo:   videoRepo,
		nsqProducer: nsqProducer,
		logger:      logger,
	}

	return handler
}
