package moltp

import (
	"errors"
	"fmt"
	"log"
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
	}

	worldindex struct {
	}

	formula struct {
		Operator   operator
		Operands   []formula
		WorldIndex worldindex
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

func operatorPreceeds(a, b *token) bool {
	if a.MuOp {
		return false
	}
	if a.UnOp {
		return true
	}
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
	return false
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

func nextToken(s string) *token {
	switch s[0] {
	case '(':
		return &token{IsLB: true, Value: "Round", Skip: 1}
	case '[':
		return &token{IsLB: true, Value: "Square", Skip: 1}
	case '{':
		return &token{IsLB: true, Value: "Curly", Skip: 1}
	case ')':
		return &token{IsRB: true, Value: "Round", Skip: 1}
	case ']':
		return &token{IsRB: true, Value: "Square", Skip: 1}
	case '}':
		return &token{IsRB: true, Value: "Curly", Skip: 1}
	case '\\':
		return matchOperator(s[1], s[2])
	default:
		return &token{IsTe: true, Value: fmt.Sprintf("%c", s[0]), Skip: 1}
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
	t := nextToken(s)
	for t != nil {
		if t.IsTe {
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
					return tokens, errors.New("missing opening brakets")
				}
			}
		}
		offset = offset + t.Skip
		segment := segment[offset:]
		if len(segment) < 1 {
			break
		} else {
			t = nextToken(segment)
		}
	}
	for i := len(ops) - 1; i >= 0; i-- {
		tokens = append(tokens, ops[i])
		// ops = ops[:len(ops)-1]
	}
	return tokens, nil
}

func parseRawFormula(rf RawFormula) (*formula, error) {
	top := &formula{}
	tokens, err := tokenize(strings.Replace(rf.Formula, " ", "", -1), 0x00)
	if err != nil {
		return top, err
	}

	log.Printf("Parsed tokens for %s\n", rf.Formula)
	for i := len(tokens) - 1; i >= 0; i-- {
		log.Printf("%d: %s\n", len(tokens)-i, tokens[i].Value)
	}
	// generateFormulaTree
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
