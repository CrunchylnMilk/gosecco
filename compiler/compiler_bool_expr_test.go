package compiler

import (
	"testing"

	"github.com/twtiger/gosecco/asm"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func BoolTest(t *testing.T) { TestingT(t) }

type BoolCompilerSuite struct{}

var _ = Suite(&BoolCompilerSuite{})

func (s *BoolCompilerSuite) Test_orExpressionBetweenEqualityComparisons(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Or{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	defaultPositiveReturn := "ret_k\t7FFF0000\n"
	defaultNegativeReturn := "ret_k\t0\n"

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+

		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	04	00	2A\n"+

		"ld_abs	18\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	00	01	2A\n"+

		defaultPositiveReturn+
		defaultNegativeReturn)
}

func (s *BoolCompilerSuite) Test_compilationOfAndExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	07	0\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	2A\n"+
		"ld_abs	18\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *BoolCompilerSuite) Test_negatedAndExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Negation{
					Operand: tree.And{
						Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
						Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	06	0\n"+
		"ld_abs	14\n"+
		"jeq_k	00	04	2A\n"+
		"ld_abs	18\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *BoolCompilerSuite) Test_compilationOfNegatedOrExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Negation{
					Operand: tree.Or{
						Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
						Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	05	00	2A\n"+
		"ld_abs	18\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *BoolCompilerSuite) Test_compilationOfNestedNegatedAndExpressionRightSide(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Negation{
						Operand: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	07	0\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	2A\n"+
		"ld_abs	18\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *BoolCompilerSuite) Test_compilationOfNestedNegatedAndExpressionLeftSide(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left: tree.Negation{
						Operand: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
					Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	05	00	2A\n"+
		"ld_abs	18\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *BoolCompilerSuite) Test_compilationOfNestedNegatedOrExpressionRightSide(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Or{
					Left: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Negation{
						Operand: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	04	00	2A\n"+
		"ld_abs	18\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *BoolCompilerSuite) Test_compilationOfNestedNegatedOrExpressionLeftSide(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Or{
					Left: tree.Negation{
						Operand: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
					Right: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	18\n"+
		"jeq_k	00	06	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	00	04	2A\n"+
		"ld_abs	10\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	14\n"+
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *BoolCompilerSuite) Test_compilationOfNegatedEqualsComparison(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Negation{
					Operand: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t05\t1\n"+
		"ld_abs\t10\n"+
		"jeq_k\t00\t02\t0\n"+
		"ld_abs\t14\n"+
		"jeq_k\t01\t00\t2A\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}

func (s *BoolCompilerSuite) Test_compilingBooleanInsideExpressionShouldPanicSinceItsAProgrammerError(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.BooleanLiteral{false},
				},
			},
		},
	}
	c.Assert(func() {
		Compile(p)
	}, Panics, "Programming error: there should never be any boolean literals left outside of the toplevel if the simplifier works correctly: syscall: write - (and (eq arg0 42) false)")
}
