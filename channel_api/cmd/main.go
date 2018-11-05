package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/diogobeda/vsp/db"
	"github.com/diogobeda/vsp/internal/domains"
	"github.com/diogobeda/vsp/internal/web"
)

func main() {
	dbSession := db.CreateSession()
	database := db.CreateDbConnection(dbSession)
	router := web.CreateRouter()

	domains.InitChannelDomain(router, database)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println(fmt.Sprintf("Server starting in port %s", port))
	log.Fatal(http.ListenAndServe(port, router))
}
