package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/dnesting/clocker"
)

var (
	noTimestamp bool
	noDelta     bool
	overwrite   bool
)

func main() {
	flag.BoolVar(&noTimestamp, "no-timestamp", false, "Do not add timestamp to output")
	flag.BoolVar(&noTimestamp, "T", false, "Do not add timestamp to output")
	flag.BoolVar(&noDelta, "no-delta", false, "Do not add delta to output")
	flag.BoolVar(&noDelta, "D", false, "Do not add delta to output")
	flag.BoolVar(&noDelta, "overwrite", false, "Overwrite the given filename instead of append")
	flag.BoolVar(&noDelta, "o", false, "Overwrite the given filename instead of append")
	flag.Parse()

	var out io.WriteCloser
	switch len(flag.Args()) {
	case 0:
		// use Stdout
	case 1:
		var err error
		if overwrite {
			out, err = os.Create(flag.Args()[0])
		} else {
			out, err = os.OpenFile(flag.Args()[0], os.O_APPEND|os.O_WRONLY, 0644)
		}
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	default:
		log.Println("Too many arguments")
		flag.Usage()
		os.Exit(1)
	}

	clocker.Run(clocker.Options{
		NoTimestamp: noTimestamp,
		NoDelta:     noDelta,
		Out:         out,
	})
}
