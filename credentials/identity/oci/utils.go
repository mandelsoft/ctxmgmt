package oci

import (
	"strings"
)

func SplitLocator(locator string) (string, string, string) {
	path := ""
	h := ""
	i := strings.Index(locator, "/")
	if i < 0 {
		h = locator
	} else {
		h = locator[:i]
		path = locator[i+1:]
	}
	i = strings.Index(h, ":")

	if i < 0 {
		return h, "", path
	}
	return h[:i], h[i+1:], path
}
