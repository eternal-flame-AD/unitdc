//go:build js && wasm

package localize

import (
	"syscall/js"
)

func getLang() (locale string) {
	lang := js.Global().Get("navigator").Get("language")
	if !lang.Truthy() {
		return ""
	}
	return lang.String()
}
