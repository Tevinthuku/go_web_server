package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"web_server/webserver"
)

func homePageHandler(w io.Writer, r *webserver.Request) {
	file, err := getFile("/index.html")
	if err != nil {
		response := webserver.NewResponse(404, []byte("Not Found"))
		response.WriteTo(w)
		return
	}
	response := webserver.NewResponse(200, file)
	response.WriteTo(w)
}

func aboutPageHandler(w io.Writer, r *webserver.Request) {
	file, err := getFile("/about.html")
	if err != nil {
		response := webserver.NewResponse(404, []byte("Not Found"))
		response.WriteTo(w)
		return
	}
	response := webserver.NewResponse(200, file)
	response.WriteTo(w)
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
