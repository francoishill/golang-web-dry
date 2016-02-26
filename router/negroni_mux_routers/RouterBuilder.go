package negroni_mux_routers

import (
	"net/http"
)

type RouterBuilder interface {
	SetController(controller Controller) RouterBuilder
	AddUrlParts(urlParts ...string) RouterBuilder
	AddMiddlewares(middlewares ...http.HandlerFunc) RouterBuilder
	AddSubrouters(subRouters ...*Router) RouterBuilder
	Build() *Router
}

type routerBuilder struct {
	urlParts    []string
	middlewares []http.HandlerFunc
	controller  Controller
	subRouters  []*Router
}

func (r *routerBuilder) SetController(controller Controller) RouterBuilder {
	r.controller = controller
	r.urlParts = controller.RelativeURLPatterns()
	return r
}

func (r *routerBuilder) AddUrlParts(urlParts ...string) RouterBuilder {
	r.urlParts = append(r.urlParts, urlParts...)
	return r
}

func (r *routerBuilder) AddMiddlewares(middlewares ...http.HandlerFunc) RouterBuilder {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *routerBuilder) AddSubrouters(subRouters ...*Router) RouterBuilder {
	r.subRouters = append(r.subRouters, subRouters...)
	return r
}

func (r *routerBuilder) Build() *Router {
	if r.controller == nil && len(r.subRouters) == 0 {
		panic("Cannot build Router if both 'controller' and 'subRouters' are defined.")
	}

	return &Router{
		r.urlParts,
		r.middlewares,
		r.controller,
		r.subRouters,
	}
}

func NewRouterBuilder() RouterBuilder {
	return &routerBuilder{
		[]string{},
		[]http.HandlerFunc{},
		nil,
		[]*Router{},
	}
}
