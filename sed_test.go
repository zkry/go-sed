package gosed

import "testing"

// TestInfo tests to see if the ending positions returned from Info
// are as expected.
func TestInfo(t *testing.T) {
	cases := []struct {
		program   string
		positions []int
	}{
		{"s/one/two/g", []int{1, 2, 5, 6, 9, 10, 11}},
	}

infoTests:
	for _, c := range cases {
		tokens := Info(c.program)
		for i := range tokens {
			if len(c.positions) == 0 {
				break
			}
			if tokens[i].End != c.positions[0] {
				t.Errorf("Info produced wrong positions:\n  Got: %v\n  Expected: %v\n", tokens[i].End, c.positions[i])
				continue infoTests
			}
			c.positions = c.positions[1:]
		}
	}
}
