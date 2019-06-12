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
