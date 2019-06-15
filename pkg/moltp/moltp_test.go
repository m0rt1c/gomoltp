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

func TestProver1(t *testing.T) {
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
			"S4: |a|_{W0:0} <-  [R8 S3]",
			"S6:  <- |a|_{2:1:0} [R7 S5]",
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

func TestProver2(t *testing.T) {
	rf := &RawFormula{OID: 0, Formula: "\\Box \\Box a \\to \\Diamond \\Diamond a"}
	prover := Prover{Debug: false}
	solution, err := prover.Prove(rf)
	if err != nil {
		t.Errorf("got error %s want nil", err)
	} else {
		out := []string{
			"S1:  <- |( ( Box ( Box a ) ) Implies ( Not ( Box ( Not ( Not ( Box ( Not a ) ) ) ) ) ) )|_{0} []",
			"S3: |( Box ( Box a ) )|_{0} <-  [R4 S1]",
			"S4: |( Box a )|_{W0:0} <-  [R8 S3]",
			"S2:  <- |( Not ( Box ( Not ( Not ( Box ( Not a ) ) ) ) ) )|_{0} [R3 S1]",
			"S6: |( Box ( Not ( Not ( Box ( Not a ) ) ) ) )|_{0} <-  [R6 S2]",
			"S7: |( Box ( Not a ) )|_{V0:0} <-  [R8 S6]",
			"S8:  <- |( Not ( Box ( Not a ) ) )|_{V0:0} [R5 S7]",
			"S9: |( Box ( Not a ) )|_{V0:0} <-  [R6 S8]",
			"S10: |( Not a )|_{V1:V0:0} <-  [R8 S9]",
			"S5: |a|_{W1:W0:0} <-  [R8 S4]",
			"S11:  <- |a|_{V1:V0:0} [R5 S10]",
			"S12:  <-  [R1 S5 S11 {W1/V1}]",
		}
		for i, o := range out {
			s := fmt.Sprintf("%s", solution[i])
			if o != s {
				t.Errorf("got %s want %s", s, o)
			}
		}

	}
}

func TestProver3(t *testing.T) {
	rf := &RawFormula{OID: 0, Formula: "\\Diamond \\Box a \\to \\Box \\Diamond a"}
	prover := Prover{Debug: false}
	solution, err := prover.Prove(rf)
	if err != nil {
		t.Errorf("got error %s want nil", err)
	} else {
		out := []string{
			"S1:  <- |( ( Not ( Box ( Not ( Box a ) ) ) ) Implies ( Box ( Not ( Box ( Not a ) ) ) ) )|_{0} []",
			"S3: |( Box a )|_{1:0} <-  [R4 S1]",
			"S4:  <- |( Box ( Not ( Box a ) ) )|_{0} [R5 S3]",
			"S5:  <- |( Not ( Box a ) )|_{1:0} [R7 S4]",
			"S6: |( Box a )|_{1:0} <-  [R6 S5]",
			"S2:  <- |( Box ( Not ( Box ( Not a ) ) ) )|_{0} [R3 S1]",
			"S8:  <- |( Not ( Box ( Not a ) ) )|_{2:0} [R7 S2]",
			"S9: |( Box ( Not a ) )|_{2:0} <-  [R6 S8]",
			"S10: |( Not a )|_{V0:2:0} <-  [R8 S9]",
			"S7: |a|_{W0:1:0} <-  [R8 S6]",
			"S11:  <- |a|_{V0:2:0} [R5 S10]",
			"S12:  <-  [R1 S7 S11 {W0/V0}]",
		}
		for i, o := range out {
			s := fmt.Sprintf("%s", solution[i])
			if o != s {
				t.Errorf("got %s want %s", s, o)
			}
		}

	}
}
