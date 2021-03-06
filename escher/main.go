// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	. "github.com/gocircuit/escher/a"
	. "github.com/gocircuit/escher/be"
	. "github.com/gocircuit/escher/circuit"
	. "github.com/gocircuit/escher/faculty"
	. "github.com/gocircuit/escher/kit/fs"
	kio "github.com/gocircuit/escher/kit/io"
	"github.com/gocircuit/escher/see"

	// Load faculties
	_ "github.com/gocircuit/escher/faculty/basic"
	"github.com/gocircuit/escher/faculty/circuit"
	_ "github.com/gocircuit/escher/faculty/cmplx"
	_ "github.com/gocircuit/escher/faculty/escher"
	_ "github.com/gocircuit/escher/faculty/http"
	_ "github.com/gocircuit/escher/faculty/index"
	_ "github.com/gocircuit/escher/faculty/io"
	_ "github.com/gocircuit/escher/faculty/math"
	_ "github.com/gocircuit/escher/faculty/model"
	fos "github.com/gocircuit/escher/faculty/os"
	"github.com/gocircuit/escher/faculty/test"
	_ "github.com/gocircuit/escher/faculty/text"
	_ "github.com/gocircuit/escher/faculty/time"
	_ "github.com/gocircuit/escher/faculty/yield"
)

// usage: escher [-a dir] [-show] address arguments...
var (
	flagSrc      = flag.String("src", "", "source directory")
	flagDiscover = flag.String("d", "", "multicast UDP discovery address for gocircuit.org faculty")
)

func main() {
	// parse flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [-src Dir] [-d NetAddress] MainCircuit Arguments...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	var flagMain string
	var flagArgs = flag.Args()
	if len(flagArgs) > 0 {
		flagMain, flagArgs = flagArgs[0], flagArgs[1:]
	}
	// parse env
	if *flagSrc == "" {
		*flagSrc = os.Getenv("ESCHER")
	}

	// initialize faculties
	fos.Init(flagArgs)
	test.Init(*flagSrc)
	circuit.Init(*flagDiscover)
	//
	index := Root()
	if *flagSrc != "" {
		index.Merge(Load(*flagSrc))
	}
	// run main
	if flagMain != "" {
		verb := see.ParseVerb(flagMain)
		if Circuit(verb).IsNil() {
			fmt.Fprintf(os.Stderr, "verb not recognized\n")
			os.Exit(1)
		}
		exec(index, Circuit(verb), false)
	}
	// standard loop
	r := kio.NewChunkReader(os.Stdin)
	for {
		chunk, err := r.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "end of session (%v)\n", err)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
		}
		src := NewSrcString(string(chunk))
		for src.Len() > 0 {
			u := see.SeeChamber(src)
			if u == nil || u.(Circuit).Len() == 0 {
				break
			}
			fmt.Fprintf(os.Stderr, "MATERIALIZING %v\n", u)
			exec(index, u.(Circuit), true)
		}
	}
}

func exec(index Index, verb Circuit, showResidue bool) {
	residue := MaterializeSystem(Circuit(verb), Circuit(index), New().Grow("Main", New()))
	if showResidue {
		fmt.Fprintf(os.Stderr, "RESIDUE %v\n\n", residue)
	}
}
