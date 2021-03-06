package negroni_mux_routers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func setupRouters(router *mux.Router, parentMiddleWare []http.HandlerFunc, routers []*Router) {
	if len(routers) == 0 {
		return
	}

	for _, rd := range routers {
		combinedMiddleWareHandlers := []http.HandlerFunc{}
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, parentMiddleWare...)
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, rd.middlewares...)

		panicOnZeroMethods := len(rd.subRouters) == 0 //Panic if we have no controller methods and also no subrouters
		for method, handler := range GetControllerMethods(rd.controller, panicOnZeroMethods) {
			for _, urlPart := range rd.urlParts {
				muxRoute := router.Handle(urlPart, newNegroniMiddleware(combinedMiddleWareHandlers, handler))
				muxRoute.Methods(method)
			}
		}

		if len(rd.urlParts) == 0 {
			setupRouters(router, combinedMiddleWareHandlers, rd.subRouters)
			continue
		}

		for _, urlPart := range rd.urlParts {
			var subRouterToUse *mux.Router
			if urlPart != "" {
				subRouterToUse = router.PathPrefix(urlPart).Subrouter()
			} else {
				subRouterToUse = router
			}
			setupRouters(subRouterToUse, combinedMiddleWareHandlers, rd.subRouters)
		}
	}
}

func RegisterRouters(router *mux.Router, baseMiddleWare []http.HandlerFunc, routers []*Router) {
	//Routers
	setupRouters(router, baseMiddleWare, routers)
}
