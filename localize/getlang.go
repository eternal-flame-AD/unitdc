//go:build !(js && wasm) && !windows

package localize

import (
	"os"
	"strings"
)

func getLang() (locale string) {
	locale = os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LANG")
	}

	if idx := strings.IndexRune(locale, '.'); idx > 0 {
		locale = locale[:idx]
	}

	return
}
