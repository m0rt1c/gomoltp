package moltp

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	sBOX     = "Box"
	sDIAMOND = "Diamond"
	sEXISTS  = "Exists"
	sFORALL  = "Forall"
	sIFF     = "Iff"
	sIMPLY   = "Implies"
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

	// TODO: report? copying the array directly is not what I intended
	// dst.Operands = src.Operands
	// even if they are two arrays and not pointers to an array they are treated as they were pointers

	dst.Operands = append([]*formula{}, src.Operands...)
	dst.Terminal = src.Terminal
	dst.Index = src.Index
	dst.Vars = append([]string{}, src.Vars...)

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
			case sIMPLY:
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
			case sIMPLY:
				return true
			case sIFF:
				return true
			}
		case sIMPLY:
			switch b.Value {
			case sOR:
				return false
			case sAND:
				return false
			case sIMPLY:
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
	sEncoding[sIMPLY] = "\\to"
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
		return &token{IsOp: true, BiOp: true, Value: sIMPLY, Skip: len("\\to")}
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
		// TODO: supoort multi letters varaibles
		if len(s) > 4 && s[1] == '(' {
			skip := 2
			vlist := []string{}
			vars := ""
			closed := false
			for i := 2; i < len(s); i++ {
				skip = skip + 1
				if s[i] == ')' {
					closed = true
					break
				}
				if s[i] == ',' {
					continue
				}
				vlist = append(vlist, fmt.Sprintf("%c", s[i]))
				if vars == "" {
					vars = fmt.Sprintf("%c", s[i])
				} else {
					vars = fmt.Sprintf("%s,%c", vars, s[i])
				}
			}
			if !closed {
				return nil, fmt.Errorf("Missing closing parenthesis for %s", s)
			}
			v := fmt.Sprintf("%c", s[0])
			return &token{IsTe: true, Value: v, Skip: skip, Vars: vlist}, nil
		}
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

func reduceFormulas(f *formula) *formula {
	for i, g := range f.Operands {
		f.Operands[i] = reduceFormulas(g)
	}
	switch f.Terminal {
	case sDIAMOND:
		// \Diamond A = \lnot \Box \lnot A
		A := f.Operands[0]
		g0 := &formula{Terminal: sNOT, Operands: []*formula{A}}
		g1 := &formula{Terminal: sBOX, Operands: []*formula{g0}}
		return &formula{Terminal: sNOT, Operands: []*formula{g1}}
	case sIFF:
		// A <-> B = ( A \to B ) \and ( B \to A ) = \lnot ( (A \to B) \to \lnot ( B \to A) )
		A := f.Operands[0]
		B := f.Operands[1]
		g0 := &formula{Terminal: sIMPLY, Operands: []*formula{B, A}}
		g1 := &formula{Terminal: sNOT, Operands: []*formula{g0}}
		g2 := &formula{Terminal: sIMPLY, Operands: []*formula{A, B}}
		g3 := &formula{Terminal: sIMPLY, Operands: []*formula{g2, g1}}
		return &formula{Terminal: sNOT, Operands: []*formula{g3}}
	case sAND:
		// A \land B = \lnot ( A \to \lnot B )
		A := f.Operands[0]
		B := f.Operands[1]
		g0 := &formula{Terminal: sNOT, Operands: []*formula{B}}
		g1 := &formula{Terminal: sIMPLY, Operands: []*formula{A, g0}}
		return &formula{Terminal: sNOT, Operands: []*formula{g1}}
	case sOR:
		// A \lor B = \lnot A \to B
		A := f.Operands[0]
		B := f.Operands[1]
		g0 := &formula{Terminal: sNOT, Operands: []*formula{A}}
		return &formula{Terminal: sIMPLY, Operands: []*formula{g0, B}}
	case sEXISTS:
		// \exists x p = \lnot \forall x \lnot p
		g0 := &formula{Terminal: sNOT, Operands: []*formula{f.Operands[len(f.Operands)-1]}}
		g1 := &formula{Terminal: sFORALL, Operands: append(f.Operands[:len(f.Operands)-1], g0), Vars: f.Vars}
		return &formula{Terminal: sNOT, Operands: []*formula{g1}}
	default:
		return f
	}
}

func genFormulasTree(tokens []*token) (*formula, error) {
	var formulas []*formula
	for _, t := range tokens {
		if t.IsOp {
			if t.MuOp {
				// We must have (1) a formula and a (2) list of variables name
				// Something like forall x \Box x -> x
				if len(formulas) < 2 {
					return formulas[0], fmt.Errorf("missing arguments for multi operator %s", t.Value)
				}
				f := &formula{}
				f.Terminal = t.Value
				// (1) This should find the formula
				m := formulas[len(formulas)-1]
				formulas = formulas[:len(formulas)-1]
				f.Operands = append(f.Operands, formulas[len(formulas)-1])
				f.Vars = append(f.Vars, formulas[len(formulas)-1].Terminal)
				formulas = formulas[:len(formulas)-1]
				// (2) this should find all the variables, mind that they are in the reversed order
				for k := len(formulas) - 1; k >= 0; k-- {
					if formulas[k].Terminal == "," {
						if k-1 < 0 {
							return formulas[0], fmt.Errorf("missing argument for multi operator %s", t.Value)
						}
						f.Operands = append([]*formula{formulas[k-1]}, f.Operands...)
						f.Vars = append(f.Vars, formulas[k-1].Terminal)
					} else {
						formulas = formulas[:k+1]
						break
					}
				}

				f.Operands = append(f.Operands, m)
				formulas = append(formulas, f)
			}
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
				f.Operands = append(f.Operands, formulas[len(formulas)-1])
				formulas = formulas[:len(formulas)-1]
				formulas = append(formulas, f)
			}
		}
		if t.IsTe {
			formulas = append(formulas, &formula{Terminal: t.Value, Vars: t.Vars})
		}
		if t.IsIn {
			if len(formulas) < 1 {
				return formulas[0], fmt.Errorf("trying to assign index %s to nothing", t.Value)
			}
			formulas[len(formulas)-1].Index = worldindex{[]*worldsymbol{&worldsymbol{Ground: true, Value: t.Value}}}
		}
	}
	return reduceFormulas(formulas[0]), nil
}

