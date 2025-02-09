package main

import (
	"log"
	"web_server/webserver"
)

func main() {

	wb := webserver.NewWebServer()
	defer wb.Close()

	RegisterRoutes(wb)
	err := wb.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
