package webserver

import (
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	StatusCode int
	Body       []byte
}

func NewResponse(statusCode int, body []byte) *Response {
	return &Response{StatusCode: statusCode, Body: body}
}

func (r *Response) WriteTo(w io.Writer) (int64, error) {
	responseAsBytes := []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", r.StatusCode, http.StatusText(r.StatusCode)))
	responseAsBytes = append(responseAsBytes, r.Body...)
	responseAsBytes = append(responseAsBytes, []byte("\r\n")...)
	n, err := w.Write(responseAsBytes)
	return int64(n), err
}