func encodeSequent(s *Sequent) (RawSequent, error) {
	rs := RawSequent{}

	rs.Name = s.Name
	rs.Left = toTextRepr(formulaArrayToString(s.Left))
	rs.Right = toTextRepr(formulaArrayToString(s.Right))
	rs.Justification = ""
	for _, j := range s.Justification {
		if rs.Justification == "" {
			rs.Justification = fmt.Sprintf("%s", j)
		} else {
			rs.Justification = fmt.Sprintf("%s, %s", rs.Justification, j)
		}
	}

	return rs, nil
}

func (p *Prover) proveFormula(f *formula) ([]*Sequent, error) {
	i := 1
	solution := []*Sequent{}
	unreduced := []*Sequent{}
	reduced := []*Sequent{}

	f.Index = worldindex{[]*worldsymbol{p.worldsKeeper.GetFreeIndividualConstant()}}

	unreduced = append(unreduced, &Sequent{Right: []*formula{f}, Name: "S1"})

	for len(unreduced) > 0 {
		if p.Debug {
			log.Println("******************************")
			log.Println("**** Applying rules loop *****")
			log.Println("******************************")
			log.Println("Unreduced:")
			for _, u := range unreduced {
				log.Printf("\t%s\n", u)
			}
			if len(solution) < 1 {
				log.Println("Solution is empty")
			} else {
				log.Println("Partial Solution:")
				for _, s := range solution {
					log.Printf("\t%s\n", s)
				}
			}
			if len(reduced) < 1 {
				log.Println("Reduced list is empty")
			} else {
				log.Println("Reduced:")
				for _, s := range reduced {
					log.Printf("\t%s\n", s)
				}
			}
		}

		pushLastInSolution := false
		last := unreduced[len(unreduced)-1]
		new := []*Sequent{}

		// Try to apply each rule
		for _, rule := range p.Rules {
			s, err := rule.applyRuleTo(last)
			if err != nil {
				return solution, err
			}
			if s != nil {
				if p.Debug {
					log.Printf("Rule %s was applied on %s\n", rule.getName(), last)
				}
				pushLastInSolution = true
				// The rule was applied successfully
				i = i + 1
				s.Name = fmt.Sprintf("S%d", i)
				s.Justification = []string{rule.getName(), last.Name}

				if len(s.Left) == 0 && len(s.Right) == 0 {
					// A solution was found
					solution = append(solution, s)
					return solution, nil
				}
				new = append(new, s)
				if p.Debug {
					log.Printf("New sequent is %s\n", s)
				}
			}
			// else the rule was not appliable
		}

		if pushLastInSolution {
			solution = append(solution, last)
		} else {
			// If no rule was appliable to the last element
			// we move it at the beginning of the reduced rules
			reduced = append(reduced, last)
		}
		unreduced = append(unreduced[:len(unreduced)-1], new...)

	}

	if p.Debug {
		log.Println("******************************")
		log.Println("***** Unreduced are over *****")
		log.Println("******************************")
		log.Println("Unreduced:")
		for _, u := range unreduced {
			log.Printf("\t%s\n", u)
		}
		if len(solution) < 1 {
			log.Println("Solution is empty")
		} else {
			log.Println("Partial Solution:")
			for _, s := range solution {
				log.Printf("\t%s\n", s)
			}
		}
		if len(reduced) < 1 {
			log.Println("Reduced list is empty")
		} else {
			log.Println("Reduced:")
			for _, s := range reduced {
				log.Printf("\t%s\n", s)
			}
		}
	}

	if len(reduced) > 1 {
		rule := p.ResolutionRule
		res, err := rule.applyRuleTo(&reduced)
		if err != nil {
			return solution, err
		}
		if len(res) > 0 {
			if p.Debug {
				log.Printf("Rule %s was applied on %s\n", rule.getName(), reduced)
			}
			s := res[2]
			// The rule was applied successfully
			i = i + 1
			s.Name = fmt.Sprintf("S%d", i)

			if p.Debug {
				log.Printf("New sequent is %s\n", s)
			}

			if len(s.Left) == 0 && len(s.Right) == 0 {
				// A solution was found
				solution = append(solution, res...)
				return solution, nil
			}
		}
	}

	if p.Debug {
		log.Println("******************************")
		log.Printf("******* %s was applied *******\n", p.ResolutionRule.getName())
		log.Println("******************************")
		log.Println("Unreduced:")
		for _, u := range unreduced {
			log.Printf("\t%s\n", u)
		}
		if len(solution) < 1 {
			log.Println("Solution is empty")
		} else {
			log.Println("Partial Solution:")
			for _, s := range solution {
				log.Printf("\t%s\n", s)
			}
		}
		if len(reduced) < 1 {
			log.Println("Reduced list is empty")
		} else {
			log.Println("Reduced:")
			for _, s := range reduced {
				log.Printf("\t%s\n", s)
			}
		}
	}

	return solution, fmt.Errorf("No solution found")
}

