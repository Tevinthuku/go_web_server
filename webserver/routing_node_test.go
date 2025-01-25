package webserver_test

import (
	"net/http"
	"testing"
	"web_server/webserver"

	"github.com/stretchr/testify/assert"
)

func TestStaticRoutingNode(t *testing.T) {
	ut := webserver.NewRoutingNode()
	ut.AddPattern("GET", "/hello/world", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	handler, err := ut.MatchMethodAndPath("GET", "/hello/world")
	assert.Nil(t, err)
	assert.NotNil(t, handler)

	handler, err = ut.MatchMethodAndPath("GET", "/hello/world2")
	assert.Nil(t, handler)
	assert.NotNil(t, err)

	handler, err = ut.MatchMethodAndPath("PATCH", "/hello/world/")
	assert.Nil(t, handler)
	assert.NotNil(t, err)

	ut.AddPattern("GET", "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, root!"))
	}))
	handler, err = ut.MatchMethodAndPath("GET", "/")
	assert.Nil(t, err)
	assert.NotNil(t, handler)
}
