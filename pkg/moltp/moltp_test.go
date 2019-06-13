package moltp

import (
	"fmt"
	"strings"
	"testing"
)

func TestReduceORFormula(t *testing.T) {
	A := &formula{Terminal: "A"}
	B := &formula{Terminal: "B"}
	f := &formula{Terminal: sOR, Operands: []*formula{A, B}}
	out := reduceFormulas(f)

	g0 := &formula{Terminal: sNOT, Operands: []*formula{A}}
	g1 := &formula{Terminal: sIMPLY, Operands: []*formula{g0, B}}

	if strings.Compare(fmt.Sprint(out), fmt.Sprint(g1)) != 0 {
		t.Errorf("got %s want %s", fmt.Sprint(out), fmt.Sprint(g1))
	}
}

func TestReduceANDFormula(t *testing.T) {
	A := &formula{Terminal: "A"}
	B := &formula{Terminal: "B"}
	f := &formula{Terminal: sAND, Operands: []*formula{A, B}}
	out := reduceFormulas(f)

	g0 := &formula{Terminal: sNOT, Operands: []*formula{B}}
	g1 := &formula{Terminal: sIMPLY, Operands: []*formula{A, g0}}
	g3 := &formula{Terminal: sNOT, Operands: []*formula{g1}}

	if strings.Compare(fmt.Sprint(out), fmt.Sprint(g3)) != 0 {
		t.Errorf("got %s want %s", fmt.Sprint(out), fmt.Sprint(g3))
	}
}

func TestProver(t *testing.T) {
	rf := &RawFormula{OID: 0, Formula: "\\Box a \\to \\Box \\Box a"}
	prover := Prover{Debug: false}
	solution, err := prover.Prove(rf)
	if err != nil {
		t.Errorf("got error %s want nil", err)
	} else {
		out := []string{
			"S1:  <- |( ( Box a ) Implies ( Box ( Box a ) ) )|_{0} []",
			"S3: |( Box a )|_{0} <-  [R4 S1]",
			"S2:  <- |( Box ( Box a ) )|_{0} [R3 S1]",
			"S5:  <- |( Box a )|_{1:0} [R7 S2]",
			"S7:  <-  [R1 S4 S6 {W0/2}]",
		}
		for i, o := range out {
			s := fmt.Sprintf("%s", solution[i])
			if o != s {
				t.Errorf("got %s want %s", s, o)
			}
		}

	}
}
