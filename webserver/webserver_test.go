package webserver_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
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
		{name: "root", path: "/static/", expectedStatus: http.StatusOK},
		{name: "index", path: "/static/index.html", expectedStatus: http.StatusOK},
		{name: "about", path: "/static/about.html", expectedStatus: http.StatusOK},
		{name: "nonexistent-file", path: "/static/nonexistent.html", expectedStatus: http.StatusNotFound},
		{name: "languages", path: "/languages", expectedStatus: http.StatusOK},
		{name: "nonexistent-url", path: "/nonexistent-url", expectedStatus: http.StatusNotFound},
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
	wb := webserver.NewWebServer()
	testRouteResistration(wb)

	ready := make(chan struct{})
	go func() {
		if err := wb.Run(addr); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()

	// Wait for server to be ready
	go func() {
		for i := 0; i < 5; i++ {
			conn, err := net.Dial("tcp", addr)
			if err == nil {
				conn.Close()
				close(ready)
				break
			}
			time.Sleep(time.Duration(i) * time.Millisecond)
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
		if _, err := response.WriteTo(w); err != nil {
			panic(fmt.Sprintf("Failed to write response: %v", err))
		}
	})

	ws.Get("/languages", func(w io.Writer, r *webserver.Request) {
		response := webserver.NewResponse(200, []byte("Languages"))
		if _, err := response.WriteTo(w); err != nil {
			panic(fmt.Sprintf("Failed to write response: %v", err))
		}
	})

	// runtime.Caller(0) returns the file path of the current source file (webserver_test.go)
	// This allows us to construct absolute paths relative to this file's location,
	// ensuring static files can be found regardless of where the tests are executed from.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current file path")
	}
	projectRoot := filepath.Join(filepath.Dir(filename), "..")
	staticFilePath := filepath.Join(projectRoot, "www")
	ws.Static("/static", staticFilePath)
}
