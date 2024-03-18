package sanitize

import "testing"

func TestPath(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		res := Path("")

		if res != "" {
			t.Fatalf("Expected \"\", got %s", res)
		}
	})

	t.Run("valid string", func(t *testing.T) {
		res := Path("this/is\\some_arbitrary\\path")

		expected := "this/is/some_arbitrary/path"
		if res != expected {
			t.Fatalf("Expected %s, got %s", expected, res)
		}
	})
}
