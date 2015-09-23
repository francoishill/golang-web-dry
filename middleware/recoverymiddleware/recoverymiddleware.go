package recoverymiddleware

// Inspiration from:
//  - https://github.com/albrow/negroni-json-recovery/blob/master/recovery.go
//  - https://github.com/codegangsta/negroni/blob/master/recovery.go

import (
	"fmt"
	"net/http"

	"github.com/unrolled/render"

	. "github.com/francoishill/golang-web-dry/errors/stacktraces/prettystacktrace"
)

var IndentJSON = false

type RecoveredErrorDetails struct {
	OriginalError interface{}
	Error         string
	StackTrace    string
}

type RecoveryResponse struct {
	StatusCode         int
	JsonResponseObject interface{}
}

type Recovery struct {
	WithRecoveredError func(errDetails *RecoveredErrorDetails) *RecoveryResponse
}
type jsonPanicError struct {
	Code  int    `json:",omitempty"` // the http response code
	Short string `json:",omitempty"` // a short explanation of the response (usually one or two words). for internal use only
	Error string `json:",omitempty"` // any errors that may have occured with the request and should be displayed to the user
	// From  string `json:",omitempty"` // the file and line number from which the error originated
}

func NewRecovery(withRecoveredError func(errDetails *RecoveredErrorDetails) *RecoveryResponse) *Recovery {
	return &Recovery{
		WithRecoveredError: withRecoveredError,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			// convert err to a string
			var errMsg string
			if e, ok := err.(error); ok {
				errMsg = e.Error()
			} else {
				errMsg = fmt.Sprint(err)
			}
			stack := GetPrettyStackTrace()

			var statusCode int = 0
			var jsonResponseObject interface{}

			if rec.WithRecoveredError != nil {
				recoveryResponse := rec.WithRecoveredError(&RecoveredErrorDetails{err, errMsg, stack})
				if recoveryResponse != nil {
					statusCode = recoveryResponse.StatusCode
					jsonResponseObject = recoveryResponse.JsonResponseObject
				}
			}

			if statusCode == 0 {
				statusCode = http.StatusInternalServerError
			}
			if jsonResponseObject == nil {
				jsonResponseObject = &jsonPanicError{statusCode, "InternalError", errMsg}
			}

			rend := render.New(render.Options{
				IndentJSON: IndentJSON,
			})
			rend.JSON(rw, statusCode, jsonResponseObject)
		}
	}()

	next(rw, r)
}
