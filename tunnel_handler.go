package proxy

import (
	"io"
	"net"
	"net/http"
	"time"
)

func closeConns(client net.Conn, dest net.Conn) {
	client.Close()
	dest.Close()
}

func (p *Proxy) tunnelHandler(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "no hijack found", 500)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for _, req := range p.req_handlers {
		resp, msg := req(r)
		if resp != nil {
			_, err := clientConn.Write([]byte("HTTP/1.1 403 Forbidden\r\nX-Proxy-Error: " + msg + "\r\n\r\n"))
			if err != nil {
				return
			}
			closeConns(clientConn, destConn)
			return
		}
	}

	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		closeConns(clientConn, destConn)
		return
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}


func transfer(dest net.Conn, client net.Conn) {
	defer closeConns(dest, client)
	io.Copy(client, dest)
}