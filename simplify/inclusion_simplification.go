package simplify

import "github.com/twtiger/go-seccomp/tree"

// AcceptInclusion implements Visitor
func (s *simplifier) AcceptInclusion(a tree.Inclusion) {
	l := Simplify(a.Left)
	pl, pok := potentialExtractValue(l)

	result := make([]tree.Numeric, len(a.Rights))
	resultVals := make([]uint32, len(a.Rights))
	resultOks := make([]bool, len(a.Rights))
	for ix, v := range a.Rights {
		result[ix] = Simplify(v)
		resultVals[ix], resultOks[ix] = potentialExtractValue(result[ix])
	}

	if pok {
		newResults := []tree.Numeric{}
		for ix, v := range result {
			if resultOks[ix] {
				if resultVals[ix] == pl {
					s.result = tree.BooleanLiteral{a.Positive}
					return
				}
			} else {
				newResults = append(newResults, v)
			}
		}
		if len(newResults) == 0 {
			s.result = tree.BooleanLiteral{!a.Positive}
		} else {
			s.result = tree.Inclusion{Positive: a.Positive, Left: l, Rights: newResults}
		}
	} else {
		s.result = tree.Inclusion{Positive: a.Positive, Left: l, Rights: result}
	}
}