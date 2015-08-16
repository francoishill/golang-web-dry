package negroniprefixrouter

import "github.com/gorilla/mux"

type routeHandler struct {
	path             string
	controller       HandlerFunc
	buildHandleRoute func(route *mux.Route)
}

type PrefixRouter struct {
	pathPrefix              string
	routeHandlers           []*routeHandler
	routeSpecificMiddleWare []HandlerFunc
	subRouters              []*PrefixRouter
}

func NewPrefixRouter(pathPrefix string) *PrefixRouter {
	return &PrefixRouter{pathPrefix: pathPrefix}
}

func (pr *PrefixRouter) AddRouteHandler(path string, controller HandlerFunc, buildHandleRoute func(route *mux.Route)) *PrefixRouter {
	pr.routeHandlers = append(pr.routeHandlers, &routeHandler{
		path:             path,
		controller:       controller,
		buildHandleRoute: buildHandleRoute,
	})
	return pr
}

func (pr *PrefixRouter) AddRouteMiddleWares(middleWares ...HandlerFunc) *PrefixRouter {
	pr.routeSpecificMiddleWare = append(pr.routeSpecificMiddleWare, middleWares...)
	return pr
}

func (pr *PrefixRouter) AddSubPrefixRouters(subPrefixRouters ...*PrefixRouter) *PrefixRouter {
	pr.subRouters = append(pr.subRouters, subPrefixRouters...)
	return pr
}

func (pr *PrefixRouter) BuildHandlers(router *mux.Router, parentMiddleWare []HandlerFunc) {
	prefixSubRouter := router.PathPrefix(pr.pathPrefix).Subrouter()

	for _, handler := range pr.routeHandlers {
		combinedMiddleWareHandlers := []HandlerFunc{}
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, parentMiddleWare...)
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, pr.routeSpecificMiddleWare...)
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, handler.controller)
		handler.buildHandleRoute(prefixSubRouter.Handle(handler.path, newNegroniMiddleware(combinedMiddleWareHandlers...)))
	}

	for _, subRouter := range pr.subRouters {
		combinedMiddleWareHandlers := []HandlerFunc{}
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, parentMiddleWare...)
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, pr.routeSpecificMiddleWare...)
		subRouter.BuildHandlers(prefixSubRouter, parentMiddleWare)
	}
}
