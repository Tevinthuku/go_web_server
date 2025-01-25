package webserver_test

import (
	"net/http"
	"testing"
	"web_server/webserver"

	"github.com/stretchr/testify/assert"
)

func TestRoutingNode(t *testing.T) {
	ut := webserver.NewRoutingNode()
	ut.AddPattern("/hello/world", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}), "GET")
	handler, err := ut.MatchMethodAndPath("GET", "/hello/world")
	assert.Nil(t, err)
	assert.NotNil(t, handler)

	handler, err = ut.MatchMethodAndPath("GET", "/hello/world2")
	assert.Nil(t, handler)
	assert.NotNil(t, err)

	handler, err = ut.MatchMethodAndPath("PATCH", "/hello/world/")
	assert.Nil(t, handler)
	assert.NotNil(t, err)
}
