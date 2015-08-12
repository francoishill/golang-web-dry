package recoverymiddleware

// Inspiration from:
//  - https://github.com/albrow/negroni-json-recovery/blob/master/recovery.go
//  - https://github.com/codegangsta/negroni/blob/master/recovery.go

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/unrolled/render"

	. "github.com/francoishill/golang-web-dry/errors/stacktraces/prettystacktrace"
)

var IndentJSON = false

type Recovery struct {
	Logger *log.Logger
}
type jsonPanicError struct {
	Code  int    `json:",omitempty"` // the http response code
	Short string `json:",omitempty"` // a short explanation of the response (usually one or two words). for internal use only
	Error string `json:",omitempty"` // any errors that may have occured with the request and should be displayed to the user
	// From  string `json:",omitempty"` // the file and line number from which the error originated
}

func NewRecovery() *Recovery {
	return &Recovery{
		Logger: log.New(os.Stdout, "[negroni-mod] ", 0),
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			stack := GetPrettyStackTrace()

			rec.Logger.Printf("PANIC: %s\n%s", err, stack)

			r := render.New(render.Options{
				IndentJSON: IndentJSON,
			})

			// convert err to a string
			var errMsg string
			if e, ok := err.(error); ok {
				errMsg = e.Error()
			} else {
				errMsg = fmt.Sprint(err)
			}

			r.JSON(rw, 500, &jsonPanicError{500, "InternalError", errMsg})
		}
	}()

	next(rw, r)
}
