package processors

import (
	"bytes"

	"github.com/juju/loggo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/diogobeda/vsp/internal/domains/video"
	"github.com/diogobeda/vsp/internal/messaging"
	"github.com/nsqio/go-nsq"
)

type TranscodableVideo struct {
	FileName string
	Content  []byte
}

type UploadProcessor struct {
	s3Uploader  *s3manager.Uploader
	nsqProducer *nsq.Producer
	logger      loggo.Logger
}

func (p UploadProcessor) HandleMessage(message *nsq.Message) error {
	p.logger.Infof("[UploadProcessor] Decoding received message")
	uv := &video.UploadableVideo{}

	decodeErr := messaging.DecodeMessage(message, uv)
	if decodeErr != nil {
		p.logger.Infof("[UploadProcessor] Error decoding received message")
		p.logger.Errorf(decodeErr.Error())
		return decodeErr
	}

	p.logger.Infof("[UploadProcessor] Uploading file to s3 bucket")
	_, uploadErr := p.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(uv.Bucket),
		Key:    aws.String(uv.FileName),
		Body:   bytes.NewReader(uv.Content),
	})

	if uploadErr != nil {
		p.logger.Infof("[UploadProcessor] Error uploading file to s3 bucket")
		p.logger.Errorf(uploadErr.Error())
		return uploadErr
	}

	p.logger.Infof("[UploadProcessor] Publishing transcode command")
	tv := &TranscodableVideo{
		FileName: uv.FileName,
		Content:  uv.Content,
	}
	tsMessage, encodeErr := messaging.EncodeMessage(tv)
	if encodeErr != nil {
		p.logger.Infof("[UploadProcessor] Error encoding TranscodableVideo message")
		p.logger.Errorf(encodeErr.Error())
		return encodeErr
	}

	p.nsqProducer.Publish("transcode", tsMessage)

	return nil
}

func NewUploadProcessor(nsqProducer *nsq.Producer, logger loggo.Logger) *UploadProcessor {
	awsSession := session.Must(session.NewSession())
	s3Uploader := s3manager.NewUploader(awsSession)

	p := &UploadProcessor{
		s3Uploader:  s3Uploader,
		logger:      logger,
		nsqProducer: nsqProducer,
	}
	return p
}
