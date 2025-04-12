package proxy

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
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
	message := "Access Reject"
	return &http.Response{
		Status: "403",
		StatusCode: http.StatusForbidden,
		Body: NewMyReader(message),
	}, "Access Reject"
}

func BlockUserAgent(agent string) FuncReqHandler {
	return func(r *http.Request) (*http.Response, string) {
		if strings.Contains(r.UserAgent(), agent) {
			msg := "Blocked User-Agent: " + agent
			return &http.Response{
				Status:     "403 Forbidden",
				StatusCode: http.StatusForbidden,
				Body:       NewMyReader(msg),
			}, msg
		}
		return &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       NewMyReader("User-Agent allowed."),
		}, "User-Agent check passed"
	}
}


func SimulateDelay(d time.Duration) FuncReqHandler {
	return func(r *http.Request) (*http.Response, string) {
		time.Sleep(d)
		msg := fmt.Sprintf("Response delayed by %s", d)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       NewMyReader(msg),
		}, msg
	}
}

var InspectRequest FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	var sb strings.Builder
	sb.WriteString("Inspecting request:\n")
	sb.WriteString("Host: " + r.Host + "\n")
	sb.WriteString("Method: " + r.Method + "\n")
	sb.WriteString("Path: " + r.URL.Path + "\n")
	sb.WriteString("Headers:\n")
	for k, v := range r.Header {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", ")))
	}

	log.Println(sb.String())

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       NewMyReader("Request inspected"),
	}, "Request inspected"
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
