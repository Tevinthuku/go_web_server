package webserver_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"web_server/webserver"
)

func TestWebServer(t *testing.T) {

	addr := ":8080"
	wb := setupTestServer(t, addr)
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

func setupTestServer(t *testing.T, addr string) *webserver.WebServer {
	wb := webserver.NewWebServer("../www")
	testRouteResistration(wb)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	// we set the listener here to avoid running the tests before the server is ready
	wb.Listener = listener
	go wb.Run(addr)

	return wb
}

func testRouteResistration(ws *webserver.WebServer) {
	ws.Get("/", func(w io.Writer, r *webserver.Request) {
		response := webserver.NewResponse(200, []byte("Hello, World!"))
		response.WriteTo(w)
	})
	ws.Get("/index.html", func(w io.Writer, r *webserver.Request) {
		response := webserver.NewResponse(200, []byte("Mock index.html!"))
		response.WriteTo(w)
	})
	ws.Get("/about.html", func(w io.Writer, r *webserver.Request) {
		response := webserver.NewResponse(200, []byte("Mock about.html!"))
		response.WriteTo(w)
	})
}
