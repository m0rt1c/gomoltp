package moltp

type (
	// Formula object holding a single formula
	Formula struct {
		OID   int    `json:"oid"`
		Left  string `json:"left"`
		Right string `json:"right"`
	}
)

// Solve givent a set of formulas it output a solution
func Solve([]Formula) ([]Formula, error) {
	return []Formula{}, nil
}
