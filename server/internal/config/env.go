package config

import "strings"

func NewReplacer() *strings.Replacer {
	replacer := strings.NewReplacer(".", "_")
	return replacer
}
