package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"web_server/webserver"
)

func homePageHandler(w io.Writer, r *webserver.Request) {
	serveHTMLFile(w, "/index.html")
}

func aboutPageHandler(w io.Writer, r *webserver.Request) {
	serveHTMLFile(w, "/about.html")
}

func personHandler(w io.Writer, r *webserver.Request) {
	id := r.UrlValues["id"]
	response := webserver.NewResponse(200, []byte("Person with id: "+id))
	response.WriteTo(w)
}

func peopleHandler(w io.Writer, r *webserver.Request) {
	response := webserver.NewResponse(200, []byte("People"))
	response.WriteTo(w)
}

func RegisterRoutes(ws *webserver.WebServer) {
	ws.Get("/", homePageHandler)
	ws.Get("/index.html", homePageHandler)
	ws.Get("/about.html", aboutPageHandler)
	ws.Get("/people/:id", personHandler)
	ws.Get("/people", peopleHandler)
}

func serveHTMLFile(w io.Writer, filepath string) {
	file, err := getFile(filepath)
	if err != nil {
		response := webserver.NewResponse(404, []byte("Not Found"))
		_, err := response.WriteTo(w)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}
	response := webserver.NewResponse(200, file)
	_, err = response.WriteTo(w)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func getFile(rawPath string) ([]byte, error) {
	path := filepath.Clean(rawPath)
	directory := filepath.Join("./www", path)
	if !strings.HasPrefix(directory, filepath.Clean("./www")) {
		return nil, errors.New("forbidden path")
	}
	body, err := os.ReadFile(directory)
	if err != nil {
		return nil, err
	}
	return body, nil
}
