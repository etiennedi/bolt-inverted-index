package main

import (
	"strings"
	"unicode"
)

type analyzer struct{}

func (a *analyzer) splitAndLowercase(in string) []string {
	parts := strings.FieldsFunc(in, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	lowercased := make([]string, len(parts))
	for i, word := range parts {
		lowercased[i] = strings.ToLower(word)
	}

	return lowercased
}
