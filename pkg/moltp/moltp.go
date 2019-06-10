package moltp

import (
	"fmt"
	"strings"
	"sync"
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
)

// Utility functions

func copyTopFormulaLevel(src *formula) *formula {
	dst := &formula{}

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

func encodeSequent(s *Sequent) (RawSequent, error) {
	rs := RawSequent{}

	rs.Left = toTextRepr(formulaArrayToString(s.Left))
	rs.Right = toTextRepr(formulaArrayToString(s.Right))

	return rs, nil
}

func proveFormula(f *formula, debugOn bool) ([]*Sequent, error) {
	i := 1
	solution := []*Sequent{}
	unreduced := []*Sequent{}
	f.Index = "0"
	unreduced = append(unreduced, &Sequent{Right: []*formula{f}, Name: "S1"})
	for len(unreduced) > 0 {
		if debugOn {
			fmt.Println("Unreduced:")
			for _, u := range unreduced {
				fmt.Printf("\t%s\n", u)
			}
			fmt.Println("Partial Solution:")
			for _, s := range solution {
				fmt.Printf("\t%s\n", s)
			}
		}

		last := unreduced[len(unreduced)-1]
		unreduced = unreduced[:len(unreduced)-1]
		// Try to apply each rule
		for _, rule := range rules {
			s, err := rule.applyRuleTo(&unreduced)
			if err != nil {
				return solution, err
			}
			if s != nil {
				if debugOn {
					fmt.Printf("Rule %s was applied\n", rule.getName())
				}
				// The rule was applied successfully
				i = i + 1
				s.Name = fmt.Sprintf("S%d", i)
				s.Justification = []string{rule.getName(), last.Name}

				if len(s.Left) == 0 && len(s.Right) == 0 {
					// A solution was found
					solution = append(solution, s)
					return solution, nil
				}
				unreduced = append(unreduced, s)
			}
			// else the rule was not appliable
		}

		solution = append(solution, last)
	}
	return solution, nil
}

// Prove givent a set of formulas it output a solution, if debugOn is true debugging messages will be printed
func Prove(rf *RawFormula, debugOn bool) ([]*Sequent, error) {
	if debugOn {
		fmt.Printf("Input:\n\t%s\n", rf.Formula)
	}
	tokens, err := tokenize(strings.Replace(rf.Formula, " ", "", -1), 0x00)
	if debugOn {
		fmt.Println("Tokens:")
		for i := len(tokens) - 1; i >= 0; i-- {
			fmt.Printf("\t%d: %s\n", len(tokens)-i, tokens[i].Value)
		}
	}
	if err != nil {
		return nil, err
	}
	top, err := genFormulasTree(tokens)
	if debugOn {
		fmt.Println("Formula:")
		fmt.Printf("\t%s\n", top)
	}
	if err != nil {
		return nil, err
	}
	s, err := proveFormula(top, debugOn)
	if debugOn {
		fmt.Println("Sequents:")
		for _, Sequent := range s {
			fmt.Printf("\t%s\n", Sequent)
		}
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// EncodeSequentSlice returns a map of latex encoded sequnets
func EncodeSequentSlice(in []*Sequent) (*map[int]RawSequent, error) {
	rawSolution := make(map[int]RawSequent)
	for i, s := range in {
		rs, err := encodeSequent(s)
		if err != nil {
			return &rawSolution, err
		}
		rawSolution[i] = rs
	}
	return &rawSolution, nil
}
