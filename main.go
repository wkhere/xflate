package main

import (
	"compress/flate"
	"flag"
	"fmt"
	"io"
	"os"
)

var prog string

func init() {
	prog = os.Args[0]
}

type config struct {
	compress bool
}

func parseFlags() config {
	var conf config
	var help bool
	flag.BoolVar(&conf.compress, "z", false, "compress (default: false)")
	flag.BoolVar(&help, "h", false, "show this help and exit")
	flag.Usage = usage
	flag.Parse()

	if help {
		flag.CommandLine.SetOutput(os.Stdout)
		usage()
		os.Exit(0)
	}

	return conf
}

func usage() {
	o := flag.CommandLine.Output()
	fmt.Fprintf(o, "Usage: %s\t\t(reads stdin, outputs to stdout)\n", prog)
	flag.PrintDefaults()
}

func main() {
	conf := parseFlags()

	switch {
	case conf.compress:
		w, err := flate.NewWriter(os.Stdout, 6)
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
