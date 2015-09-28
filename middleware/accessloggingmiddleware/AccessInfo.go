package accessloggingmiddleware

import (
	"time"
)

type AccessInfo struct {
	HttpMethod string
	RemoteAddr string
	RemoteIP   string
	RequestURI string
	Proxies    []string
	UserAgent  string
}

type StartAccessInfo struct {
	*AccessInfo
	StartTime time.Time
}

type EndAccessInfo struct {
	*AccessInfo
	GotPanic   bool
	Status     int
	StatusText string
	EndTime    time.Time
	Duration   time.Duration
}
