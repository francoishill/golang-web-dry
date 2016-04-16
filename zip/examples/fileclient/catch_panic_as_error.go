package fileclient

import (
	"errors"
	"fmt"
)

func CatchPanicAsError(errPointer *error) {
	r := recover()
	if r == nil {
		return
	}

	switch t := r.(type) {
	case error:
		*errPointer = t
	case string:
		*errPointer = errors.New(t)
	default:
		*errPointer = fmt.Errorf("%#v", t)
	}
}
