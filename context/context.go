package context

// Context represents information that will be passed to each
// page when rendered. E.g. it allows you to render the title
// in a <title> tag inside of an ace template.
type Context map[string]interface{}

// context is a private copy of the global context. It should
// only be accessed or modified by its methods.
var context = Context{}

// GetContext returns the context
func GetContext() Context {
	return context
}

// Add adds val to the context identified by key
func Add(key string, val interface{}) {
	context[key] = val
}

// Copy context returns a copy of the context. Modifying
// it will not change the original context. This function
// can be used to create a per-page or per-post context.
func CopyContext() Context {
	contextCopy := Context{}
	for k, v := range context {
		contextCopy[k] = v
	}
	return contextCopy
}
