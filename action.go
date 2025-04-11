package proxy

import "net/http"



var AlwaysReject FuncReqHandler = func(r *http.Request) (*http.Response, string) {
	return nil, "always reject"
}