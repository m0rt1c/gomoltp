package moltp

import (
	"fmt"
	"strings"
	"sync"
)

type (
	// RawFormula object holding a single unparsed formula encoded using a TEX notation
	RawFormula struct {
		OID     int    `json:"oid"`
		Formula string `json:"formula"`
	}

	// RawSequent object holding a single unparsed sequent
	// Left and right parts are encoded using a TEX notation
	RawSequent struct {
		Left  string `json:"left"`
		Right string `json:"right"`
	}

	sequent struct {
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

	operator interface {
		apply(ops []*formula) bool
	}

	or struct {
	}

	formula struct {
		Operator operator
		Operands []*formula
		Terminal string
		Index    string
	}

	inferenceRule interface {
		applyRuleTo(s *sequent) (*sequent, error)
	}

	r1  struct{}
	r2  struct{}
	r3  struct{}
	r4  struct{}
	r5  struct{}
	r6  struct{}
	r7  struct{}
	r8  struct{}
	r9  struct{}
	r10 struct{}
)

const (
	sBOX     = "Box"
	sDIAMOND = "Diamond"
	sEXISTS  = "Exists"
	sFORALL  = "Forall"
	sIFF     = "Iff"
	sIMPLIES = "Implies"
	sAND     = "And"
	sOR      = "Or"
	sNOT     = "Not"
)

var (
	sEInit    = &sync.Once{}
	sEncoding = make(map[string]string)
	rules     = []inferenceRule{r1{}, r2{}, r3{}, r4{}, r5{}, r6{}, r7{}, r8{}, r9{}, r10{}}
)

// this functions rapresenting inference rules returns
// 1) a sequent and a nil if the rule was applied successfully. The returned sequent is the result of applying the rule
// 2) nil and nil if the sequent was not appliable
// 3) nil and an error if there was some sort of error

// R1: If S,|p|_{i} <- T and S' <- |q|_{j}, T' and |p|_{i} and |q|_{j}
// unify with unification O then S_{O} U S'_{O} <- T_{O} U T'_{O}
func (r r1) applyRuleTo(s *sequent) (*sequent, error) {
	return nil, nil
}

// R2: If S,|(p->q)|_{i} <- T then S,|q|_{i}<-|p|_{i},T
func (r r2) applyRuleTo(s *sequent) (*sequent, error) {
	l := len(s.Left)
	if l < 1 {
		return nil, nil
	}
	f := s.Left[l-1]
	if f.Terminal == sIMPLIES {
		n := &sequent{}

		t := copyTopFormulaLevel(f.Operands[1])
		t.Index = f.Index
		n.Left = append(s.Left[:l-1], t)

		t = copyTopFormulaLevel(f.Operands[0])
		t.Index = f.Index
		n.Right = append([]*formula{t}, s.Right...)

		n.Justification = append(s.Justification, s.Name)
		return n, nil
	}
	return nil, nil
}

// R3: If S <- |(p->q)|_{i},T then S <- |q|_{i},T
func (r r3) applyRuleTo(s *sequent) (*sequent, error) {
	l := len(s.Right)
	if l < 1 {
		return nil, nil
	}
	f := s.Right[0]
	if f.Terminal == sIMPLIES {
		n := &sequent{}

		t := copyTopFormulaLevel(f.Operands[1])
		t.Index = f.Index
		n.Right = append([]*formula{t}, s.Right...)
		n.Left = s.Left

		return n, nil
	}
	return nil, nil
}

// R4: If S <- |(p->q)|_{i},T then S,|p|_{i} <- T
func (r r4) applyRuleTo(s *sequent) (*sequent, error) {
	l := len(s.Right)
	if l < 1 {
		return nil, nil
	}
	f := s.Right[0]
	if f.Terminal == sIMPLIES {
		n := &sequent{}

		t := copyTopFormulaLevel(f.Operands[0])
		t.Index = f.Index
		n.Left = append(s.Left, t)
		n.Right = s.Right

		return n, nil
	}
	return nil, nil
}

func (r r5) applyRuleTo(s *sequent) (*sequent, error) {
	return nil, nil
}

func (r r6) applyRuleTo(s *sequent) (*sequent, error) {
	return nil, nil
}

func (r r7) applyRuleTo(s *sequent) (*sequent, error) {
	return nil, nil
}

func (r r8) applyRuleTo(s *sequent) (*sequent, error) {
	return nil, nil
}

func (r r9) applyRuleTo(s *sequent) (*sequent, error) {
	return nil, nil
}

func (r r10) applyRuleTo(s *sequent) (*sequent, error) {
	return nil, nil
}

// Utility functions

func copyTopFormulaLevel(src *formula) *formula {
	dst := &formula{}

	dst.Operator = src.Operator
	dst.Operands = src.Operands
	dst.Terminal = src.Terminal
	dst.Index = src.Index

	return dst
}

func formulaArrayToString(a []*formula) string {
	out := ""
	for _, f := range a {
		if out == "" {
			out = fmt.Sprintf("%s", f)
		} else {
			out = fmt.Sprintf("%s, %s", a, f)
		}
	}
	return out
}

func (s *sequent) String() string {
	return fmt.Sprintf("%s <- %s", formulaArrayToString(s.Left), formulaArrayToString(s.Right))
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

func (or) apply(ops []*formula) bool {
	return false
}

func operatorPreceeds(a, b *token) bool {
	if a.MuOp {
		return false
	}
	if a.UnOp && !b.UnOp {
		return true
	}
	if a.BiOp && b.BiOp {
		switch a.Value {
		case sOR:
			switch b.Value {
			case sOR:
				return false
			case sAND:
				return true
			case sIMPLIES:
				return true
			case sIFF:
				return true
			}
		case sAND:
			switch b.Value {
			case sOR:
				return false
			case sAND:
				return false
			case sIMPLIES:
				return true
			case sIFF:
				return true
			}
		case sIMPLIES:
			switch b.Value {
			case sOR:
				return false
			case sAND:
				return false
			case sIMPLIES:
				return false
			case sIFF:
				return true
			}
		case sIFF:
			return false
		}
	}
	return false
}

func matchIndex(s string) (*token, error) {
	if s[1] == '{' {
		for j := 2; j < len(s); j++ {
			if s[j] == '}' {
				return &token{IsIn: true, Value: fmt.Sprintf("%s", s[1:j]), Skip: j + 2}, nil
			}
		}
		return nil, fmt.Errorf("missing closing } in index")
	}
	return &token{IsIn: true, Value: fmt.Sprintf("%c", s[1]), Skip: 2}, nil
}

// TODO: Find a better way to init this object
func initsEncoding() {
	sEncoding[sBOX] = "\\Box"
	sEncoding[sDIAMOND] = "\\Diamond"
	sEncoding[sEXISTS] = "\\exists"
	sEncoding[sFORALL] = "\\forall"
	sEncoding[sIFF] = "\\iff"
	sEncoding[sIMPLIES] = "\\to"
	sEncoding[sAND] = "\\land"
	sEncoding[sOR] = "\\lor"
	sEncoding[sNOT] = "\\lnot"
}

func toTextRepr(s string) string {
	sEInit.Do(initsEncoding)
	for k, v := range sEncoding {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

// the len(string) was left here instead of the immediate value to better understand from where the value came
func matchOperator(o, t byte) *token {
	switch o {
	case 'B':
		return &token{IsOp: true, UnOp: true, Value: sBOX, Skip: len("\\Box")}
	case 'D':
		return &token{IsOp: true, UnOp: true, Value: sDIAMOND, Skip: len("\\Diamond")}
	case 'e':
		return &token{IsOp: true, MuOp: true, Value: sEXISTS, Skip: len("\\exists")}
	case 'f':
		return &token{IsOp: true, MuOp: true, Value: sFORALL, Skip: len("\\forall")}
	case 'i':
		return &token{IsOp: true, BiOp: true, Value: sIFF, Skip: len("\\iff")}
	case 't':
		return &token{IsOp: true, BiOp: true, Value: sIMPLIES, Skip: len("\\to")}
	}
	switch t {
	case 'a':
		return &token{IsOp: true, BiOp: true, Value: sAND, Skip: len("\\land")}
	case 'o':
		return &token{IsOp: true, BiOp: true, Value: sOR, Skip: len("\\lor")}
	case 'n':
		return &token{IsOp: true, UnOp: true, Value: sNOT, Skip: len("\\lnot")}
	}
	return nil
}

func genOperator(s string) operator {
	op := or{}
	return op
}

func nextToken(s string) (*token, error) {
	switch s[0] {
	case '(':
		return &token{IsLB: true, Value: "Round", Skip: 1}, nil
	case '[':
		return &token{IsLB: true, Value: "Square", Skip: 1}, nil
	case '{':
		return &token{IsLB: true, Value: "Curly", Skip: 1}, nil
	case ')':
		return &token{IsRB: true, Value: "Round", Skip: 1}, nil
	case ']':
		return &token{IsRB: true, Value: "Square", Skip: 1}, nil
	case '}':
		return &token{IsRB: true, Value: "Curly", Skip: 1}, nil
	case '\\':
		return matchOperator(s[1], s[2]), nil
	case '_':
		return matchIndex(s)
	default:
		return &token{IsTe: true, Value: fmt.Sprintf("%c", s[0]), Skip: 1}, nil
	}
}

// Based on Shunting Yard Algorithm
// 1.  While there are tokens to be read:
// 2.        Read a token
// 3.        If it's a terminal add it to the token queue
// 4.        If it's an operator
// 5.               While there's an operator on the top of the operator stack with greater precedence:
// 6.                       Pop operators from the operator stack onto the token queue
// 7.               Push the current operator onto the operator stack
// 8.        If it's a left bracket push it onto the operators stack
// 9.        If it's a right bracket
// 10.            While there's not a left bracket at the top of the stack:
// 11.                     Pop operators from the stack onto the output queue.
// 12.             Pop the left bracket from the stack and discard it
// 13. While there are operators on the stack, pop them to the queue
func tokenize(s string, term byte) ([]*token, error) {
	var tokens []*token
	var ops []*token
	var offset int
	segment := s
	t, err := nextToken(s)
	if err != nil {
		return tokens, err
	}
	for t != nil {
		if t.IsTe || t.IsIn {
			tokens = append(tokens, t)
		}
		if t.IsOp {
			for len(ops) > 0 {
				k := ops[len(ops)-1]
				if k.IsOp && operatorPreceeds(k, t) {
					ops = ops[:len(ops)-1]
					tokens = append(tokens, k)
				} else {
					break
				}
			}
			ops = append(ops, t)
		}
		if t.IsLB {
			ops = append(ops, t)
		}
		if t.IsRB {
			for len(ops) > 0 {
				k := ops[len(ops)-1]
				if k.IsLB && k.Value == t.Value {
					ops = ops[:len(ops)-1]
					break
				} else {
					ops = ops[:len(ops)-1]
					tokens = append(tokens, k)
				}
				if len(ops) == 0 {
					return tokens, fmt.Errorf("missing opening brakets")
				}
			}
		}
		offset = offset + t.Skip
		segment := segment[offset:]
		if len(segment) < 1 {
			break
		} else {
			t, err = nextToken(segment)
			if err != nil {
				return tokens, err
			}
		}
	}
	for i := len(ops) - 1; i >= 0; i-- {
		tokens = append(tokens, ops[i])
		// ops = ops[:len(ops)-1]
	}
	return tokens, nil
}

func genFormulasTree(tokens []*token) (*formula, error) {
	var formulas []*formula
	for _, t := range tokens {
		if t.IsOp {
			if t.BiOp {
				if len(formulas) < 2 {
					return formulas[0], fmt.Errorf("missing argument for binary operator %s", t.Value)
				}
				f := &formula{}
				f.Terminal = t.Value
				f.Operator = genOperator(t.Value)
				f.Operands = append(f.Operands, formulas[len(formulas)-2:]...)
				formulas = formulas[:len(formulas)-2]
				formulas = append(formulas, f)
			}
			if t.UnOp {
				if len(formulas) < 1 {
					return formulas[0], fmt.Errorf("missing argument for unary operator %s", t.Value)
				}
				f := &formula{}
				f.Terminal = t.Value
				f.Operator = genOperator(t.Value)
				f.Operands = append(f.Operands, formulas[len(formulas)-1:]...)
				formulas = formulas[:len(formulas)-1]
				formulas = append(formulas, f)
			}
		}
		if t.IsTe {
			formulas = append(formulas, &formula{Terminal: t.Value})
		}
		if t.IsIn {
			if len(formulas) < 1 {
				return formulas[0], fmt.Errorf("trying to assign idex %s to nothing", t.Value)
			}
			formulas[len(formulas)-1].Index = t.Value
		}
	}
	return formulas[0], nil
}

func encodeSequent(s *sequent) (*RawSequent, error) {
	rs := &RawSequent{}

	rs.Left = toTextRepr(formulaArrayToString(s.Left))
	rs.Right = toTextRepr(formulaArrayToString(s.Right))

	return rs, nil
}

func proveFormula(f *formula) (*map[int]*sequent, error) {
	i := 0
	solution := make(map[int]*sequent)
	unreduced := []*sequent{}
	f.Index = "0"
	unreduced = append(unreduced, &sequent{Right: []*formula{f}, Name: "S0"})
	for len(unreduced) > 0 {
		ruleWasApplied := false
		// Try to apply each rule
		for _, rule := range rules {
			last := unreduced[len(unreduced)-1]
			s, err := rule.applyRuleTo(last)
			if err != nil {
				return &solution, err
			}
			if s != nil {
				// The rule was applied successfully
				i = i + 1
				s.Name = fmt.Sprintf("S%d", i)
				solution[i] = last
				if len(s.Left) == 0 && len(s.Right) == 0 {
					// A solution was found
					i = i + 1
					solution[i] = s
					return &solution, nil
				}
				unreduced = append(unreduced, s)
				ruleWasApplied = true
			}
			// else the rule was not appliable
		}
		if !ruleWasApplied {
			return &solution, fmt.Errorf("no rule was applied")
		}
	}
	return &solution, nil
}

// Prove givent a set of formulas it output a solution, if debugOn is true debugging messages will be printed
func Prove(rf *RawFormula, debugOn bool) (*map[int]*RawSequent, error) {
	if debugOn {
		fmt.Printf("Input:\n%s\n", rf.Formula)
	}
	tokens, err := tokenize(strings.Replace(rf.Formula, " ", "", -1), 0x00)
	if debugOn {
		fmt.Println("Tokens:")
		for i := len(tokens) - 1; i >= 0; i-- {
			fmt.Printf("%d: %s\n", len(tokens)-i, tokens[i].Value)
		}
	}
	if err != nil {
		return nil, err
	}
	top, err := genFormulasTree(tokens)
	if debugOn {
		fmt.Println("Formula:")
		fmt.Println(top)
	}
	if err != nil {
		return nil, err
	}
	s, err := proveFormula(top)
	if debugOn {
		fmt.Println("Sequents:")
		for key := 0; key < len(*s); key++ {
			sequent := (*s)[key]
			fmt.Printf("%s: %s <- %s %v.\n", sequent.Name, formulaArrayToString(sequent.Left), formulaArrayToString(sequent.Right), sequent.Justification)
		}
	}
	if err != nil {
		return nil, err
	}
	rawSolution := make(map[int]*RawSequent)
	for key, sequent := range *s {
		rs, err := encodeSequent(sequent)
		if err != nil {
			return &rawSolution, nil
		}
		rawSolution[key] = rs
	}
	return &rawSolution, nil
}
