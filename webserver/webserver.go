package webserver

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type WebServer struct {
	listener net.Listener
	rn       *routingNode
}

func NewWebServer() *WebServer {
	return &WebServer{rn: NewRoutingNode()}
}

func (ws *WebServer) Run(addr string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ws.listener = listener

	ws.start()
	return nil

}

func (ws *WebServer) start() {
	for {
		conn, err := ws.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Listener closed")
				break
			}
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
		if errors.Is(err, io.EOF) {
			log.Println("Client closed connection")
			return
		}
		log.Println("Error reading request line:", err)
		response := NewResponse(http.StatusBadRequest, []byte("Bad Request"))
		_, err := response.WriteTo(conn)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}
	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) != 3 {
		log.Println("Invalid request line:", requestLine)
		response := NewResponse(http.StatusBadRequest, []byte("Bad Request"))
		_, err := response.WriteTo(conn)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}
	// the METHOD is the first part and the PATH is the second part of the request line
	// request structure: GET /path HTTP/1.1
	method, rawPath := requestLineParts[0], requestLineParts[1]
	handler, err := ws.rn.MatchMethodAndPath(method, rawPath)
	if err != nil {
		response := NewResponse(http.StatusNotFound, []byte("Not Found"))
		_, err := response.WriteTo(conn)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}
	req := Request{
		Method:    method,
		Path:      rawPath,
		UrlValues: handler.DynamicContent,
	}
	handler.Handler(conn, &req)

}
