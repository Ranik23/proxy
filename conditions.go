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

func stripScheme(s string) string {
	if strings.HasPrefix(s, "http://") {
		return strings.TrimPrefix(s, "http://")
	}
	if strings.HasPrefix(s, "https://") {
		return strings.TrimPrefix(s, "https://")
	}
	return s
}

func stripPort(s string) string {
	if i := strings.IndexByte(s, ':'); i != -1 {
		return s[:i]
	}
	return s
}

func isLoopback(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
