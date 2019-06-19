package moltp

import (
	"fmt"
	"strconv"
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
		Value string   // token symbol value
		Vars  []string // variable, used for functions and unviversal ops
		IsTe  bool     // is terminal
		IsIn  bool     // is an index for a terminal
		IsLB  bool     // is left braket
		IsRB  bool     // is right braket
		IsOp  bool     // is operator
		UnOp  bool     // is unary operator
		BiOp  bool     // is binary oprator
		MuOp  bool     // is miltiple arguments operator
		IsCo  bool     // is comma for multiple args operators
		Skip  int      // how many char was have to be skipped from input
	}

	unification struct {
		Map map[string]string
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
		Vars     []string
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
		ter := f.Terminal
		if len(f.Vars) > 0 {
			vars := ""
			for _, v := range f.Vars {
				if vars == "" {
					vars = v
				} else {
					vars = fmt.Sprintf("%s,%s", vars, v)
				}
			}
			ter = fmt.Sprintf("%s(%s)", ter, vars)
		}
		if len(f.Index.Symbols) < 1 {
			return fmt.Sprintf("%s", ter)
		}
		return fmt.Sprintf("|%s|_{%s}", ter, &f.Index)
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
	for k, v := range u.Map {
		if out == "" {
			out = fmt.Sprintf("%s/%s", k, v)
		} else {
			out = fmt.Sprintf("%s,%s/%s", out, k, v)
		}
	}
	return fmt.Sprintf("{%s}", out)
}

func compose(m, n *unification) *unification {
	if m == nil {
		return n
	}
	if n == nil {
		return m
	}
	for k, v := range m.Map {
		n.Map[k] = v
	}
	return n
}

func unify(f, g *formula) *unification {
	u := &unification{Map: make(map[string]string)}
	l1 := len(f.Vars)
	l2 := len(g.Vars)
	if l1 < 1 || l2 < 1 {
		return u
	}
	l := l1
	if l2 < l1 {
		l = l2
	}
	// TODO: change this
	// In short we need to substite all non free variables usually not numbers
	// with numbers
	for i := 0; i < l; i++ {
		a := f.Vars[i]
		_, err1 := strconv.Atoi(a)
		b := g.Vars[i]
		_, err2 := strconv.Atoi(b)
		if err1 == nil && err2 != nil {
			u.Map[b] = a
		}
		if err1 != nil && err2 == nil {
			u.Map[a] = b
		}
	}
	return u
}

func (R *relation) munify(f, g *formula) *unification {
	m := unify(f, g)
	n := R.wunify(&f.Index, &g.Index)
	return compose(m, n)
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
	if f.Terminal == sFORALL {
		f.Operands[len(f.Operands)-1] = u.applyUnification(f.Operands[len(f.Operands)-1])
	} else {
		for i, o := range f.Operands {
			f.Operands[i] = u.applyUnification(o)
		}
	}

	t := copyTopFormulaLevel(f)
	changes := false
	if len(t.Operands) == 0 {
		n, ok := u.Map[t.Terminal]
		if ok {
			changes = true
			t.Terminal = n
		}
		for i, v := range t.Vars {
			n, ok := u.Map[v]
			if ok {
				changes = true
				t.Vars[i] = n
			}
		}
	}
	if !changes {
		return f
	}
	return t
}

func (u *unification) applyUnifications(fs []*formula) []*formula {
	for i, f := range fs {
		fs[i] = u.applyUnification(f)
	}
	return fs
}

func (R *relation) findUnification(s0, s1 *worldsymbol) *unification {
	u := &unification{Map: make(map[string]string)}
	// TODO: we need to change this
	_, err := strconv.Atoi(s1.Value)
	if err != nil {
		return u
	}
	u.Map[s0.Value] = s1.Value
	return u
}

func (R *relation) wunify(i, j *worldindex) *unification {
	if start(i).Value == "0" && start(j).Value == "0" {
		if end(i).Ground && end(j).Ground && end(i).Value == end(j).Value {
			return &unification{Map: make(map[string]string)}
		}
		if (end(i).Ground && !end(j).Ground || !end(i).Ground && end(j).Ground) && R.Serial {
			o := R.findUnification(end(i), end(j))
			if o != nil {
				return o
			}
		}
		if !end(i).Ground && !end(j).Ground && R.Serial {
			o := R.findUnification(end(j), end(i))
			if o != nil {
				return o
			}
			o = R.findUnification(end(i), end(j))
			if o != nil {
				return o
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
		if !s.Ground {
			if vars == "" {
				vars = s.Value
			} else {
				vars = fmt.Sprintf("%s,%s", vars, s.Value)
			}
		}
	}
	for _, s := range f.Vars {
		if vars == "" {
			vars = s
		} else {
			vars = fmt.Sprintf("%s,%s", vars, s)
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

// GetAllFreeVars finds free vars in all the subformulas
func (f *formula) GetAllFreeVars(nonFreeVars *map[string]bool) []string {
	if nonFreeVars == nil {
		f := make(map[string]bool)
		nonFreeVars = &f
	}
	sub := []string{}
	for _, o := range f.Operands {
		sub = append(sub, o.GetAllFreeVars(nonFreeVars)...)
	}
	for _, v := range sub {
		_, ok := (*nonFreeVars)[v]
		if !ok {
			sub = append(sub, v)
		}
	}
	return sub
}
