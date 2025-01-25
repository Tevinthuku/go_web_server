package webserver

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type WebServer struct {
	listener net.Listener
	rootDir  string
	rn       *routingNode
}

func NewWebServer(rootDir string) *WebServer {
	return &WebServer{rootDir: rootDir, rn: NewRoutingNode()}
}

func (ws *WebServer) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ws.listener = listener
	go ws.start()
	return nil
}

func (ws *WebServer) start() {
	for {
		conn, err := ws.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go ws.handleConnection(conn)
	}
}

func (ws *WebServer) Close() error {
	return ws.listener.Close()
}

func (ws *WebServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	response := ws.handleRequest(conn)
	_, err := response.WriteTo(conn)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func (ws *WebServer) handleRequest(conn net.Conn) *Response {
	// Read the request line
	// GET /path HTTP/1.1
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading request line:", err)
		return NewResponse(http.StatusBadRequest, []byte("Bad Request"))
	}
	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) != 3 {
		log.Println("Invalid request line:", requestLine)
		return NewResponse(http.StatusBadRequest, []byte("Bad Request"))
	}
	// the path is the second part of the request line
	// request structure: GET /path HTTP/1.1
	rawPath := requestLineParts[1]
	path := filepath.Clean(rawPath)
	if path == "/" {
		path = "/index.html"
	}
	directory := filepath.Join(ws.rootDir, path)
	if !strings.HasPrefix(directory, filepath.Clean(ws.rootDir)) {
		return NewResponse(http.StatusForbidden, []byte("Forbidden path"))
	}
	body, err := os.ReadFile(directory)
	if err != nil {
		log.Println("Error reading file:", err)
		return NewResponse(http.StatusNotFound, []byte("File not found"))
	}
	return NewResponse(http.StatusOK, body)
}
