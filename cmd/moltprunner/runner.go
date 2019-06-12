package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gomoltp/pkg/moltp"
)

var (
	debugOn bool
	formula string
)

func init() {
	flag.StringVar(&formula, "f", "\\Box ( a \\to b ) \\to  ( \\Box a \\to \\Box b )", "Formula to be solved.")
	flag.BoolVar(&debugOn, "v", false, "Swith for log printing")
}

func main() {
	flag.Parse()
	rf := &moltp.RawFormula{OID: 0, Formula: formula}
	prover := moltp.Prover{Debug: debugOn}
	solution, err := prover.Prove(rf)
	if err != nil {
		log.Println(err)
	} else {
		for _, s := range solution {
			fmt.Printf("%s\n", s)
		}
	}
}
