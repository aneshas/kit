package http

// Adapter represents HandlerFunc middleware adapter
type Adapter func(HandlerFunc) HandlerFunc

// AdaptHandlerFunc decorates given HandlerFunc with provided adapters
func AdaptHandlerFunc(hf HandlerFunc, adapters ...Adapter) HandlerFunc {
	for _, adapter := range adapters {
		hf = adapter(hf)
	}
	return hf
}
