package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	var (
		defaultLevel = defaultAction.compressLevel

		file1 = func(s string) func(*action) bool {
			return func(a *action) bool { return a.fileIn == s }
		}
		file2 = func(s string) func(*action) bool {
			return func(a *action) bool { return a.fileOut == s }
		}
		compress = func(level int) func(*action) bool {
			return func(a *action) bool {
				return a.compress == true && a.compressLevel == level
			}
		}
		decompress = func(a *action) bool { return a.compress == false }
		hasHelp    = func(a *action) bool { return a.help != nil }
		none       = func(_ *action) bool { return false }
		and        = func(f, g func(*action) bool) func(*action) bool {
			return func(a *action) bool { return f(a) && g(a) }
		}

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
			{"", and(file1("-"), file2("-")), ""},
			{"-z", and(file1("-"), file2("-")), ""},
			{"-d", and(file1("-"), file2("-")), ""},
			{"-", and(file1("-"), file2("-")), ""},
			{"-z -", and(file1("-"), file2("-")), ""},
			{"-d -", and(file1("-"), file2("-")), ""},
			{"- x", and(file1("-"), file2("x")), ""},
			{"x -", and(file1("x"), file2("-")), ""},
			{"- -", and(file1("-"), file2("-")), ""},
			{"-z - -", and(file1("-"), file2("-")), ""},
			{"-d - -", and(file1("-"), file2("-")), ""},
			{"foo", file1("foo"), ""},
			{"-4 foo", and(compress(4), file1("foo")), ""},
			{"foo -4", and(compress(4), file1("foo")), ""},
			{"foo bar", and(file1("foo"), file2("bar")), ""},
			{"foo", and(file1("foo"), file2("foo.deflate")), ""},
			{"-z foo", and(file1("foo"), file2("foo.deflate")), ""},
			{"-d foo.deflate", and(file1("foo.deflate"), file2("foo")), ""},
			{"-d foo", none, "unable to guess 2nd file name"},
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
