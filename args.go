package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseArgs(args []string) (a action, _ error) {
	a = defaultAction

	errs := runArgp(&a, args)
	if len(errs) > 0 {
		return a, errs[0]
	}
	return a, nil
}

// argp struct & api entry

type argp struct {
	input  []string
	pos    int
	lastw  int
	hasArg bool

	a    *action
	errs []error
}

func runArgp(a *action, input []string) []error {
	p := argp{
		input: input,
		a:     a,
	}
	p.run()
	return p.errs
}

// basic primitives

type argpStateFn func(*argp) argpStateFn

func (p *argp) run() {
	for st := argpStart; st != nil; {
		st = st(p)
	}
}

func (p *argp) emit(f func(*action)) {
	f(p.a)
}

func (p *argp) emitError(err error) {
	p.errs = append(p.errs, err)
}

func (p *argp) emitErrorf(format string, a ...any) {
	p.emitError(fmt.Errorf(format, a...))
}

func (p *argp) read() (s string) {
	if len(p.input[p.pos:]) == 0 {
		p.lastw = 0
		return
	}
	s = p.input[p.pos]
	p.lastw = 1
	p.pos++
	return s
}

func (p *argp) backup() {
	p.pos -= p.lastw
}

// finalizing helpers

func (p *argp) final(f func(*action)) argpStateFn {
	p.emit(f)
	return nil
}

func (p *argp) error(err error) argpStateFn {
	p.emitError(err)
	return nil
}

func (p *argp) errorf(format string, a ...any) argpStateFn {
	p.emitErrorf(format, a...)
	return nil
}

// state functions

func argpStart(p *argp) argpStateFn {
	var flag struct {
		z, d       bool
		zset, dset bool
		lvl        int
	}
loop:
	for {
		switch s := p.read(); {
		case s == "":
			break loop

		case p.parseBoolFlag(&flag.z, s, "-z", "--compress"):
			flag.zset = true

		case p.parseBoolFlag(&flag.d, s, "-d", "--decompress"):
			flag.dset = true

		case p.parseIntFlag(&flag.lvl, s, "--level"):
			p.a.compressLevel = flag.lvl

		case s == "-h" || s == "--help":
			return p.final(toplevelHelp)

		case len(s) > 1 && s[0] == '-':
			return p.errorf("unknown flag: %s", s)

		default:
			return p.errorf("unexpected args, use stdin/stdout")
		}
	}

	switch {
	case flag.zset && flag.dset && flag.z == flag.d:
		return p.errorf("conflicting flags -z=%v and -d=%v", flag.z, flag.d)
	case flag.zset:
		p.a.compress = flag.z
	case flag.dset:
		p.a.compress = !flag.d
	}
	return nil
}

// parsing bits, for case expressions

func (p *argp) parseBoolFlag(dest *bool, s, short, long string) bool {
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

func (p *argp) parseIntFlag(dest *int, s, long string) bool {
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

// emited bits

func toplevelHelp(a *action) {
	a.help = func() {
		fmt.Printf(usage, defaultAction.compressLevel)
	}
}

const usage = `Usage: xflate [flags]  (reads stdin, outputs to stdout)
  -z, --compress     compress (default true)
  -d, --decompress   decompress
  -N, --level=N      compress level, -2..9 (default %d)
  -h, --help         show this help and exit
`
