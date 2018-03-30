package utilsx

import (
	"fmt"
	"testing"
)

func TestWordCount(t *testing.T) {
	case1Want := make(map[string]int, 4)
	case1Want["Australia"] = 2
	case1Want["Canada"] = 2
	case1Want["Germany"] = 1
	case1Want["Japan"] = 1

	cases := []struct {
		in   string
		want map[string]int
	}{
		{"Australia Canada Germany Australia Japan Canada", case1Want},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := WordCount(tc.in)
			for i, ii := range out {
				if case1Want[i] != ii {
					t.Errorf("\nout:  %#v\nwant: %#v\n", ii, case1Want[i])
				}
			}
		})
	}
}
