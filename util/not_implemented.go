package util

import "fmt"

type NotImplementedError struct {
	Feature string
}

func (e NotImplementedError) Error() string {
	if e.Feature != "" {
		return fmt.Sprintf("not implemented: %s", e.Feature)
	} else {
		return "not implemented"
	}
}
