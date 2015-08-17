package gorillamuxurlparams

import (
	"net/http"
	"strconv"

	. "github.com/francoishill/golang-web-dry/errors/checkerror"

	"github.com/gorilla/mux"
)

func MustGetUrlParamValue_String(r *http.Request, paramName string) string {
	vars := mux.Vars(r)
	paramValue, varFound := vars[paramName]
	if !varFound {
		panic(paramName + " cannot be found from URL")
	}
	return paramValue
}

func MustGetUrlParamValue_Int64(r *http.Request, paramName string) int64 {
	paramValue := MustGetUrlParamValue_String(r, paramName)
	intVal, err := strconv.ParseInt(paramValue, 10, 64)
	CheckError(err)
	return intVal
}
