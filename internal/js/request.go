package js

import (
	"net/http"
	"net/url"

	"github.com/dop251/goja"
	"github.com/samber/lo"
)

// Request implements a simplified Web API Request interface for use in JavaScript.
// It provides read-only access to request properties.
type Request struct {
	method  string
	url     *url.URL
	headers *Headers
}

// NewRequest creates a new Request object from an HTTP request.
// The URL is cloned so the value JS observes is a snapshot, independent of any
// later rewriting of r.URL by middleware or the router.
func NewRequest(r *http.Request) *Request {
	return &Request{
		method:  r.Method,
		url:     r.URL.Clone(),
		headers: NewReadonlyHeaders(r.Header),
	}
}

// Method returns the HTTP method.
func (r *Request) Method() string {
	return r.method
}

// URL returns the request URL as a string.
func (r *Request) URL() string {
	return r.url.String()
}

// Headers returns the request headers.
func (r *Request) Headers() *Headers {
	return r.headers
}

// ToJSObject creates a goja object with Request properties and methods.
// The object is sealed to prevent modification of its structure.
func (r *Request) ToJSObject(vm *goja.Runtime) *goja.Object {
	obj := vm.NewObject()

	// Read-only properties via getters
	lo.Must0(obj.DefineAccessorProperty("method", vm.ToValue(func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(r.method)
	}), nil, goja.FLAG_FALSE, goja.FLAG_TRUE))

	lo.Must0(obj.DefineAccessorProperty("url", vm.ToValue(func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(r.URL())
	}), nil, goja.FLAG_FALSE, goja.FLAG_TRUE))

	lo.Must0(obj.DefineAccessorProperty("headers", vm.ToValue(func(call goja.FunctionCall) goja.Value {
		return r.headers.ToJSObject(vm)
	}), nil, goja.FLAG_FALSE, goja.FLAG_TRUE))

	// Seal the object to prevent adding/removing properties
	seal := lo.Must(vm.RunString("Object.seal"))
	sealFn := lo.Must(goja.AssertFunction(seal))
	lo.Must(sealFn(goja.Undefined(), obj))

	return obj
}
