package main

import (
	"compress/flate"
	"flag"
	"io"
	"log"
	"os"
)

func init() { log.SetFlags(0) }

func parseFlags() {
	flag.Usage = func() {
		log.Println("Usage: deflate\t\t(reads stdin, outputs to stdout)")
	}
	flag.Parse()
}

func main() {
	parseFlags()

	_, err := io.Copy(os.Stdout, flate.NewReader(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
}
