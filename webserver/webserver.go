package webserver

import (
	"bufio"
	"log"
	"net"
	"net/http"
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

func (ws *WebServer) Run(addr string) (<-chan struct{}, error) {
	ready := make(chan struct{})
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return ready, err
	}
	ws.listener = listener
	go ws.start()
	// Signal that the server is ready
	close(ready)
	return ready, nil
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
	ws.handleRequest(conn)
}

func (ws *WebServer) handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	// Read the request line
	// GET /path HTTP/1.1
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading request line:", err)
		response := NewResponse(http.StatusBadRequest, []byte("Bad Request"))
		response.WriteTo(conn)
		return
	}
	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) != 3 {
		log.Println("Invalid request line:", requestLine)
		response := NewResponse(http.StatusBadRequest, []byte("Bad Request"))
		response.WriteTo(conn)
		return
	}
	// the METHOD is the first part and the PATH is the second part of the request line
	// request structure: GET /path HTTP/1.1
	method, rawPath := requestLineParts[0], requestLineParts[1]
	handler, err := ws.rn.MatchMethodAndPath(method, rawPath)
	if err != nil {
		response := NewResponse(http.StatusNotFound, []byte("Not Found"))
		response.WriteTo(conn)
		return
	}
	req := Request{
		Method:    method,
		Path:      rawPath,
		UrlValues: handler.DynamicContent,
	}
	handler.Handler(conn, &req)

}
