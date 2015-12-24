package negroni_mux_routers

import (
	"fmt"
	"net/http"
)

type Controller interface {
	RelativeURLPatterns() []string
}

type optionsHandler interface {
	Options(w http.ResponseWriter, r *http.Request)
}
type getHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}
type headHandler interface {
	Head(w http.ResponseWriter, r *http.Request)
}
type postHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
}
type putHandler interface {
	Put(w http.ResponseWriter, r *http.Request)
}
type deleteHandler interface {
	Delete(w http.ResponseWriter, r *http.Request)
}

func GetControllerMethods(controller Controller, panicOnZeroMethods bool) map[string]ControllerMethod {
	m := make(map[string]ControllerMethod)

	cnt := 0

	if o, ok := controller.(optionsHandler); ok {
		cnt++
		m["OPTIONS"] = o.Options
	}

	if g, ok := controller.(getHandler); ok {
		cnt++
		m["GET"] = g.Get
	}

	if h, ok := controller.(headHandler); ok {
		cnt++
		m["HEAD"] = h.Head
	}

	if h, ok := controller.(postHandler); ok {
		cnt++
		m["POST"] = h.Post
	}

	if h, ok := controller.(putHandler); ok {
		cnt++
		m["PUT"] = h.Put
	}

	if h, ok := controller.(deleteHandler); ok {
		cnt++
		m["DELETE"] = h.Delete
	}

	if panicOnZeroMethods && cnt == 0 {
		panic(fmt.Sprintf("Controller '%#T' must have at least one exposed method 'Get', 'Put', etc.", controller))
	}

	return m
}
