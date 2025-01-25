package webserver

import (
	"fmt"
	"net/http"
	"strings"
)

type RoutingNode struct {
	path_segment string
	// handlers is a map of HTTP methods to handlers
	handlers map[string]http.Handler
	children map[string]*RoutingNode
}

func NewRoutingNode() *RoutingNode {
	return &RoutingNode{path_segment: "/", children: make(map[string]*RoutingNode)}
}

func (ut *RoutingNode) AddPattern(url string, handler http.Handler, method string) {
	split_url := strings.Split(url, "/")
	fmt.Println(split_url)
	current_node := ut
	for _, segment := range split_url {
		if _, ok := current_node.children[segment]; ok {
			current_node = current_node.children[segment]
		} else {
			current_node.children[segment] = &RoutingNode{
				path_segment: segment,
				handlers:     make(map[string]http.Handler),
				children:     make(map[string]*RoutingNode),
			}
			current_node = current_node.children[segment]
		}
	}
	current_node.handlers[method] = handler
}

func (ut *RoutingNode) MatchMethodAndPath(method, path string) (http.Handler, error) {
	split_url := strings.Split(path, "/")
	current_node := ut
	for _, segment := range split_url {
		if _, ok := current_node.children[segment]; !ok {
			return nil, fmt.Errorf("no handler found for path %s", path)
		}
		current_node = current_node.children[segment]
	}
	handler, ok := current_node.handlers[method]
	if !ok {
		return nil, fmt.Errorf("no handler found for method %s", method)
	}
	return handler, nil
}
