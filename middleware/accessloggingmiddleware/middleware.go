package accessloggingmiddleware

// Inspiration from:
//  - https://github.com/codegangsta/negroni/blob/master/logger.go
//  - Beego's context/input: https://github.com/astaxie/beego/blob/8f7246e17b504c858592a28a87e96fa7537a5aaf/context/input.go

import (
	"net/http"

	"github.com/francoishill/golang-web-dry/requests/requestproxyutils"
)

type AccessInfo struct {
	HttpMethod string
	RemoteAddr string
	RemoteIP   string
	RequestURI string
	Proxies    []string
	UserAgent  string
}

type middleware struct {
	withAccessInfo func(info *AccessInfo)
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if m.withAccessInfo != nil {
		m.withAccessInfo(&AccessInfo{
			HttpMethod: r.Method,
			RemoteAddr: r.RemoteAddr,
			RemoteIP:   requestproxyutils.IP(r),
			RequestURI: r.RequestURI,
			Proxies:    requestproxyutils.Proxy(r),
			UserAgent:  r.UserAgent(),
		})
	}

	next(w, r)
}

func NewAccessLoggingMiddleware(withAccessInfo func(info *AccessInfo)) *middleware {
	return &middleware{withAccessInfo}
}
