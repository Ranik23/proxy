package proxy

import (
	"log"
	"net/http"
)

type RequestCondition interface {
	HandleReq(r *http.Request) bool
}

type FuncRequestCondititon func(r *http.Request) bool

func (f FuncRequestCondititon) HandleReq(r *http.Request) bool {
	return f(r)
}

type RequestHandler interface {
	Handle(r *http.Request) (*http.Response, string)
}

type FuncReqHandler func(r *http.Request) (*http.Response, string)

func (f FuncReqHandler) Handle(r *http.Request) (*http.Response, string) {
	return f(r)
}

type ReqProxyConds struct {
	proxy *Proxy
	conds []RequestCondition
}

func (pcond *ReqProxyConds) Do(handler RequestHandler) {
	pcond.proxy.req_handlers = append(pcond.proxy.req_handlers,
		FuncReqHandler(func(r *http.Request) (*http.Response, string) {
			log.Printf("[proxy] Handling request for host: %s, path: %s", r.Host, r.URL.Path)
			for i, cond := range pcond.conds {
				if !cond.HandleReq(r) {
					log.Printf("[proxy] Condition #%d failed for request to %s%s", i + 1, r.Host, r.URL.Path)
					return nil, "Condition failed"
				}
				log.Printf("[proxy] Condition #%d passed", i+1)
			}
			log.Printf("[proxy] All conditions passed, invoking handler")
			return handler.Handle(r)
		}))
}
