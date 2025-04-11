package proxy



import (
	"net"
	"net/http"
	"strings"
)

var IsLocalHost FuncRequestCondititon = func(r *http.Request) bool {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}

	if isLoopback(ip) {
		return true
	}

	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		for _, forwardedIP := range strings.Split(forwarded, ",") {
			forwardedIP = strings.TrimSpace(forwardedIP)
			if isLoopback(forwardedIP) {
				return true
			}
		}
	}

	return false
}

func isLoopback(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
