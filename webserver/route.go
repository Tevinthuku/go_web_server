package webserver

import (
	"net/http"
)

func (ws *WebServer) Get(path string, handler WebServerHandler) {
	ws.rn.AddPattern(http.MethodGet, path, handler)
}

func (ws *WebServer) Post(path string, handler WebServerHandler) {
	ws.rn.AddPattern(http.MethodPost, path, handler)
}

func (ws *WebServer) Put(path string, handler WebServerHandler) {
	ws.rn.AddPattern(http.MethodPut, path, handler)
}

func (ws *WebServer) Delete(path string, handler WebServerHandler) {
	ws.rn.AddPattern(http.MethodDelete, path, handler)
}

func (ws *WebServer) Patch(path string, handler WebServerHandler) {
	ws.rn.AddPattern(http.MethodPatch, path, handler)
}

// relativePath is the base URL path prefix for serving static files (e.g., "/static")
// dir is the local filesystem directory path containing the static files to serve
func (ws *WebServer) Static(relativePath, dir string) {
	handler := staticFsHandler(dir)
	urlPath := relativePath + "/:filepath"
	ws.rn.AddPattern(http.MethodGet, urlPath, handler)
}
