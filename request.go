package proxy

import "net/http"


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
		for _, cond := range pcond.proxy.req_conds {
			if !cond.HandleReq(r) {
				return nil, "condition failed"
			}
		}
		return handler.Handle(r) 
	}))
} 



