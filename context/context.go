package context

type Context map[string]interface{}

var context = Context{}

func GetContext() Context {
	return context
}

func Add(key string, val interface{}) {
	context[key] = val
}