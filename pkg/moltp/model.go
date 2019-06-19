package moltp

import (
	"fmt"
	"strings"
)

type (
	// RawFormula object holding a single unparsed formula encoded using a TEX notation
	RawFormula struct {
		OID     int    `json:"oid"`
		Formula string `json:"formula"`
	}

	// RawSequent object holding a single unparsed Sequent
	// Left and right parts are encoded using a TEX notation
	RawSequent struct {
		Name          string `json:"name"`
		Left          string `json:"left"`
		Right         string `json:"right"`
		Justification string `json:"just"`
	}

	// Prover object holding the prover state
	Prover struct {
		Debug          bool
		Rules          []inferenceRule
		ResolutionRule resolutionRule
		R              *relation
		worldsKeeper   *worldskeeper
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
		Map map[worldsymbol]worldsymbol
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
		Ground bool
	}

	worldindex struct {
		Symbols []*worldsymbol
	}

	worldskeeper struct {
		nextConst    int
		nextVar      string
		nextFunction string
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
		// Workaround when multioperators have 2 arguments
		// TODO: Find a better way to handle this case
		if f.Terminal != sFORALL && f.Terminal != sEXISTS {
			if len(f.Index.Symbols) < 1 {
				return fmt.Sprintf("( %s %s %s )", f.Operands[0], f.Terminal, f.Operands[1])
			}
			return fmt.Sprintf("|( %s %s %s )|_{%s}", f.Operands[0], f.Terminal, f.Operands[1], &f.Index)
		}
		// going to default
		fallthrough
	default:
		// In multi operator formulas ( forall ) the first n elements are variables name
		// the last one a formula
		k := ""
		for _, o := range f.Operands[:len(f.Operands)-1] {
			if k == "" {
				k = fmt.Sprintf("%s", o)
			} else {
				k = fmt.Sprintf("%s, %s", k, o)
			}
		}
		k = fmt.Sprintf("( %s ) %s", k, f.Operands[len(f.Operands)-1])

		if len(f.Index.Symbols) < 1 {
			return fmt.Sprintf("( %s %s )", f.Terminal, k)
		}
		return fmt.Sprintf("|( %s %s )|_{%s}", f.Terminal, k, &f.Index)
	}
}

func (s *worldsymbol) String() string {
	return s.Value
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

func (u *unification) String() string {
	out := ""
	if len(u.Map) > 0 {
		for k, v := range u.Map {
			out = fmt.Sprintf("%s/%s,", &k, &v)
		}
	}
	return fmt.Sprintf("{%s}", strings.TrimSuffix(out, ","))
}

func (R *relation) munify(f, g *formula) *unification {
	n := R.wunify(&f.Index, &g.Index)
	if n != nil {
		return n
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

func end(i *worldindex) *worldsymbol {
	if len(i.Symbols) < 1 {
		return nil
	}
	return i.Symbols[0]
}

func start(i *worldindex) *worldsymbol {
	l := len(i.Symbols)
	if l < 1 {
		return nil
	}
	return i.Symbols[l-1]
}

func (p *Prover) initProver() {
	if p.R == nil {
		p.R = &relation{Serial: true}
	}
	if p.worldsKeeper == nil {
		p.worldsKeeper = &worldskeeper{nextVar: "w", nextConst: 0, nextFunction: "f"}
	}
	if len(p.Rules) == 0 {
		// TODO make this look better
		p.Rules = []inferenceRule{
			r2{Name: "R2"},
			r3{Name: "R3"},
			r4{Name: "R4"},
			r5{Name: "R5"},
			r6{Name: "R6"},
			r7{Name: "R7", worldsKeeper: p.worldsKeeper},
			r8{Name: "R8", worldsKeeper: p.worldsKeeper},
			r9{Name: "R9", worldsKeeper: p.worldsKeeper},
			r10{Name: "R10", worldsKeeper: p.worldsKeeper},
		}
	}
	if p.ResolutionRule == nil {
		p.ResolutionRule = r1{Name: "R1", R: p.R}
	}
}

func (u *unification) applyUnification(f *formula) *formula {
	newSymbols := []*worldsymbol{}
	for _, s := range f.Index.Symbols {
		newSymbols = append(newSymbols, &worldsymbol{Value: u.Map[*s].Value, Ground: u.Map[*s].Ground})
	}
	return f
}

func (u *unification) applyUnifications(fs []*formula) []*formula {
	for _, f := range fs {
		u.applyUnification(f)
	}
	return fs
}

func (s *substitution) compose(u *unification) *unification {
	u.Map[*s.Old] = *s.New
	return u
}

func (R *relation) findUnification(s0, s1 *worldsymbol) *unification {
	u := &unification{Map: make(map[worldsymbol]worldsymbol)}
	return u
}

func (R *relation) wunify(i, j *worldindex) *unification {
	if start(i).Value == "0" && start(j).Value == "0" {
		if end(i).Ground && end(j).Ground && end(i).Value == end(j).Value {
			return &unification{Map: make(map[worldsymbol]worldsymbol)}
		}
		if (end(i).Ground && !end(j).Ground || !end(i).Ground && end(j).Ground) && R.Serial {
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

func (k *worldskeeper) GetFreeIndividualConstant() *worldsymbol {
	old := fmt.Sprintf("%d", k.nextConst)
	k.nextConst = k.nextConst + 1
	return &worldsymbol{Value: old, Ground: true}
}

func (k *worldskeeper) GetSkolemFunctionOf(f *formula) *worldsymbol {
	old := k.nextFunction
	switch k.nextFunction[0] {
	case 'f':
		k.nextFunction = "g"
	case 'g':
		k.nextFunction = "h"
	case 'h':
		k.nextFunction = "f'"
	default:
		k.nextFunction = fmt.Sprintf("%s'", k.nextFunction)
	}
	for i := 0; i < len(old)-1; i++ {
		k.nextFunction = k.nextFunction + "'"
	}
	// TODO: Implement corret world index value
	vars := ""
	for _, s := range f.Index.Symbols {
		if s.Ground {
			if vars == "" {
				vars = s.Value
			} else {
				vars = fmt.Sprintf("%s,%s", vars, s.Value)
			}
		}
	}
	return &worldsymbol{Value: fmt.Sprintf("%s(%s)", old, vars), Ground: true}
}

func (k *worldskeeper) GetWorldVariable() *worldsymbol {
	old := k.nextVar
	switch k.nextVar[0] {
	case 'w':
		k.nextVar = "v"
	case 'v':
		k.nextVar = "u"
	case 'u':
		k.nextVar = "w'"
	default:
		k.nextVar = fmt.Sprintf("%s'", k.nextVar)
	}
	for i := 0; i < len(old)-1; i++ {
		k.nextVar = k.nextVar + "'"
	}
	return &worldsymbol{Value: old, Ground: false}
}
