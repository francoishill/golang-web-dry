package prettystacktrace

import (
	"bytes"
	"fmt"
	"runtime"
)

func GetPrettyStackTrace() string {
	var buf bytes.Buffer
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buf.WriteString(fmt.Sprintln(fmt.Sprintf("%s:%d", file, line)))
	}
	return buf.String()
}
