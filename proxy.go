package proxy

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
)

type Proxy struct {
	server		  			*http.Server


	req_handlers			[]FuncReqHandler
	resp_handlers			[]FuncResponseHandler


	req_conds				[]RequestCondition
	resp_conds				[]ResponseCondition


	cert 					string
	key 					string
	proto					string
}

func NewProxy(address string, cert string, key string, proto string) *Proxy {
	srv := &http.Server{
		Addr: address,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	return &Proxy{
		server: srv,
		cert: cert,
		key: key,
		proto: proto,
	}
}

func (p *Proxy) Run() error {
		
	p.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			p.tunnelHandler(w, r)
		} else {
			p.handlerHTTP(w, r)
		}
	})

	if p.proto == "http" {
		log.Println("Starting the server...")
		if err := p.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
	
	log.Println("Starting the server...")
	if err := p.server.ListenAndServeTLS(p.cert, p.key); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (p *Proxy) OnRequest(conds ...RequestCondition) *ReqProxyConds {
	return &ReqProxyConds{proxy: p, conds: conds}
}

func (p *Proxy) OnResponse(conds ...ResponseCondition) *RespProxyCond {
	return &RespProxyCond{proxy: p, conds: conds}
}

