package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/diogobeda/vsp/db"
	"github.com/diogobeda/vsp/internal/domains"
	"github.com/diogobeda/vsp/internal/web"
	"github.com/juju/loggo"
)

func main() {
	dbSession := db.CreateSession()
	database := db.CreateDbConnection(dbSession)
	logger := loggo.GetLogger("video_api")
	router := web.CreateRouter()

	loggo.ConfigureLoggers(`video_api=DEBUG`)

	domains.InitVideoDomain(router, database, logger)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println(fmt.Sprintf("Server starting in port %s", port))
	log.Fatal(http.ListenAndServe(port, router))
}
