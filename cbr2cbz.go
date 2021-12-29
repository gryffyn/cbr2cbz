package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"git.neveris.one/gryffyn/cbr2cbz/ar"
	"github.com/jessevdk/go-flags"
)

func main() {
	type Opts struct {
		Rename         bool   `short:"r" long:"rename" description:"Renames files with incorrect extensions"`
		Positional     struct {
			CBR string `positional-arg-name:"<INPUT CBR>" required:"true"`
			CBZ string `positional-arg-name:"<OUTPUT CBZ>"`
		} `positional-args:"yes"`
	}

	var opts Opts
	_, err := flags.Parse(&opts)

	if e, ok := err.(*flags.Error); ok {
		if e.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	outfile := opts.Positional.CBZ
	if opts.Positional.CBZ == "" {
		outfile = fmt.Sprintf("%s.%s", strings.TrimSuffix(opts.Positional.CBR, path.Ext(opts.Positional.CBR)), "cbz")
	}

	rt := ar.Rar{}
	filetype, err := rt.CheckRAR(opts.Positional.CBR)
	if opts.Rename {
		if filetype == "zip" {
			fterr := os.Rename(opts.Positional.CBR, outfile)
			ErrExit(fterr)
		} else {
			ErrExit(err)
		}
	} else {
		ErrExit(err)

		of, err := os.Create(outfile)
		ErrExit(err)
		rz := ar.Zip{}
		rz.Open(of)

		err = rt.Walk(opts.Positional.CBR, func(f ar.File) error {
			err := rz.WriteFile(f)
			return err
		})
		if err != nil {
			log.Fatal(err)
		}

		rz.Close()
		of.Close()
	}

	rt.Close()
}

func ErrExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}