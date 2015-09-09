package ziputils

type SimpleLogger interface {
	Debug(msg string, msgArgs ...interface{})
}
