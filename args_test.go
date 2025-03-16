package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	var (
		defaultLevel = defaultAction.compressLevel

		compress = func(level int) func(*action) bool {
			return func(a *action) bool {
				return a.compress == true && a.compressLevel == level
			}
		}
		decompress = func(a *action) bool { return a.compress == false }
		hasHelp    = func(a *action) bool { return a.help != nil }
		none       = func(_ *action) bool { return false }

		tab = []struct {
			input string
			want  func(*action) bool
			errs  string
		}{
			{"", compress(defaultLevel), ""},
			{"-h", hasHelp, ""},
			{"-q", none, "unknown flag"},
			{"-z -d", none, "conflicting flags"},
			{"-z=1 -d=1", none, "conflicting flags"},
			{"-z=0 -d=0", none, "conflicting flags"},
			{"-z=true -d=true", none, "conflicting flags"},
			{"-z=false -d=false", none, "conflicting flags"},
			{"-z", compress(defaultLevel), ""},
			{"-d", decompress, ""},
			{"-1", compress(1), ""},
			{"-2", compress(2), ""},
			{"-9", compress(9), ""},
			{"foo", none, "unexpected args"},
			{"foo bar", none, "unexpected args"},
		}
	)

	for i, tc := range tab {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a, err := parseArgs(strings.Fields(tc.input))

			switch {
			case err != nil && tc.errs == "":
				t.Errorf("unexpected error: %s", err)

			case err != nil && tc.errs != "":
				if !strings.Contains(err.Error(), tc.errs) {
					t.Errorf("expected error with %q; having: %q", tc.errs, err)
				}

			case err == nil && tc.errs != "":
				t.Errorf("expected error with %q", tc.errs)

			default:
				if !(tc.want(&a)) {
					t.Errorf("predicate failed on %+v", a)
				}
			}
		})
	}
}
