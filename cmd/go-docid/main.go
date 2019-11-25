package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cfhamlet/os-go-docid/cmd/go-docid/internal/version"
	"github.com/cfhamlet/os-go-docid/docid"
)

type fileList []string

func (f *fileList) Set(v string) error {
	*f = fileList(strings.Split(v, ","))
	return nil
}

func (f *fileList) String() string {
	return strings.Join(*f, ",")
}

var (
	h     bool
	v     bool
	vv    bool
	flist fileList
)

func init() {
	flag.Var(&flist, "f", "comma separated file list (default: stdin)")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&vv, "V", false, "show verbos info and exit")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "go-docid %s\n\n", version.Version())
	fmt.Fprintf(os.Stderr, "usage: go-docid [-hvV] [-f file list]\n\n")
	fmt.Fprintf(os.Stderr, "options:\n")
	flag.PrintDefaults()
}

func readAndProcess(f *os.File) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		b := input.Bytes()
		d, err := docid.New(b)
		if err != nil {
			os.Stdout.Write([]byte("E\t"))
		} else {
			os.Stdout.Write([]byte(d.String()))
			os.Stdout.Write([]byte("\t"))
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
	}

}

func main() {
	flag.Parse()
	if v {
		fmt.Fprintf(os.Stderr, "go-docid %s\n", version.Version())
		os.Exit(0)
	} else if vv {
		fmt.Fprintln(os.Stderr, version.VerbosInfo())
		os.Exit(0)
	}
	if len(flist) <= 0 {
		readAndProcess(os.Stdin)
	} else {
		for _, file := range flist {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warm: %s %v\n", file, err)
				continue
			}
			readAndProcess(f)
			f.Close()
		}
	}
}
