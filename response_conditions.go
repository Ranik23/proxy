package proxy

import "net/http"


func StatusIs(code int) FuncResponseCondition {
	return func(resp *http.Response) bool {
		return resp.StatusCode == code
	}
}

func HasHeader(key string) FuncResponseCondition {
	return func(resp *http.Response) bool {
		_, ok := resp.Header[key]
		return ok
	}
}
