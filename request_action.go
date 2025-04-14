package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)


var AlwaysReject FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	message := "Access Reject"
	return &http.Response{
		Status: "403",
		StatusCode: http.StatusForbidden,
		Body: io.NopCloser(strings.NewReader(message)),
	}, "Access Reject"
}

func BlockUserAgent(agent string) FuncReqHandler {
	return func(r *http.Request) (*http.Response, string) {
		if strings.Contains(r.UserAgent(), agent) {
			msg := "Blocked User-Agent: " + agent
			return &http.Response{
				Status:     "403 Forbidden",
				StatusCode: http.StatusForbidden,
				Body:       io.NopCloser(strings.NewReader(msg)),
			}, msg
		}
		return &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:      io.NopCloser(strings.NewReader("User-Agent allowed.")),
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
			Body:       io.NopCloser(strings.NewReader(msg)),
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
		Body:       io.NopCloser(strings.NewReader("Request inspected")),
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
		Body:       io.NopCloser(strings.NewReader(message)),
	}, "Personalized hello sent"
}
