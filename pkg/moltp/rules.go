package moltp

import "fmt"

type (
	inferenceRule interface {
		getName() string
		applyRuleTo(s *Sequent) (*Sequent, error)
	}
	resolutionRule interface {
		getName() string
		applyRuleTo(sequents *[]*Sequent) ([]*Sequent, error)
	}
	r1 struct {
		Name string
		R    *relation
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
		R    *relation
	}
	r7 struct {
		Name         string
		worldsKeeper *worldskeeper
	}
	r8 struct {
		Name         string
		worldsKeeper *worldskeeper
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
func (r r1) applyRuleTo(sequents *[]*Sequent) ([]*Sequent, error) {
	for _, s1 := range *sequents {
		l1 := len(s1.Left)
		if l1 < 1 {
			continue
		}
		f1 := s1.Left[l1-1]
		if len(f1.Operands) == 0 { // This means it is an atomic formula
			for _, s2 := range *sequents {
				l2 := len(s2.Right)
				if l2 < 1 {
					continue
				}
				f2 := s2.Right[0]
				if len(f2.Operands) == 0 {
					g := r.R.munify(f1, f2)
					if g != nil {
						n := &Sequent{}

						t1 := g.applyUnifications(s1.Left[:l1-1])
						t2 := g.applyUnifications(s2.Left)
						n.Left = append(t1, t2...)

						t1 = g.applyUnifications(s1.Right)
						t2 = g.applyUnifications(s2.Right[1:])
						n.Right = append(t1, t2...)

						n.Justification = []string{r.Name, s1.Name, s2.Name, fmt.Sprintf("%s", g)}

						return []*Sequent{s1, s2, n}, nil
					}
				}
			}
		}
	}
	return []*Sequent{}, nil
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
	if f.Terminal == sIMPLY {
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
	if f.Terminal == sIMPLY {
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
	if f.Terminal == sIMPLY {
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
		ns := &worldsymbol{Ground: true}
		if f.Index.isGround() {
			ns.Value = fmt.Sprintf("%d", r.worldsKeeper.NextIndex)
		} else {
			// TODO: Implement corret world index value
			ns.Value = fmt.Sprintf("f(%s)", r.worldsKeeper.NextVar)
			r.worldsKeeper.updateNextVariable()
			// r.worldsKeeper.NextVar
		}
		r.worldsKeeper.NextIndex = r.worldsKeeper.NextIndex + 1
		t.Index.Symbols = append([]*worldsymbol{ns}, f.Index.Symbols...)
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
		if f.Index.isGround() {
			t.Index.Symbols = append([]*worldsymbol{&worldsymbol{Value: r.worldsKeeper.NextVar, Index: 0, Ground: false}}, f.Index.Symbols...)
			r.worldsKeeper.updateNextVariable()
		} else {
			t.Index.Symbols = append([]*worldsymbol{&worldsymbol{Value: start(&f.Index).Value, Index: start(&f.Index).Index + 1, Ground: false}}, f.Index.Symbols...)
		}
		// n.Left = append(s.Left[:l-1], t) TODO: WTF!!!!
		b := []*formula{}
		for _, p := range s.Left[:l-1] {
			b = append(b, p)
		}
		n.Left = append(b, t)
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
