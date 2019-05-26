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
	var rf moltp.RawFormula
	rf.OID = 0
	rf.Formula = formula
	solution, err := moltp.Prove(rf, debugOn)
	if err != nil {
		log.Println(err)
	} else {
		for key, sequent := range *solution {
			fmt.Printf("%d: %s <- %s\n", key, sequent.Left, sequent.Right)
		}
	}
}
