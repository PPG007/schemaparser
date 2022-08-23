package main

import (
	"fmt"
	"strings"
)

func UpperFirst(s string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(s[0:1]), s[1:])
}

func ToSnakeCase(s string) string {
	result := ""
	var lastRune rune
	for i, v := range s {
		if i == 0 {
			lastRune = v
			result = string(v)
			continue
		}
		if IsUpper(v) && !IsUpper(lastRune) {
			result = fmt.Sprintf("%s_%s", result, string(v+32))
		} else {
			result = fmt.Sprintf("%s%s", result, string(v))
		}
		lastRune = v
	}
	return result
}

func IsUpper(s rune) bool {
	if s >= 65 && s <= 90 {
		return true
	}
	return false
}
