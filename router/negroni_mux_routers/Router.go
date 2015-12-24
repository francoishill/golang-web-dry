package negroni_mux_routers

import (
	"net/http"
)

type Router struct {
	urlParts    []string
	middlewares []http.HandlerFunc
	controller  Controller
	subRouters  []*Router
}
