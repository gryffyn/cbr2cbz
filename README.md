# cbr2cbz
[![Build Status](https://ci.neveris.one/api/badges/gryffyn/cbr2cbz/status.svg)](https://ci.neveris.one/gryffyn/cbr2cbz)  
Converts CBR (`application/vnd.comicbook-rar`) files to CBZ (`application/vnd.comicbook+zip`) files.

## Installation
With go:  
`go install git.neveris.one/gryffyn/cbr2cbz@latest`

Using precompiled binaries:  
Download the binary from `Releases`.

## Usage
```
$ cbr2cbz -h
Usage:
  cbr2cbz [OPTIONS] [<INPUT CBR>] [<OUTPUT CBZ>]

Application Options:
  -r, --rename        Renames files with incorrect extensions

Help Options:
  -h, --help          Show this help message
```

## License
This project is under the MIT license.  
Sections of `ar/rar.go` and `ar/util.go` under MIT from [mholt/archiver](https://github.com/mholt/archiver).