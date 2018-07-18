package processors

import (
	"github.com/nareix/joy4/av/transcode"
	"bytes"

	"github.com/diogobeda/vsp/internal/messaging"
	"github.com/juju/loggo"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/cgo/ffmpeg"
	"github.com/nsqio/go-nsq"
)

type TranscodeProcessor struct {
	logger loggo.Logger
}

func (p *TranscodeProcessor) HandleMessage(message *nsq.Message) error {
	p.logger.Infof("[TranscodeProcessor] Decoding nsq message")
	tv := &TranscodableVideo{}
	decodeErr := messaging.DecodeMessage(message, tv)
	if decodeErr != nil {
		p.logger.Infof("[TranscodeProcessor] Error decoding nsq message")
		p.logger.Errorf(decodeErr.Error())
		return decodeErr
	}

	videoFile := bytes.NewReader(tv.Content)
	transcoder := &transcode.Demuxer{
		Options: transcode.Options{
			FindAudioDecoderEncoder: p.FindAudioCodec
		}
	}

	return nil
}

func (p *TranscodeProcessor) FindAudioCodec(stream av.AudioCodecData) (bool, av.AudioDecoder, av.AudioEncoder, error) {
	p.logger.Infof("[TranscodeProcessor] Creating audio decoder")
	dec, decodeErr := ffmpeg.NewAudioDecoder(stream)
	if decodeErr != nil {
		p.logger.Infof("[TranscodeProcessor] Error creating audio decoder")
		p.logger.Errorf(decodeErr.Error())
		return true, nil, nil, decodeErr
	}

	p.logger.Infof("[TranscodeProcessor] Creating audio encoder")
	enc, encodeErr := ffmpeg.NewAudioEncoderByName("libfdk_aac")
	if encodeErr != nil {
		p.logger.Infof("[TranscodeProcessor] Error creating audio encoder")
		p.logger.Errorf(encodeErr.Error())
		return true, nil, nil, encodeErr
	}

	enc.SetSampleRate(stream.SampleRate())
	enc.SetChannelLayout(av.CH_STEREO)
	enc.SetBitrate(12000)
	enc.SetOption("profile", "HE-AACv2")

	return true, dec, enc, nil
}

func NewTranscodeProcessor(logger loggo.Logger) *TranscodeProcessor {
	tp := &TranscodeProcessor{
		logger: logger,
	}

	return tp
}
