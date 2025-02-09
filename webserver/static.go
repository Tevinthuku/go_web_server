package webserver

import (
	"io"
	"log"
	"net/http"
)

func staticFsHandler(dir string) WebServerHandler {
	fs := http.Dir(dir)

	return func(w io.Writer, r *Request) {
		filepath := r.UrlValues["filepath"]
		if filepath == "" {
			filepath = "index.html"
		}

		file, err := fs.Open(filepath)
		if err != nil {
			response := NewResponse(404, []byte("Not Found"))
			_, err := response.WriteTo(w)
			if err != nil {
				log.Println("Error writing response:", err)
			}
			return
		}
		defer file.Close()
		file_bytes, err := io.ReadAll(file)
		if err != nil {
			response := NewResponse(500, []byte("Internal Server Error"))
			_, err := response.WriteTo(w)
			if err != nil {
				log.Println("Error writing response:", err)
			}
			return
		}
		response := NewResponse(http.StatusOK, file_bytes)
		_, err = response.WriteTo(w)
		if err != nil {
			log.Println("Error writing response:", err)
		}
	}
}
