package helpers

import "strings"

func TrimStringSpaces(items []string) []string {
	r := make([]string, 0, len(items))
	for _, i := range items {
		r = append(r, strings.TrimSpace(i))
	}
	return r
}
