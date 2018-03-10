package main

import (
	"compress/flate"
	"io"
	"log"
	"os"
)

func init() { log.SetFlags(0) }

func main() {
	_, err := io.Copy(os.Stdout, flate.NewReader(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
}
