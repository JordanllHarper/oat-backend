package main

import "strings"

func sliceUnorderedRemove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func trimPreserveNewline(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return r == ' '
	})
}
