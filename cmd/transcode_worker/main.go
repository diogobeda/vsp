package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diogobeda/vsp/internal/processors"
	"github.com/juju/loggo"
	"github.com/nsqio/go-nsq"
)

func main() {
	logger := loggo.GetLogger("transcode_worker")
	loggo.ConfigureLoggers(`transcode_worker=DEBUG`)

	logger.Infof("Creating nsq consumer")
	nsqConfig := nsq.NewConfig()
	consumer, consumerErr := nsq.NewConsumer("transcode", "transcode", nsqConfig)

	if consumerErr != nil {
		logger.Errorf(consumerErr.Error())
		log.Fatal()
	}

	logger.Infof("Adding TranscodeProcessor handler")
	consumer.ChangeMaxInFlight(200)
	consumer.AddConcurrentHandlers(processors.NewTranscodeProcessor(logger), 20)

	logger.Infof("Connecting to nsqdlookup")
	connectionErr := consumer.ConnectToNSQLookupd(os.Getenv("NSQ_LOOKUPD_URL"))

	if connectionErr != nil {
		logger.Errorf(connectionErr.Error())
		log.Fatal()
	}

	shutdownChan := make(chan os.Signal, 2)
	signal.Notify(shutdownChan, syscall.SIGINT)

	for {
		select {
		case <-consumer.StopChan:
			return
		case <-shutdownChan:
			consumer.Stop()
		}
	}
}
