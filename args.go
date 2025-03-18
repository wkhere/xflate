package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseArgs(args []string) (a action, _ error) {
	a = defaultAction

	var p pstate

	var flag struct {
		z, d, t, f bool
		zset, dset bool
		lvl        int
	}
	var rest []string

loop:
	for ; len(args) > 0 && p.err == nil; args = args[1:] {
		switch s := args[0]; {
		case s == "":
			break loop

		case p.parseBoolFlag(&flag.z, s, "-z", "--compress"):
			flag.zset = true

		case p.parseBoolFlag(&flag.d, s, "-d", "--decompress"):
			flag.dset = true

		case p.parseIntFlag(&flag.lvl, s, "--level"):
			a.compressLevel = flag.lvl

		case p.parseBoolFlag(&flag.f, s, "-f", "--force"):
			a.force = flag.f

		case p.parseBoolFlag(&flag.t, s, "-t", "--test"): // ok

		case s == "-h" || s == "--help":
			a.help = toplevelHelp
			return a, nil

		case len(s) > 1 && s[0] == '-':
			p.emitErrorf("unknown flag: %s", s)

		default:
			rest = append(rest, s)
		}
	}

	if p.err != nil {
		return a, p.err
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

type pstate struct {
	err error
}

func (p *pstate) emitErrorf(format string, a ...any) {
	// saving only the first error:
	if p.err == nil {
		p.err = fmt.Errorf(format, a...)
	}
}

// parsing bits, for case expressions

func (p *pstate) parseBoolFlag(dest *bool, s, short, long string) bool {
	var n int
	if strings.HasPrefix(s, long) {
		n = len(long)
	} else if strings.HasPrefix(s, short) {
		n = len(short)
	} else {
		return false
	}
	flag, s := s[:n], s[n:]
	if s == "" {
		*dest = true
		return true
	}
	if s[0] != '=' {
		return false
	}
	v, err := strconv.ParseBool(s[1:])
	if err != nil {
		p.emitErrorf("flag %s: %w", flag, err)
		return false
	}
	*dest = v
	return true
}

func (p *pstate) parseIntFlag(dest *int, s, long string) bool {
	var n int
	var dashN bool
	if strings.HasPrefix(s, long) {
		n = len(long)
	} else if s[0] == '-' {
		n = 1
		dashN = true
	} else {
		return false
	}
	flag, s := s[:n], s[n:]

	if dashN {
		if s == "" {
			return false
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		*dest = v
		return true
	}

	if s == "" {
		p.emitErrorf("flag %s: expected value", flag)
		return false
	}
	if s[0] != '=' {
		return false
	}
	v, err := strconv.Atoi(s[1:])
	if err != nil {
		p.emitErrorf("flag %s: %w", flag, err)
		return false
	}
	*dest = v
	return true
}

// help

func toplevelHelp() {
	fmt.Printf(usage, defaultAction.compressLevel)
}

const usage = `Usage: xflate [FLAGS] [FILE1] [FILE2]
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
