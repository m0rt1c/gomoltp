package moltp

import "fmt"

type (
	// RawFormula object holding a single unparsed formula encoded using a TEX notation
	RawFormula struct {
		OID     int    `json:"oid"`
		Formula string `json:"formula"`
	}

	// RawSequent object holding a single unparsed Sequent
	// Left and right parts are encoded using a TEX notation
	RawSequent struct {
		Left  string `json:"left"`
		Right string `json:"right"`
	}

	// Prover object holding the prover state
	Prover struct {
		debugOn bool
	}

	// Sequent object holding a Sequent
	Sequent struct {
		Name          string
		Justification []string
		Left          []*formula
		Right         []*formula
	}

	token struct {
		Value string // token symbol value
		IsTe  bool   // is terminal
		IsIn  bool   // is an index for a terminal
		IsLB  bool   // is left braket
		IsRB  bool   // is right braket
		IsOp  bool   // is operator
		UnOp  bool   // is unary operator
		BiOp  bool   // is binary oprator
		MuOp  bool   // is miltiple arguments operator
		IsCo  bool   // is comma for multiple args operators
		Skip  int    // how many char was have to be skipped from input
	}

	unification struct{}

	substitution struct{}

	worldsymbol struct {
		Value  string
		Index  int
		Ground bool
	}

	worldindex struct {
		Symbols []worldsymbol
	}

	formula struct {
		Operands []*formula
		Terminal string
		Index    string
	}
)

func (s *Sequent) String() string {
	return fmt.Sprintf("%s: %s <- %s %v",
		s.Name,
		formulaArrayToString(s.Left),
		formulaArrayToString(s.Right),
		s.Justification)
}

func (f *formula) String() string {
	switch len(f.Operands) {
	case 0:
		if len(f.Index) < 1 {
			return fmt.Sprintf("%s", f.Terminal)
		}
		return fmt.Sprintf("%s_{%s}", f.Terminal, f.Index)
	case 1:
		if len(f.Index) < 1 {
			return fmt.Sprintf("( %s %s )", f.Terminal, f.Operands[0])
		}
		return fmt.Sprintf("|( %s %s )|_{%s}", f.Terminal, f.Operands[0], f.Index)
	case 2:
		if len(f.Index) < 1 {
			return fmt.Sprintf("( %s %s %s )", f.Operands[0], f.Terminal, f.Operands[1])
		}
		return fmt.Sprintf("|( %s %s %s )|_{%s}", f.Operands[0], f.Terminal, f.Operands[1], f.Index)
	default:
		k := ""
		for _, o := range f.Operands {
			if k == "" {
				k = fmt.Sprintf("%s", o)
			} else {
				k = fmt.Sprintf("%s, %s", k, o)
			}
		}
		if len(f.Index) < 1 {
			return fmt.Sprintf("( %s %s )", f.Terminal, k)
		}
		return fmt.Sprintf("|( %s %s )|_{%s}", f.Terminal, k, f.Index)
	}
}

func (i *worldindex) ground() bool {
	for _, s := range i.Symbols {
		if !s.Ground {
			return false
		}
	}
	return true
}

func (i *worldindex) start() worldsymbol {
	if len(i.Symbols) < 1 {
		// TODO should we return an error here?
		return worldsymbol{}
	}
	return i.Symbols[0]
}

func (i *worldindex) end() worldsymbol {
	l := len(i.Symbols)
	if l < 1 {
		// TODO should we return an error here?
		return worldsymbol{}
	}
	return i.Symbols[l-1]
}

func (i *worldindex) wunify(j *worldindex) *unification {
	if i.start().Value == "0" && j.start().Value == "0" {
		if i.end().Ground && j.end().Ground && i.end().Value == j.end().Value {
			return &unification{}
		}
		// TODO implement the rest
	}
	return nil
}

func (f *formula) munify(g *formula) *substitution {
	return nil
}

func (g *substitution) applySubstitutionTo(fs []*formula) []*formula {
	return fs
}
