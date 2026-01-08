package js

import (
	"net/http"
	"sort"
	"strings"

	"github.com/dop251/goja"
	"github.com/samber/lo"
)

// Headers implements the Web API Headers interface for use in JavaScript.
// It provides case-insensitive header access with get/set/has/delete methods.
type Headers struct {
	headers  http.Header
	readonly bool
}

// NewHeaders creates a new Headers object from http.Header.
func NewHeaders(h http.Header, readonly bool) *Headers {
	if h == nil {
		h = make(http.Header)
	}
	return &Headers{headers: h, readonly: readonly}
}

// NewReadonlyHeaders creates a read-only Headers object.
func NewReadonlyHeaders(h http.Header) *Headers {
	return NewHeaders(h, true)
}

// NewWritableHeaders creates a writable Headers object.
func NewWritableHeaders(h http.Header) *Headers {
	return NewHeaders(h, false)
}

// Get returns the value for a header name, or null if not found.
// Case-insensitive lookup.
func (h *Headers) Get(name string) any {
	value := h.headers.Get(name)
	if value == "" {
		return nil
	}
	return value
}

// Has returns true if the header exists.
// Case-insensitive lookup.
func (h *Headers) Has(name string) bool {
	return h.headers.Get(name) != ""
}

// Set sets a header value, replacing any existing value.
// Only works on writable Headers.
func (h *Headers) Set(name, value string) {
	if h.readonly {
		return
	}
	h.headers.Set(name, value)
}

// Append adds a value to a header, allowing multiple values.
// Only works on writable Headers.
func (h *Headers) Append(name, value string) {
	if h.readonly {
		return
	}
	h.headers.Add(name, value)
}

// Delete removes a header.
// Only works on writable Headers.
func (h *Headers) Delete(name string) {
	if h.readonly {
		return
	}
	h.headers.Del(name)
}

// Entries returns all header entries as [name, value] pairs.
func (h *Headers) Entries() [][2]string {
	var entries [][2]string
	for name, values := range h.headers {
		for _, value := range values {
			entries = append(entries, [2]string{strings.ToLower(name), value})
		}
	}
	// Sort for consistent ordering
	sort.Slice(entries, func(i, j int) bool {
		return entries[i][0] < entries[j][0]
	})
	return entries
}

// Keys returns all header names.
func (h *Headers) Keys() []string {
	keys := lo.Uniq(lo.Map(lo.Keys(h.headers), func(s string, _ int) string {
		return strings.ToLower(s)
	}))
	sort.Strings(keys)
	return keys
}

// Values returns all header values.
func (h *Headers) Values() []string {
	return lo.Flatten(lo.Values(h.headers))
}

// ForEach calls a function for each header entry.
func (h *Headers) ForEach(vm *goja.Runtime, callback goja.Callable) error {
	entries := h.Entries()
	for _, entry := range entries {
		_, err := callback(goja.Undefined(),
			vm.ToValue(entry[1]), // value
			vm.ToValue(entry[0]), // key
			vm.ToValue(h),        // this
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToHTTPHeader returns the underlying http.Header.
func (h *Headers) ToHTTPHeader() http.Header {
	return h.headers
}

// ToJSObject creates a goja object with Headers methods.
func (h *Headers) ToJSObject(vm *goja.Runtime) *goja.Object {
	obj := vm.NewObject()

	// Standard methods
	lo.Must0(obj.Set("get", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Null()
		}
		name := call.Argument(0).String()
		result := h.Get(name)
		if result == nil {
			return goja.Null()
		}
		return vm.ToValue(result)
	}))

	lo.Must0(obj.Set("has", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(false)
		}
		name := call.Argument(0).String()
		return vm.ToValue(h.Has(name))
	}))

	lo.Must0(obj.Set("set", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return goja.Undefined()
		}
		name := call.Argument(0).String()
		value := call.Argument(1).String()
		h.Set(name, value)
		return goja.Undefined()
	}))

	lo.Must0(obj.Set("append", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return goja.Undefined()
		}
		name := call.Argument(0).String()
		value := call.Argument(1).String()
		h.Append(name, value)
		return goja.Undefined()
	}))

	lo.Must0(obj.Set("delete", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		name := call.Argument(0).String()
		h.Delete(name)
		return goja.Undefined()
	}))

	lo.Must0(obj.Set("keys", func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(h.Keys())
	}))

	lo.Must0(obj.Set("values", func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(h.Values())
	}))

	lo.Must0(obj.Set("entries", func(call goja.FunctionCall) goja.Value {
		entries := h.Entries()
		result := make([]any, len(entries))
		for i, entry := range entries {
			result[i] = []string{entry[0], entry[1]}
		}
		return vm.ToValue(result)
	}))

	lo.Must0(obj.Set("forEach", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		callback, ok := goja.AssertFunction(call.Argument(0))
		if !ok {
			return goja.Undefined()
		}
		lo.Must0(h.ForEach(vm, callback))
		return goja.Undefined()
	}))

	return obj
}
