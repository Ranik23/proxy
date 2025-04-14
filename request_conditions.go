package proxy

import (
	"log"
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


func Is(hostname string) FuncRequestCondititon {
	hostname = stripScheme(hostname)
	hostname = stripPort(hostname)

	return func(r *http.Request) bool {
		reqHost := stripScheme(r.Host)
		reqHost = stripPort(reqHost)

		match := reqHost == hostname

		log.Printf("[cond] Checking host: request=%q, expected=%q → match=%v", reqHost, hostname, match)

		return match
	}
}
