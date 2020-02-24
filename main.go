// Xflate decompresses or compresses deflate stream.
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
	prog = os.Args[0]
	usageHead = fmt.Sprintf(
		"Usage: %s\t(reads stdin, outputs to stdout)", prog)
}

type config struct {
	compress      bool
	compressLevel int
}

func parseFlags(args []string) config {
	var conf config
	var help bool

	flag := pflag.NewFlagSet("flags", pflag.ContinueOnError)
	flag.SortFlags = false

	flag.BoolVarP(&conf.compress, "compress", "z", false,
		"compress (default false -- means decompress)")
	// TODO: complement with --decompress -d flag
	flag.IntVarP(&conf.compressLevel, "level", "n", 6,
		"compress level, -2..9")
	flag.BoolVarP(&help, "help", "h", false,
		"show this help and exit")
	flag.Usage = func() { usage(os.Stderr, flag) }

	err := flag.Parse(args)
	if err != nil {
		usage(os.Stderr, flag)
		os.Exit(2)
	}
	if help {
		usage(os.Stdout, flag)
		os.Exit(0)
	}

	return conf
}

func usage(w io.Writer, f *pflag.FlagSet) {
	fmt.Fprintln(w, usageHead)
	f.SetOutput(w)
	f.PrintDefaults()
}

func main() {
	conf := parseFlags(os.Args[1:])

	switch {
	case conf.compress:
		w, err := flate.NewWriter(os.Stdout, conf.compressLevel)
		if err != nil {
			die(fmt.Errorf("failed creating compress writer: %v", err))
		}
		defer w.Close()
		_, err = io.Copy(w, os.Stdin)
		if err != nil {
			die(fmt.Errorf("compress: %v", err))
		}

	default:
		r := flate.NewReader(os.Stdin)
		defer r.Close()
		_, err := io.Copy(os.Stdout, r)
		if err != nil {
			die(fmt.Errorf("decompress: %v", err))
		}
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, prog+":", err)
	os.Exit(1)
}
