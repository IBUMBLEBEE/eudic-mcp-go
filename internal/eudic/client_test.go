package eudic

import "testing"

func TestIDString(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{nil, "0"},
		{"", "0"},
		{"132", "132"},
		{float64(42), "42"},
		{0, "0"},
	}
	for _, c := range cases {
		if got := IDString(c.in); got != c.want {
			t.Fatalf("IDString(%v)=%q want %q", c.in, got, c.want)
		}
	}
}
