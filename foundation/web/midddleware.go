package web

type MidFunc func(handler HandlerFunc) HandlerFunc

func wrapMiddleware(mw []MidFunc, handler HandlerFunc) HandlerFunc {
	// loog backwards through the middleware invoking each one. Replace
	// the handler with the new wrapped handler. Looping backwards ensures
	// that the first middleware of the dlice is the first to be executed by requests.

	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := mw[i]
		if mwFunc != nil {
			handler = mwFunc(handler)
		}
	}

	return handler
}
