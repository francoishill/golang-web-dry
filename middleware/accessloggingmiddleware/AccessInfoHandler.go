package accessloggingmiddleware

type AccessInfoHandler interface {
	OnStart(info *StartAccessInfo)
	OnEnd(info *EndAccessInfo)
}

type simpleAccessInfoHandler struct {
	onStart func(info *StartAccessInfo)
	onEnd   func(info *EndAccessInfo)
}

func (s *simpleAccessInfoHandler) OnStart(info *StartAccessInfo) {
	s.onStart(info)
}
func (s *simpleAccessInfoHandler) OnEnd(info *EndAccessInfo) {
	s.onEnd(info)
}

func NewSimpleAccessInfoHandler(onStart func(info *StartAccessInfo), onEnd func(info *EndAccessInfo)) AccessInfoHandler {
	return &simpleAccessInfoHandler{onStart, onEnd}
}
