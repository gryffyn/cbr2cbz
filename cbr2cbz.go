package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"git.gryffyn.io/gryffyn/cbr2cbz/ar"
	"git.gryffyn.io/gryffyn/cbr2cbz/tweaks"
	"github.com/jessevdk/go-flags"
)

func main() {
	type Opts struct {
		Scale   int `short:"s" long:"scale" description:"Scales images to percent (1-100)"`
		Quality int `short:"q" long:"quality" description:"JPEG quality setting (1-100)"`
		Positional     struct {
			CBR string `positional-arg-name:"<INPUT CBR>" required:"true"`
			CBZ string `positional-arg-name:"<OUTPUT CBZ>"`
		} `positional-args:"yes"`
	}

	var opts Opts
	_, err := flags.Parse(&opts)

	// sets default JPEG quality to 91
	if opts.Quality == 0 {opts.Quality = 91}

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
	if filetype == "zip" {
		fterr := os.Rename(opts.Positional.CBR, outfile)
		ErrExit(fterr)
	} else {
		ErrExit(err)

		of, err := os.Create(outfile)
		ErrExit(err)
		rz := ar.Zip{}
		rz.Open(of)

		err = rt.Walk(opts.Positional.CBR, func(f ar.File) error {
			if opts.Scale > 0 && opts.Scale <= 100 {
				dest, _ := ioutil.TempFile("", "cbr2cbz")
				tweaks.ResizeImage(&f.ReadCloser, dest, opts.Scale, opts.Quality)
			} else if opts.Quality > 0 && opts.Quality <= 100 && opts.Scale == 0 {
				dest, _ := ioutil.TempFile("", "cbr2cbz")
				tweaks.ImgQuality(&f.ReadCloser, dest, opts.Quality)
			}

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

