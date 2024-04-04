package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func parseArgs(args []string) (conf config, _ error) {
	const usageHead = "Usage: xflate\t(reads stdin, outputs to stdout)"
	var (
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
		return conf, err
	}
	if len(flag.Args()) > 0 {
		return conf, fmt.Errorf("unexpected args, use stdin/stdout")
	}
	if help {
		conf.help = func() {
			flag.SetOutput(os.Stdout)
			flag.Usage()
		}
		return conf, nil
	}

	if flag.Changed("compress") && flag.Changed("decompress") &&
		conf.compress == decompress {
		return conf, fmt.Errorf("conflicting flags -z and -d")
	}
	if flag.Changed("decompress") {
		conf.compress = !decompress
	}

	return conf, nil
}
