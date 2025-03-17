// xflate compresses or decompresses deflate stream.
package main

import (
	"compress/flate"
	"fmt"
	"io"
	"os"
)

const fileExt = ".deflate"

type action struct {
	fileIn        string
	fileOut       string
	compress      bool
	compressLevel int
	force         bool

	help func()
}

var defaultAction = action{
	compress:      true,
	compressLevel: 6,
}

func openIn(path string) (*os.File, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func openOut(path string, force bool) (*os.File, error) {
	if path == "-" {
		return os.Stdout, nil
	}
	var wflag int
	if force {
		wflag = os.O_TRUNC
	} else {
		wflag = os.O_EXCL
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|wflag, 0644)
}

func run(a action) (rerr error) {
	in, err := openIn(a.fileIn)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := openOut(a.fileOut, a.force)
	if err != nil {
		return err
	}
	defer safeClose(out, &rerr)

	switch {
	case a.compress:
		w, err := flate.NewWriter(out, a.compressLevel)
		if err != nil {
			return fmt.Errorf("failed creating compress writer: %w", err)
		}

		_, err = io.Copy(w, in)
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("compress closing %w", err)
		}

	case !a.compress:
		r := flate.NewReader(in)
		defer safeClose(r, &rerr)

		_, err := io.Copy(out, r)
		if err != nil {
			return fmt.Errorf("decompress: %w", err)
		}
	}

	return nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if conf.help != nil {
		conf.help()
		os.Exit(0)
	}

	err = run(conf)
	if err != nil {
		die(1, err)
	}

}

func safeClose(c io.Closer, errp *error) {
	cerr := c.Close()
	if cerr != nil && *errp == nil {
		*errp = cerr
	}
}

func die(code int, err error) {
	fmt.Fprintln(os.Stderr, "xflate:", err)
	os.Exit(code)
}
