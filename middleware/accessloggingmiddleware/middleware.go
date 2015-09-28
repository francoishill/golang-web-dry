package accessloggingmiddleware

// Inspiration from:
//  - https://github.com/codegangsta/negroni/blob/master/logger.go
//  - Beego's context/input: https://github.com/astaxie/beego/blob/8f7246e17b504c858592a28a87e96fa7537a5aaf/context/input.go

import (
	"github.com/codegangsta/negroni"
	"net/http"
	"time"

	"github.com/francoishill/golang-web-dry/requests/requestproxyutils"
)

type middleware struct {
	handler AccessInfoHandler
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accessInfo := &AccessInfo{
		HttpMethod: r.Method,
		RemoteAddr: r.RemoteAddr,
		RemoteIP:   requestproxyutils.IP(r),
		RequestURI: r.RequestURI,
		Proxies:    requestproxyutils.Proxy(r),
		UserAgent:  r.UserAgent(),
	}

	if m.handler != nil {
		startTime := time.Now()
		m.handler.OnStart(&StartAccessInfo{accessInfo, startTime})

		defer func() {
			r := recover()

			endTime := time.Now()
			duration := endTime.Sub(startTime)

			res := w.(negroni.ResponseWriter)

			m.handler.OnEnd(&EndAccessInfo{accessInfo, r != nil, res.Status(), http.StatusText(res.Status()), endTime, duration})

			if r != nil {
				panic(r)
			}
		}()
	}

	next(w, r)
}

func NewAccessLoggingMiddleware(handler AccessInfoHandler) *middleware {
	return &middleware{handler}
}
