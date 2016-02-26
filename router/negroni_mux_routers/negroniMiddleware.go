package negroni_mux_routers

import (
	"github.com/codegangsta/negroni"
	"net/http"
)

type NegroniHandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func newNegroniMiddleware(middleWareFuncs []http.HandlerFunc, controllerMethod http.HandlerFunc) *negroni.Negroni {
	negroniHandlers := []negroni.Handler{}
	for ind, _ := range middleWareFuncs {
		negroniHandlers = append(negroniHandlers, negroni.Wrap(middleWareFuncs[ind]))
	}
	negroniHandlers = append(negroniHandlers, negroni.Wrap(http.HandlerFunc(controllerMethod)))
	return negroni.New(negroniHandlers...)
}
