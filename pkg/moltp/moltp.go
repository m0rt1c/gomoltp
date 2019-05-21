package moltp

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
	// symbolNOT     = "\\lnot"
	// symbolOR      = "\\lor"
	// symbolAND     = "\\land"
	symbolIMPLIES = "\\to"
	symbolALL     = "\\forall"
	symbolEXISTS  = "\\exists"
	symbolBOX     = "\\Box"
	symbolDIAMOND = "\\Diamond"
)

var (
	terminals = []string{}
)

func parseRawSequent(rs RawSequent) formula {
	return formula{}
}

// Prove givent a set of formulas it output a solution
func Prove(rf RawFormula) ([]RawSequent, error) {
	return []RawSequent{}, nil
}
