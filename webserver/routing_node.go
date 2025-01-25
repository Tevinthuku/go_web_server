package webserver

import (
	"fmt"
	"net/http"
	"strings"
)

type routingNode struct {
	path_segment string
	is_dynamic   bool
	// handlers is a map of HTTP methods to handlers
	handlers map[string]http.Handler
	children map[string]*routingNode
}

func NewRoutingNode() *routingNode {
	return newBlankRoutingNode("/", false)
}

func (ut *routingNode) AddPattern(method string, url string, handler http.Handler) {
	split_url := strings.Split(url, "/")
	current_node := ut
	for _, segment := range split_url {
		if _, ok := current_node.children[segment]; ok {
			current_node = current_node.children[segment]
		} else {
			is_dynamic := strings.HasPrefix(segment, ":")
			if is_dynamic {
				segment = segment[1:]
			}
			current_node.children[segment] = newBlankRoutingNode(segment, is_dynamic)
			current_node = current_node.children[segment]
		}
	}
	current_node.handlers[method] = handler
}

type handlerWithDynamicContent struct {
	Handler        http.Handler
	DynamicContent map[string]interface{}
}

func (ut *routingNode) MatchMethodAndPath(method, path string) (*handlerWithDynamicContent, error) {
	split_url := strings.Split(path, "/")
	dynamic_content := map[string]interface{}{}
	current_node := ut
	for _, segment := range split_url {
		node, ok := current_node.children[segment]
		if ok {
			current_node = node
			continue
		} else {
			found := false
			for child, node := range current_node.children {
				if node.is_dynamic {
					dynamic_content[child] = segment
					current_node = node
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("no handler found for path %s", path)
			}
		}
	}
	handler, ok := current_node.handlers[method]
	if !ok {
		return nil, fmt.Errorf("no handler found for method %s", method)
	}
	return &handlerWithDynamicContent{
		Handler:        handler,
		DynamicContent: dynamic_content,
	}, nil
}

func newBlankRoutingNode(path_segment string, is_dynamic bool) *routingNode {
	return &routingNode{path_segment: path_segment, is_dynamic: is_dynamic, handlers: make(map[string]http.Handler), children: make(map[string]*routingNode)}
}
