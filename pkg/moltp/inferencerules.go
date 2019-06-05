package moltp

import "fmt"

var (
	rules = []inferenceRule{r1{Name: "R1"}, r2{Name: "R2"}, r3{Name: "R3"}, r4{Name: "R4"}, r5{Name: "R5"}, r6{Name: "R6"}, r7{Name: "R7"}, r8{Name: "R8"}, r9{Name: "R9"}, r10{Name: "R10"}}
)

type (
	inferenceRule interface {
		getName() string
		applyRuleTo(s *Sequent) (*Sequent, error)
	}

	r1 struct {
		Name string
	}
	r2 struct {
		Name string
	}
	r3 struct {
		Name string
	}
	r4 struct {
		Name string
	}
	r5 struct {
		Name string
	}
	r6 struct {
		Name string
	}
	r7 struct {
		Name string
	}
	r8 struct {
		Name string
	}
	r9 struct {
		Name string
	}
	r10 struct {
		Name string
	}
)

// this functions rapresenting inference rules returns
// 1) a Sequent and a nil if the rule was applied successfully. The returned Sequent is the result of applying the rule
// 2) nil and nil if the Sequent was not appliable
// 3) nil and an error if there was some sort of error

// R1: If S,|p|_{i} <- T and S' <- |q|_{j}, T' and |p|_{i} and |q|_{j}
// unify with unification O then S_{O} U S'_{O} <- T_{O} U T'_{O}
func (r r1) applyRuleTo(s *Sequent) (*Sequent, error) {
	return nil, nil
}
func (r r1) getName() string {
	return r.Name
}

// R2: If S,|(p->q)|_{i} <- T then S,|q|_{i}<-|p|_{i},T
func (r r2) applyRuleTo(s *Sequent) (*Sequent, error) {
	l := len(s.Left)
	if l < 1 {
		return nil, nil
	}
	f := s.Left[l-1]
	if f.Terminal == sIMPLIES {
		n := &Sequent{}

		t := copyTopFormulaLevel(f.Operands[1])
		t.Index = f.Index
		n.Left = append(s.Left[:l-1], t)

		t = copyTopFormulaLevel(f.Operands[0])
		t.Index = f.Index
		n.Right = append([]*formula{t}, s.Right...)

		return n, nil
	}
	return nil, nil
}
func (r r2) getName() string {
	return r.Name
}

// R3: If S <- |(p->q)|_{i},T then S <- |q|_{i},T
func (r r3) applyRuleTo(s *Sequent) (*Sequent, error) {
	l := len(s.Right)
	if l < 1 {
		return nil, nil
	}
	f := s.Right[0]
	if f.Terminal == sIMPLIES {
		n := &Sequent{}

		t := copyTopFormulaLevel(f.Operands[1])
		t.Index = f.Index
		n.Right = append([]*formula{t}, s.Right[1:]...)
		n.Left = s.Left

		return n, nil
	}
	return nil, nil
}
func (r r3) getName() string {
	return r.Name
}

// R4: If S <- |(p->q)|_{i},T then S,|p|_{i} <- T
func (r r4) applyRuleTo(s *Sequent) (*Sequent, error) {
	l := len(s.Right)
	if l < 1 {
		return nil, nil
	}
	f := s.Right[0]
	if f.Terminal == sIMPLIES {
		n := &Sequent{}

		t := copyTopFormulaLevel(f.Operands[0])
		t.Index = f.Index
		n.Left = append(s.Left, t)
		n.Right = s.Right[1:]

		return n, nil
	}
	return nil, nil
}
func (r r4) getName() string {
	return r.Name
}

// R5: If S,| not p|_{i} <- T then S <- |p|_{i},T
func (r r5) applyRuleTo(s *Sequent) (*Sequent, error) {
	l := len(s.Left)
	if l < 1 {
		return nil, nil
	}
	f := s.Left[l-1]
	if f.Terminal == sNOT {
		n := &Sequent{}

		t := copyTopFormulaLevel(f.Operands[0])
		t.Index = f.Index
		n.Right = append([]*formula{t}, s.Right...)
		n.Left = s.Left[:l-1]

		return n, nil
	}
	return nil, nil
}
func (r r5) getName() string {
	return r.Name
}

// R6: If S <- |not p|_{i},T then S,|p|_{i} <- T
func (r r6) applyRuleTo(s *Sequent) (*Sequent, error) {
	l := len(s.Right)
	if l < 1 {
		return nil, nil
	}
	f := s.Right[0]
	if f.Terminal == sNOT {
		n := &Sequent{}

		t := copyTopFormulaLevel(f.Operands[0])
		t.Index = f.Index
		n.Left = append(s.Left, t)
		n.Right = s.Right[1:]

		return n, nil
	}
	return nil, nil
}
func (r r6) getName() string {
	return r.Name
}

// R7: If S <- | Box p|_{i},T then S <- |p|_{n:i},T
func (r r7) applyRuleTo(s *Sequent) (*Sequent, error) {
	l := len(s.Right)
	if l < 1 {
		return nil, nil
	}
	f := s.Right[0]
	if f.Terminal == sBOX {
		n := &Sequent{}

		t := copyTopFormulaLevel(f.Operands[0])
		// TODO: Implement corret world index value
		t.Index = fmt.Sprintf("w:%s", f.Index)
		n.Left = s.Left
		n.Right = append([]*formula{t}, s.Right[1:]...)

		return n, nil
	}
	return nil, nil
}
func (r r7) getName() string {
	return r.Name
}

// R8: If S,|Box p|_{i} <- T then S,|p|_{w:i} <- T
func (r r8) applyRuleTo(s *Sequent) (*Sequent, error) {
	l := len(s.Left)
	if l < 1 {
		return nil, nil
	}
	f := s.Left[l-1]
	if f.Terminal == sBOX {
		n := &Sequent{}

		t := copyTopFormulaLevel(f.Operands[0])
		// TODO: Implement corret world index value
		t.Index = f.Index
		n.Left = append(s.Left[:l-1], t)
		n.Right = s.Right

		return n, nil
	}
	return nil, nil
}
func (r r8) getName() string {
	return r.Name
}

func (r r9) applyRuleTo(s *Sequent) (*Sequent, error) {
	return nil, nil
}
func (r r9) getName() string {
	return r.Name
}

func (r r10) applyRuleTo(s *Sequent) (*Sequent, error) {
	return nil, nil
}
func (r r10) getName() string {
	return r.Name
}
