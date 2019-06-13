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
		Debug bool
		Rules []inferenceRule
		R     *relation
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

	unification struct {
		List []substitution
	}

	substitution struct {
		Old *worldsymbol
		New *worldsymbol
	}

	relation struct {
		Serial bool
	}

	worldsymbol struct {
		Value  string
		Index  int
		Ground bool
	}

	worldindex struct {
		Symbols []*worldsymbol
	}

	formula struct {
		Operands []*formula
		Terminal string
		Index    worldindex
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
		if len(f.Index.Symbols) < 1 {
			return fmt.Sprintf("%s", f.Terminal)
		}
		return fmt.Sprintf("|%s|_{%s}", f.Terminal, &f.Index)
	case 1:
		if len(f.Index.Symbols) < 1 {
			return fmt.Sprintf("( %s %s )", f.Terminal, f.Operands[0])
		}
		return fmt.Sprintf("|( %s %s )|_{%s}", f.Terminal, f.Operands[0], &f.Index)
	case 2:
		if len(f.Index.Symbols) < 1 {
			return fmt.Sprintf("( %s %s %s )", f.Operands[0], f.Terminal, f.Operands[1])
		}
		return fmt.Sprintf("|( %s %s %s )|_{%s}", f.Operands[0], f.Terminal, f.Operands[1], &f.Index)
	default:
		k := ""
		for _, o := range f.Operands {
			if k == "" {
				k = fmt.Sprintf("%s", o)
			} else {
				k = fmt.Sprintf("%s, %s", k, o)
			}
		}
		if len(f.Index.Symbols) < 1 {
			return fmt.Sprintf("( %s %s )", f.Terminal, k)
		}
		return fmt.Sprintf("|( %s %s )|_{%s}", f.Terminal, k, &f.Index)
	}
}

func (s *worldsymbol) String() string {
	if s.Ground {
		return s.Value
	}
	return fmt.Sprintf("%s_%d", s.Value, s.Index)
}

func (i *worldindex) String() string {
	switch len(i.Symbols) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%s", i.Symbols[0])
	default:
		out := fmt.Sprintf("%s", i.Symbols[0])
		for _, k := range i.Symbols[1:] {
			out = fmt.Sprintf("%s:%s", out, k)
		}
		return out
	}
}

func unify(f, g *formula) *unification {
	return nil
}

func (R *relation) munify(f, g *formula) *substitution {
	o := unify(f, g)
	if o != nil {
		n := R.wunify(&f.Index, &g.Index)
		if n != nil {
			return &substitution{}
		}
	}
	return nil
}

func (i *worldindex) parent(s *worldsymbol) *worldsymbol {
	for k, p := range i.Symbols {
		if p == s {
			if k < len(i.Symbols)+1 {
				return i.Symbols[len(i.Symbols)+1]
			}
			return nil
		}
	}
	return nil
}

func (i *worldindex) parentIndex(s *worldsymbol) []*worldsymbol {
	for k, p := range i.Symbols {
		if p == s {
			return i.Symbols[k:]
		}
	}
	return []*worldsymbol{}
}

func (i *worldindex) isGround() bool {
	for _, s := range i.Symbols {
		if !s.Ground {
			return false
		}
	}
	return true
}

func start(i *worldindex) *worldsymbol {
	if len(i.Symbols) < 1 {
		return nil
	}
	return i.Symbols[0]
}

func end(i *worldindex) *worldsymbol {
	l := len(i.Symbols)
	if l < 1 {
		return nil
	}
	return i.Symbols[l-1]
}

func (p *Prover) initRules() {
	if p.R == nil {
		p.R = &relation{Serial: true}
	}
	if len(p.Rules) == 0 {
		// TODO make this look better
		p.Rules = []inferenceRule{r1{Name: "R1", R: p.R}, r2{Name: "R2", R: p.R}, r3{Name: "R3", R: p.R}, r4{Name: "R4", R: p.R}, r5{Name: "R5", R: p.R}, r6{Name: "R6", R: p.R}, r7{Name: "R7", R: p.R}, r8{Name: "R8", R: p.R}, r9{Name: "R9", R: p.R}, r10{Name: "R10", R: p.R}}
	}
}

func (s *substitution) applySubstitutionTo(fs []*formula) []*formula {
	return fs
}

func (s *substitution) compose(u *unification) *unification {
	return u
}

func (R *relation) findUnification(s0, s1 *worldsymbol) *unification {
	return &unification{}
}

func (R *relation) wunify(i, j *worldindex) *unification {
	if start(i).Value == "0" && start(j).Value == "0" {
		if end(i).Ground && end(j).Ground && end(i).Value == end(j).Value {
			return &unification{}
		}
		if end(i).Ground && !end(j).Ground && R.Serial {
			o := R.findUnification(end(j), end(i))
			if o != nil {
				s := &substitution{Old: end(i), New: end(j)}
				return s.compose(o)
			}
		}
		if !end(i).Ground && !end(j).Ground && R.Serial {
			o := R.findUnification(end(j), end(i))
			if o != nil {
				s := &substitution{Old: end(i), New: end(j)}
				return s.compose(o)
			}
			o = R.findUnification(end(i), end(j))
			if o != nil {
				s := &substitution{Old: end(i), New: end(j)}
				return s.compose(o)
			}
		}
	}
	return nil
}
