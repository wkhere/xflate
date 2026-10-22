package main

import (
	"fmt"
	"strings"

	"codeberg.org/wkhere/argp"
	pc "codeberg.org/wkhere/argp/parsecommon"
)

func parseArgs(args []string) (a action, _ error) {
	a = defaultAction

	p := &argp.PState{Args: args}

	var flag struct {
		z, d, t    bool
		zset, dset bool
	}
	var rest []string

	for ; len(p.Args) > 0 && p.Err == nil; p.Next() {
		switch arg := p.Arg(); {

		case p.MatchBoolFlag("-z", "--compress"):
			flag.zset = pc.ParseBool(&flag.z, p)

		case p.MatchBoolFlag("-d", "--decompress"):
			flag.dset = pc.ParseBool(&flag.d, p)

		case p.MatchDashIntFlag() || p.MatchLongFlagExpr("--level"):
			pc.ParseInt(&a.compressLevel, p)

		case p.MatchBoolFlag("-f", "--force"):
			pc.ParseBool(&a.force, p)

		case p.MatchBoolFlag("-t", "--test"):
			pc.ParseBool(&flag.t, p)

		case p.MatchJustFlag("-h", "--help"):
			a.help = toplevelHelp
			return a, nil

		case p.MatchAnyFlag():
			p.Errorf("unknown flag: %s", p.Flag)

		default:
			rest = append(rest, arg)
		}
	}

	if p.Err != nil {
		return a, p.Err
	}

	switch {
	case flag.t && (flag.zset || flag.dset):
		return a, fmt.Errorf("conflicting flags: use -t without -z or -d")
	case flag.zset && flag.dset && flag.z == flag.d:
		return a, fmt.Errorf("conflicting flags -z=%v and -d=%v", flag.z, flag.d)
	case flag.zset:
		a.compress = flag.z
	case flag.dset:
		a.compress = !flag.d
	}

	if flag.t {
		switch len(rest) {
		case 0:
			a.fileIn = "-"
		case 1:
			a.fileIn = rest[0]
		default:
			return a, fmt.Errorf("-t accepts at most one file name")
		}
		a.fileOut = discard
		a.compress = false
		return a, nil
	}

	switch len(rest) {
	case 0:
		a.fileIn = "-"
		a.fileOut = "-"
	case 1:
		switch f1 := rest[0]; {
		case f1 == "-":
			a.fileIn = "-"
			a.fileOut = "-"
		case a.compress:
			a.fileIn = f1
			a.fileOut = f1 + fileExt
		case !a.compress && strings.HasSuffix(f1, fileExt):
			a.fileIn = f1
			a.fileOut = f1[:len(f1)-len(fileExt)]
		default:
			return a, fmt.Errorf("unable to guess 2nd file name")
		}
	case 2:
		if rest[0] == rest[1] && rest[0] != "-" {
			return a, fmt.Errorf("files must be different")
		}
		a.fileIn, a.fileOut = rest[0], rest[1]
	default:
		return a, fmt.Errorf("too many file args")
	}

	return a, nil
}

// help

func toplevelHelp() {
	fmt.Printf(usage, defaultAction.compressLevel)
}

const usage = `Compress/decompress files according to the deflate format.

NOTE: The raw deflate format differs from both the zlib and gzip formats.

Usage: xflate [FLAGS] [FILE1] [FILE2]
    -z, --compress     compress (default true)
    -d, --decompress   decompress
    -N, --level=N      compress level, -2..9 (default %d)
    -f, --force        force overwriting FILE2
    -t, --test         test compressed FILE1
    -h, --help         show this help and exit
When compressing and only FILE1 is given, FILE2 is FILE1.deflate .
When decompressing and only FILE1 is given, FILE2 tries to strip .deflate from it.
FILE1 or FILE2 (or both) can be "-", meaning stdin and stdout.
NOTE: FILE1 is not deleted afterwards.
`
