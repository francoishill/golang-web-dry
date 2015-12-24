package negroni_mux_routers

import (
	"net/http"
)

type ControllerMethod func(w http.ResponseWriter, r *http.Request)
