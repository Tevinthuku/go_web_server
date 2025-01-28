package webserver_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
	"web_server/webserver"

	"github.com/stretchr/testify/assert"
)

func TestWebServer(t *testing.T) {

	addr := ":8080"
	wb := setupTestServer(t, addr)
	defer wb.Close()
	client := &http.Client{}

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
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost%s%s", addr, test.path), nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			assert.Equal(t, test.expectedStatus, resp.StatusCode)
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
	ready := make(chan struct{})
	go func() {
		if err := wb.Run(addr); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()

	// Wait for server to be ready
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			close(ready)
		}
	}()
	select {
	case <-ready:
	case <-time.After(2 * time.Second):
		t.Fatal("server failed to start within timeout")
	}

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
