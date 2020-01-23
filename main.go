package main

import (
	"compress/flate"
	"flag"
	"fmt"
	"io"
	"os"
)

func parseFlags() {
	var help bool
	flag.BoolVar(&help, "h", false, "show this help and exit")
	flag.Usage = usage
	flag.Parse()

	if help {
		flag.CommandLine.SetOutput(os.Stdout)
		usage()
		os.Exit(0)
	}
}

func usage() {
	o := flag.CommandLine.Output()
	prog := os.Args[0]
	fmt.Fprintf(o, "Usage: %s\t\t(reads stdin, outputs to stdout)\n", prog)
	flag.PrintDefaults()
}

func main() {
	parseFlags()

	_, err := io.Copy(os.Stdout, flate.NewReader(os.Stdin))
	if err != nil {
		die(fmt.Errorf("deflate: %v", err))
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
