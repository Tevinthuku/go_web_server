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
		response := NewResponse(http.StatusOK, nil)
		_, err = response.WriteHeaderTo(w)
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}
		// stream file content directly to the writer
		if _, err := io.Copy(w, file); err != nil {
			log.Println("Error writing response:", err)
			return
		}
	}
}
