package proxy

import (
	"io"
	"log"
	"net/http"
)

func (p *Proxy) handlerHTTP(w http.ResponseWriter, r *http.Request) {
	for _, req := range p.req_handlers {
		resp, msg := req(r)
		log.Println(msg)
		if resp == nil {
			return
		} else {
			defer resp.Body.Close()
			copyHeader(w.Header(), resp.Header)
			w.WriteHeader(resp.StatusCode)
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				log.Printf("Error copying response body: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Printf("Error in roundtrip: %v", err)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
