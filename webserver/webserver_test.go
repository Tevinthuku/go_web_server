package webserver_test

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"web_server/webserver"
)

func TestWebServer(t *testing.T) {

	addr := ":8080"

	wb := webserver.NewWebServer("../www")
	err := wb.Run(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer wb.Close()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{name: "root", path: "/", expectedStatus: http.StatusOK},
		{name: "index", path: "/index.html", expectedStatus: http.StatusOK},
		{name: "about", path: "/about.html", expectedStatus: http.StatusOK},
		{name: "nonexistent", path: "/nonexistent.html", expectedStatus: http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			conn.Write([]byte(fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\n\r\n", test.path, addr)))

			response := make([]byte, 1024)
			_, err = conn.Read(response)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(response), fmt.Sprintf("%d %s", test.expectedStatus, http.StatusText(test.expectedStatus))) {
				t.Errorf("Expected %d %s, got %s", test.expectedStatus, http.StatusText(test.expectedStatus), string(response))
			}
		})
	}
}
