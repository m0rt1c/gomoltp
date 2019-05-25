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

	// RawSequent object holding a single unparsed sequent
	// Left and right parts are encoded using a TEX notation
	RawSequent struct {
		OID   int    `json:"oid"`
		Left  string `json:"left"`
		Right string `json:"right"`
	}

	sequent struct {
		OID   int
		Left  []formula
		Right []formula
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

func (f *formula) printFormula() {
	switch len(f.Operands) {
	case 0:
		if len(f.Index) < 0 {
			fmt.Printf("%s ", f.Terminal)
		} else {
			fmt.Printf("%s%s ", f.Terminal, f.Index)
		}
	case 1:
		fmt.Printf("%s ", f.Terminal)
		f.Operands[0].printFormula()
	case 2:
		f.Operands[0].printFormula()
		fmt.Printf("%s ", f.Terminal)
		f.Operands[1].printFormula()
	default:
		fmt.Printf("%s ", f.Terminal)
		for _, o := range f.Operands {
			o.printFormula()
		}
	}
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
				return &token{IsIn: true, Value: fmt.Sprintf("_%s", s[1:j+1]), Skip: j + 2}, nil
			}
		}
		return nil, fmt.Errorf("missing closing } in index")
	}
	return &token{IsIn: true, Value: fmt.Sprintf("_%c", s[1]), Skip: 2}, nil
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

func (or) apply(ops []*formula) bool {
	return false
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

func parseRawFormula(rf RawFormula) (*formula, error) {
	top := &formula{}
	tokens, err := tokenize(strings.Replace(rf.Formula, " ", "", -1), 0x00)
	if err != nil {
		return top, err
	}

	fmt.Printf("Parsed tokens for %s\n", rf.Formula)
	for i := len(tokens) - 1; i >= 0; i-- {
		fmt.Printf("%d: %s\n", len(tokens)-i, tokens[i].Value)
	}

	top, err = genFormulasTree(tokens)
	fmt.Printf("Parsed formula for %s\n", rf.Formula)
	top.printFormula()

	return top, err
}

func encodeSequent(s sequent) (RawSequent, error) {
	return RawSequent{}, nil
}

func proveFormula(f *formula) ([]sequent, error) {
	return []sequent{}, nil
}

// Prove givent a set of formulas it output a solution
func Prove(rf RawFormula) ([]RawSequent, error) {
	var solution []RawSequent
	f, err := parseRawFormula(rf)
	if err != nil {
		return solution, err
	}
	s, err := proveFormula(f)
	if err != nil {
		return solution, err
	}
	for _, sequent := range s {
		rs, err := encodeSequent(sequent)
		if err != nil {
			return solution, nil
		}
		solution = append(solution, rs)
	}
	return solution, nil
}
