package sanitize

import "strings"

func Path(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
