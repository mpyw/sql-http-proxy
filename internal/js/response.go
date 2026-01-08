package js

import (
	"net/http"

	"github.com/dop251/goja"
	"github.com/samber/lo"
)

// Response implements a simplified Web API Response interface for use in JavaScript.
// It provides writable access to response properties for post-transform.
type Response struct {
	status     int
	statusText string
	headers    *Headers
}

// NewResponse creates a new Response object with default values.
func NewResponse() *Response {
	return &Response{
		status:     http.StatusOK,
		statusText: "OK",
		headers:    NewWritableHeaders(nil),
	}
}

// Status returns the HTTP status code.
func (r *Response) Status() int {
	return r.status
}

// SetStatus sets the HTTP status code.
// Status must be in the valid HTTP range (100-599).
// Invalid values are ignored.
func (r *Response) SetStatus(status int) {
	if status < 100 || status > 599 {
		return // Ignore invalid status codes
	}
	r.status = status
	// Update statusText to match common status codes
	r.statusText = http.StatusText(status)
	if r.statusText == "" {
		r.statusText = "Unknown"
	}
}

// StatusText returns the HTTP status text.
func (r *Response) StatusText() string {
	return r.statusText
}

// SetStatusText sets the HTTP status text.
func (r *Response) SetStatusText(text string) {
	r.statusText = text
}

// Headers returns the response headers.
func (r *Response) Headers() *Headers {
	return r.headers
}

// Ok returns true if status is in the 200-299 range.
func (r *Response) Ok() bool {
	return r.status >= 200 && r.status < 300
}

// ToHTTPHeader returns the underlying http.Header for writing to response.
func (r *Response) ToHTTPHeader() http.Header {
	return r.headers.ToHTTPHeader()
}

// ToJSObject creates a goja object with Response properties and methods.
// The object is sealed to prevent modification of its structure.
func (r *Response) ToJSObject(vm *goja.Runtime) *goja.Object {
	obj := vm.NewObject()

	// status: get/set with type validation
	lo.Must0(obj.DefineAccessorProperty("status",
		vm.ToValue(func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(r.status)
		}),
		vm.ToValue(func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				status := call.Argument(0).ToInteger()
				r.SetStatus(int(status))
			}
			return goja.Undefined()
		}),
		goja.FLAG_FALSE, goja.FLAG_TRUE))

	// statusText: get/set with type validation
	lo.Must0(obj.DefineAccessorProperty("statusText",
		vm.ToValue(func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(r.statusText)
		}),
		vm.ToValue(func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				r.SetStatusText(call.Argument(0).String())
			}
			return goja.Undefined()
		}),
		goja.FLAG_FALSE, goja.FLAG_TRUE))

	// headers: read-only getter (but the Headers object itself is writable)
	headersObj := r.headers.ToJSObject(vm)
	lo.Must0(obj.DefineAccessorProperty("headers",
		vm.ToValue(func(call goja.FunctionCall) goja.Value {
			return headersObj
		}),
		nil, goja.FLAG_FALSE, goja.FLAG_TRUE))

	// ok: computed read-only getter
	lo.Must0(obj.DefineAccessorProperty("ok",
		vm.ToValue(func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(r.Ok())
		}),
		nil, goja.FLAG_FALSE, goja.FLAG_TRUE))

	// Seal the object to prevent adding/removing properties
	seal := lo.Must(vm.RunString("Object.seal"))
	sealFn := lo.Must(goja.AssertFunction(seal))
	lo.Must(sealFn(goja.Undefined(), obj))

	return obj
}
