package proxy

import (
	"net"
	"strings"
)


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