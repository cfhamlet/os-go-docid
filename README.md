# os-go-docid

[![Build Status](https://www.travis-ci.org/cfhamlet/os-go-docid.svg?branch=master)](https://www.travis-ci.org/cfhamlet/os-go-docid)
[![codecov](https://codecov.io/gh/cfhamlet/os-go-docid/branch/master/graph/badge.svg)](https://codecov.io/gh/cfhamlet/os-go-docid)
[![Documentation](https://godoc.org/github.com/cfhamlet/os-go-docid/docid?status.svg)](https://godoc.org/github.com/cfhamlet/os-go-docid/docid)


DocID for Golang. 

Python version is [here]( https://github.com/cfhamlet/os-docid ).

## Install

You can get the library with ``go get``

```
go get -u github.com/cfhamlet/os-go-docid
```

The binary command line tool can be build from source, you should always build the latest released version , [Releases List]( https://github.com/cfhamlet/os-go-docid/releases )

```
git clone -b v0.0.2 https://github.com/cfhamlet/os-go-docid.git
cd os-go-docid
make install
```

## Document

Read the [GoDoc](https://godoc.org/github.com/cfhamlet/os-go-docid/docid ).

## Usage

### APIs

```
import "github.com/cfhamlet/os-go-docid/docid"

url  := "http://www.example.com/"
d, e := docid.New(url)
```

### Command line

```
$ go-docid -h
go-docid development

usage: go-docid [-hvV] [-f file list]

options:
  -V    show verbos version info and exit
  -f value
        comma separated file list (default: stdin)
  -v    show version and exit
```

## License
  MIT licensed.

