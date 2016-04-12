package parser

import (
	"github.com/twtiger/go-seccomp/tree"
	. "gopkg.in/check.v1"
)

type RulesSuite struct{}

var _ = Suite(&RulesSuite{})

func (s *RulesSuite) Test_parsesSimpleRule(c *C) {
	result, _ := parseExpression("1")

	c.Assert(result, DeepEquals, tree.BooleanLiteral{true})
}

func (s *RulesSuite) Test_parsesAlmostSimpleRule(c *C) {
	result, _ := parseExpression("arg0 > 0")

	c.Assert(result, DeepEquals, tree.Comparison{
		Left:  tree.Argument{0},
		Op:    tree.GT,
		Right: tree.NumericLiteral{0},
	})
}

func (s *RulesSuite) Test_parseAnotherRule(c *C) {
	result, _ := parseExpression("arg0 == 4")

	c.Assert(result, DeepEquals, tree.Comparison{
		Left:  tree.Argument{0},
		Op:    tree.EQL,
		Right: tree.NumericLiteral{4},
	})
}

func (s *RulesSuite) Test_parseYetAnotherRule(c *C) {
	result, _ := parseExpression("arg0 == 4 || arg0 == 5")

	c.Assert(tree.ExpressionString(result), Equals, "(or (eq arg0 4) (eq arg0 5))")
	c.Assert(result, DeepEquals, tree.Or{
		Left: tree.Comparison{
			Left:  tree.Argument{0},
			Op:    tree.EQL,
			Right: tree.NumericLiteral{4},
		},
		Right: tree.Comparison{
			Left:  tree.Argument{0},
			Op:    tree.EQL,
			Right: tree.NumericLiteral{5},
		},
	})
}

func (s *RulesSuite) Test_parseExpressionWithMultiplication(c *C) {
	result, _ := parseExpression("arg0 == 12 * 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (mul 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithAddition(c *C) {
	result, _ := parseExpression("arg0 == 12 + 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (plus 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithDivision(c *C) {
	result, _ := parseExpression("arg0 == 12 / 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (div 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithSubtraction(c *C) {
	result, _ := parseExpression("arg0 == 12 - 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (minus 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryAnd(c *C) {
	result, _ := parseExpression("arg0 == 0 & 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (binand 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryOr(c *C) {
	result, _ := parseExpression("arg0 == 0 | 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (binor 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryXor(c *C) {
	result, _ := parseExpression("arg0 == 0 ^ 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (binxor 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryNegation(c *C) {
	c.Skip("not yet implemented, check binary negation syntax")
	result, _ := parseExpression("arg0 == ^0")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (bnot 0))")
}

func (s *RulesSuite) Test_parseAExpressionLeftShift(c *C) {
	result, _ := parseExpression("arg0 == 2 << 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (lsh 2 1))")
}

func (s *RulesSuite) Test_parseAExpressionRightShift(c *C) {
	result, _ := parseExpression("arg0 == 2 >> 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (rsh 2 1))")
}

func (s *RulesSuite) Test_parseAExpressionWithModulo(c *C) {
	result, _ := parseExpression("arg0 == 12 % 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (mod 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithBooleanAnd(c *C) {
	result, _ := parseExpression("arg0 == 0 && arg1 == 0")
	c.Assert(tree.ExpressionString(result), Equals, "(and (eq arg0 0) (eq arg1 0))")
}

func (s *RulesSuite) Test_parseAExpressionWithBooleanNegation(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("!arg0")
	c.Assert(tree.ExpressionString(result), Equals, "(not arg0)")
}

func (s *RulesSuite) Test_parseAExpressionWithNotEqual(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 != arg1")
	c.Assert(tree.ExpressionString(result), Equals, "(neq arg0 arg1")
}

func (s *RulesSuite) Test_parseAExpressionWithGreaterThanOrEqualTo(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 >= arg1")
	c.Assert(tree.ExpressionString(result), Equals, "(geq arg0 arg1")
}

func (s *RulesSuite) Test_parseAExpressionWithLessThan(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 < arg1")
	c.Assert(tree.ExpressionString(result), Equals, "(lss arg0 arg1")
}

func (s *RulesSuite) Test_parseAExpressionWithLessThanOrEqualTo(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 <= arg1")
	c.Assert(tree.ExpressionString(result), Equals, "(leq arg0")
}

func (s *RulesSuite) Test_parseAExpressionWithBitSets(c *C) {
	c.Skip("not yet implemented, check syntax against how we use binary and")
	result, _ := parseExpression("arg0 & val")
	c.Assert(tree.ExpressionString(result), Equals, "(set arg0")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusion(c *C) {
	c.Skip("not yet implemented, check syntax about set")
	result, _ := parseExpression("in(arg0, 1, 2)")
	c.Assert(tree.ExpressionString(result), Equals, "(in arg0 {1, 2}")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusionLargerSet(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("in(arg0, 1, 2, 3, 4)")
	c.Assert(tree.ExpressionString(result), Equals, "(in arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusionWithWhitespace(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("in(arg0, 1,   2,   3,   4)")
	c.Assert(tree.ExpressionString(result), Equals, "(in arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithNotInclusion(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("notIn(arg0, 1, 2)")
	c.Assert(tree.ExpressionString(result), Equals, "(notIn arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithNotInclusionLargerSet(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("notin(arg0, 1, 2, 3, 4)")
	c.Assert(tree.ExpressionString(result), Equals, "(notin arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithTrue(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("true")
	c.Assert(tree.ExpressionString(result), Equals, "1")
}

func (s *RulesSuite) Test_parseAExpressionWithFalse(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("false")
	c.Assert(tree.ExpressionString(result), Equals, "0")
}

func (s *RulesSuite) Test_parseAExpressionWith0AsFalse(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("0")
	c.Assert(tree.ExpressionString(result), Equals, "0")
}

func (s *RulesSuite) Test_parseAExpressionWithNestedOperatorsWithParens(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 == (12 + 3) * 2")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (* (+ 12 3) 2)))")
}

func (s *RulesSuite) Test_parseAExpressionWithNestedOperators(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 == 12 + 3 * 2")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (+ 12 (* 3 2)))")
}

func (s *RulesSuite) Test_parseAExpressionWithInvalidArithmeticOperator(c *C) {
	c.Skip("not yet implemented, error handling")
	result, _ := parseExpression("arg0 == 12 _ 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (add 12 3))")
}

//	result, _ := doParse("read2: arg0 > 0")

//	c.Assert(tree.ExpressionString(result), DeepEquals, "(read2 (gt (arg 0) (literal 0)))")

// func (s *RulesSuite) Test_parsesSlightlyMoreComplicatedRule(c *C) {
// 	result, _ := doParse("write: arg1 == 42 || arg0 + 1 == 15 && (arg3 == 1 || arg4 == 2)")

// 	c.Assert(result, DeepEquals, []rule{
// 		rule{
// 			syscall: "write",
// 			expression: orExpr{
// 				Left: equalsComparison{
// 					Left: argumentNode{index: 1},
// 					Right: literalNode{value: 42},
// 				},
// 				Right: andExpr{
// 					Left: equalsComparison{
// 						Left: addition{
// 							Left: argumentNode{index: 0},
// 							Right: literalNode{value: 1},
// 						},
// 						Right: literalNode{value: 15},
// 					},
// 					Right: orExpr{
// 						Left: equalsComparison{
// 							Left: argumentNode{index: 3},
// 							Right: literalNode{value: 1},
// 						},
// 						Right: equalsComparison{
// 							Left: argumentNode{index: 4},
// 							Right: literalNode{value: 2},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// }
