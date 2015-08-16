package negroniprefixrouter

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func newNegroniMiddleware(funcs ...HandlerFunc) *negroni.Negroni {
	negroniHandlers := []negroni.Handler{}
	for ind, _ := range funcs {
		negroniHandlers = append(negroniHandlers, negroni.HandlerFunc(funcs[ind]))
	}
	return negroni.New(negroniHandlers...)
}