// Prove givent a set of formulas it output a solution, if debugOn is true debugging messages will be printed
func (p *Prover) Prove(rf *RawFormula) ([]*Sequent, error) {
	p.initProver()
	if p.Debug {
		log.Println("Input:")
		log.Printf("\t%s\n", rf.Formula)
	}
	tokens, err := tokenize(strings.Replace(rf.Formula, " ", "", -1), 0x00)
	if p.Debug {
		log.Println("Tokens:")
		for i := len(tokens) - 1; i >= 0; i-- {
			t := tokens[i]
			if len(t.Vars) == 0 {
				log.Printf("\t%d: %s\n", len(tokens)-i, t.Value)
			} else {
				log.Printf("\t%d: %s Vars: %s\n", len(tokens)-i, t.Value, t.Vars)
			}
		}
	}
	if err != nil {
		return nil, err
	}
	top, err := genFormulasTree(tokens)
	if p.Debug {
		log.Println("Formula:")
		log.Printf("\t%s\n", top)
	}
	if err != nil {
		return nil, err
	}
	s, err := p.proveFormula(top)
	if p.Debug {
		log.Println("Sequents:")
		for _, Sequent := range s {
			log.Printf("\t%s\n", Sequent)
		}
	}
	if err != nil {
		return s, err
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
