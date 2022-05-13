// xflate compresses or decompresses deflate stream.
package main

import (
	"compress/flate"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

var prog, usageHead string

func init() {
	prog = "xflate"
	usageHead = fmt.Sprintf(
		"Usage: %s\t(reads stdin, outputs to stdout)", prog)
}

type config struct {
	compress      bool
	compressLevel int
}

func parseFlags(args []string) config {
	var (
		conf       config
		decompress bool
		help       bool
	)

	flag := pflag.NewFlagSet("flags", pflag.ContinueOnError)
	flag.SortFlags = false

	flag.BoolVarP(&conf.compress, "compress", "z", true, "compress")
	flag.BoolVarP(&decompress, "decompress", "d", false, "decompress")
	flag.IntVarP(&conf.compressLevel, "level", "n", 6,
		"compress level, -2..9")
	flag.BoolVarP(&help, "help", "h", false,
		"show this help and exit")
	flag.Usage = func() {
		fmt.Fprintln(flag.Output(), usageHead)
		flag.PrintDefaults()
	}

	err := flag.Parse(args)
	if err != nil {
		die(2, err)
	}
	if help {
		flag.SetOutput(os.Stdout)
		flag.Usage()
		die(0)
	}

	if flag.Changed("compress") && flag.Changed("decompress") &&
		conf.compress == decompress {
		die(2, "conflicting flags -z and -d")
	}
	if flag.Changed("decompress") {
		conf.compress = !decompress
	}

	return conf
}

func main() {
	conf := parseFlags(os.Args[1:])

	switch {
	case conf.compress:
		w, err := flate.NewWriter(os.Stdout, conf.compressLevel)
		if err != nil {
			die(1, "failed creating compress writer:", err)
		}
		_, err = io.Copy(w, os.Stdin)
		if err != nil {
			die(1, "compress:", err)
		}
		err = w.Close()
		if err != nil {
			die(1, "compress closing", err)
		}

	case !conf.compress:
		r := flate.NewReader(os.Stdin)
		defer r.Close()
		_, err := io.Copy(os.Stdout, r)
		if err != nil {
			die(1, "decompress:", err)
		}
	}
}

func die(exitcode int, msgs ...interface{}) {
	if len(msgs) > 0 {
		fmt.Fprint(os.Stderr, prog, ": ")
		fmt.Fprintln(os.Stderr, msgs...)
	}
	os.Exit(exitcode)
}
