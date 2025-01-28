package webserver

import (
	"fmt"
	"io"
	"strings"
)

type routingNode struct {
	pathSegment string
	isDynamic   bool
	// handlers is a map of HTTP methods to handlers
	handlers map[string]WebServerHandler
	children map[string]*routingNode
}

type WebServerHandler func(w io.Writer, r *Request)

type Request struct {
	Method    string
	Path      string
	UrlValues map[string]string
}

func NewRoutingNode() *routingNode {
	return newBlankRoutingNode("/", false)
}

func (ut *routingNode) AddPattern(method string, url string, handler WebServerHandler) {
	if method == "" || url == "" || handler == nil {
		panic("method, url and handler must not be empty")
	}
	splitURL := strings.Split(url, "/")
	currentNode := ut
	for _, segment := range splitURL {
		if _, ok := currentNode.children[segment]; ok {
			currentNode = currentNode.children[segment]
		} else {
			isDynamic := strings.HasPrefix(segment, ":")
			if isDynamic {
				segment = segment[1:]
			}
			currentNode.children[segment] = newBlankRoutingNode(segment, isDynamic)
			currentNode = currentNode.children[segment]
		}
	}
	currentNode.handlers[method] = handler
}

type handlerWithDynamicContent struct {
	Handler        WebServerHandler
	DynamicContent map[string]string
}

func (ut *routingNode) MatchMethodAndPath(method, path string) (*handlerWithDynamicContent, error) {
	if method == "" || path == "" {
		return nil, fmt.Errorf("method and path must not be empty")
	}
	split_url := strings.Split(path, "/")
	dynamicContent := make(map[string]string)
	currentNode := ut
	for _, segment := range split_url {
		node, ok := currentNode.children[segment]
		if ok {
			currentNode = node
			continue
		} else {
			currentNode = currentNode.getChildDynamicRoutingNode()
			if currentNode == nil {
				return nil, fmt.Errorf("no handler found for path %s", path)
			} else {
				dynamicContent[currentNode.pathSegment] = segment
			}
		}
	}
	handler, ok := currentNode.handlers[method]
	if !ok {
		return nil, fmt.Errorf("no handler found for method %s", method)
	}
	return &handlerWithDynamicContent{
		Handler:        handler,
		DynamicContent: dynamicContent,
	}, nil
}

func (ut *routingNode) getChildDynamicRoutingNode() *routingNode {
	for _, child := range ut.children {
		if child.isDynamic {
			return child
		}
	}
	return nil
}

func newBlankRoutingNode(pathSegment string, isDynamic bool) *routingNode {
	return &routingNode{pathSegment: pathSegment, isDynamic: isDynamic, handlers: make(map[string]WebServerHandler), children: make(map[string]*routingNode)}
}
