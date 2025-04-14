package proxy

import "net/http"



type ResponseCondition interface {
	HandleResp(resp *http.Response) bool
}

type FuncResponseCondition func(resp *http.Response) bool

func (f FuncResponseCondition) HandleResp(resp *http.Response) bool {
	return f(resp)
}


type ResponseHandler interface {
	Handle(resp *http.Response) *http.Response
}
type FuncResponseHandler func(resp *http.Response) *http.Response

func (f FuncResponseHandler) Handle(resp *http.Response) *http.Response {
	return f(resp)
}


type RespProxyCond struct {
	proxy *Proxy
	conds []ResponseCondition
}

func (rcond *RespProxyCond) Do(handler FuncResponseHandler) {
	rcond.proxy.resp_handlers = append(rcond.proxy.resp_handlers, 
	FuncResponseHandler(func(resp *http.Response) *http.Response {
		for _, cond := range rcond.conds {
			if !cond.HandleResp(resp) {
				return nil
			}
		}
		return handler.Handle(resp)
	}))
} 
