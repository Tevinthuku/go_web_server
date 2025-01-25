package webserver_test

import (
	"net/http"
	"testing"
	"web_server/webserver"

	"github.com/stretchr/testify/assert"
)

func TestStaticRoutingNode(t *testing.T) {
	rn := webserver.NewRoutingNode()
	rn.AddPattern("GET", "/hello/world", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	handler, err := rn.MatchMethodAndPath("GET", "/hello/world")
	assert.Nil(t, err)
	assert.NotNil(t, handler)

	handler, err = rn.MatchMethodAndPath("GET", "/hello/world2")
	assert.Nil(t, handler)
	assert.NotNil(t, err)

	handler, err = rn.MatchMethodAndPath("PATCH", "/hello/world/")
	assert.Nil(t, handler)
	assert.NotNil(t, err)

	rn.AddPattern("GET", "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, root!"))
	}))
	handler, err = rn.MatchMethodAndPath("GET", "/")
	assert.Nil(t, err)
	assert.NotNil(t, handler)
}

func TestDynamicRoutingNode(t *testing.T) {
	rn := webserver.NewRoutingNode()

	rn.AddPattern("GET", "/people/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Specific person!"))
	}))

	rn.AddPattern("GET", "/people/list", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("People list!"))
	}))

	handler, err := rn.MatchMethodAndPath("GET", "/people/123")
	assert.Nil(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, handler.DynamicContent["id"], "123")

	handler, err = rn.MatchMethodAndPath("GET", "/people/list")
	assert.Nil(t, err)
	assert.NotNil(t, handler)

	rn.AddPattern("GET", "/people/:id/duplicate", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Duplicate person!"))
	}))

	handler, err = rn.MatchMethodAndPath("GET", "/people/123/duplicate")
	assert.Nil(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, handler.DynamicContent["id"], "123")
}
