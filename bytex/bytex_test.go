package bytex

import "testing"

func TestByteIndexFold(t *testing.T) {
	for i, tc := range []struct {
		haystack, needle string
		want             int
	}{
		{"body", "body", 0},
	} {
		got := ByteIndexFold([]byte(tc.haystack), []byte(tc.needle))
		if got != tc.want {
			t.Errorf("%d. got %d, wanted %d.", i, got, tc.want)
		}
	}
}

func TestByteHasPrefixFold(t *testing.T) {
	for i, tc := range []struct {
		haystack, needle string
		want             bool
	}{
		{"body", "body", true},
	} {
		got := ByteHasPrefixFold([]byte(tc.haystack), []byte(tc.needle))
		if got != tc.want {
			t.Errorf("%d. got %t, wanted %t.", i, got, tc.want)
		}
	}
}
