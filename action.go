package proxy

import (
	"fmt"
	"net/http"
	"strings"
)


type MyReader struct {
	strings.Reader
}

func NewMyReader(s string) *MyReader {
	return &MyReader{
		Reader: *strings.NewReader(s),
	}
}

func (r *MyReader) Close() error {
	return nil
}


var AlwaysReject FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	return nil, "always reject"
}

var PersonalizedHello FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}
	message := fmt.Sprintf("Hello, %s!", name)
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       NewMyReader(message),
	}, "Personalized hello sent"
}


var JSONResponse FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	response := `{"message": "Hello, world", "status": "success"}`
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       NewMyReader(response),
	}, "JSON response sent"
}

var BadRequest FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	return &http.Response{
		Status:     "400 Bad Request",
		StatusCode: http.StatusBadRequest,
		Body:       NewMyReader("Bad request!"),
	}, "Bad request error"
}

var RedirectHandler FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	return &http.Response{
		Status:     "301 Moved Permanently",
		StatusCode: http.StatusMovedPermanently,
		Header:     map[string][]string{"Location": {"http://example.com"}},
		Body:       NewMyReader("redirect"),
	}, "Redirecting to example.com"
}
