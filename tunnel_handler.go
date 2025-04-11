package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func (p *Proxy) tunnelHandler(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10 * time.Second)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	hikacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "no hijack found", 500)
		return
	}

	client_conn, _, err := hikacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for _, req := range p.req_handlers {
		resp, msg := req(r)
		log.Println(msg)
		if resp == nil {
			return
		}
	}

	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(dest net.Conn, client net.Conn) {
	defer client.Close()
	defer dest.Close()
	io.Copy(client, dest)
}