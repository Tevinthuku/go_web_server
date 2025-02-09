package main

import (
	"io"
	"log"
	"web_server/webserver"
)

func personHandler(w io.Writer, r *webserver.Request) {
	id := r.UrlValues["id"]
	response := webserver.NewResponse(200, []byte("Person with id: "+id))
	_, err := response.WriteTo(w)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func peopleHandler(w io.Writer, r *webserver.Request) {
	response := webserver.NewResponse(200, []byte("People"))
	_, err := response.WriteTo(w)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func RegisterRoutes(ws *webserver.WebServer) {
	ws.Get("/people/:id", personHandler)
	ws.Get("/people", peopleHandler)
	ws.Static("/static", "./www")
}
