package parser

import (
	"fmt"
	"testing"

	"../ast"
	"../lexer"
)

func TestAsStatements(t *testing.T) {
	input := `
	as x = 5;
	as y = 10;
	as z = 895678;
	`

	// Initialize a lexer, parser and a program.
	lex := lexer.New(input)
	par := New(lex)
	program := par.Parse()
	checkParseErrors(t, par)

	if program == nil {
		t.Fatalf("Parse returned nil.")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program doesn't contain 3 statements, got: %d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"z"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testAsStatements(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testAsStatements(t *testing.T, statement ast.Statement, expectedIdentifier string) bool {
	if statement.TokenLiteral() != "as" {
		t.Errorf("Statement's token literal is not 'as'. Got: %q",
			statement.TokenLiteral())
		return false
	}

	asStatement, ok := statement.(*ast.DeclareStatement)
	if !ok {
		t.Errorf("Statement not a DeclareStatement. Got: %q", statement)
		return false
	}

	if asStatement.Name.Value != expectedIdentifier {
		t.Errorf("Statement doesn't contain an expected identifier %q. Got: %q",
			expectedIdentifier, asStatement.Name.Value)
		return false
	}

	if asStatement.Name.TokenLiteral() != expectedIdentifier {
		t.Errorf("Statement doesn't contain an expected token literal %q. Got: %q",
			expectedIdentifier, asStatement.Name)
		return false
	}

	return true
}

func checkParseErrors(t *testing.T, par *Parser) {
	errors := par.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser had %d errors.", len(errors))
	for _, msg := range errors {
		t.Errorf("Parse error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
	ret 5;
	ret 10;
	ret 959854;
	`

	lex := lexer.New(input)
	par := New(lex)
	program := par.Parse()
	checkParseErrors(t, par)

	if len(program.Statements) != 3 {
		t.Fatalf("Program doesn't contain 3 statements, got: %d",
			len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Statement not a ReturnStatement. Got: %T", statement)
			continue
		}
		if returnStatement.TokenLiteral() != "ret" {
			t.Errorf("ReturnStatement's token literal not 'rt'. Got: %q",
				returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "add;"

	lex := lexer.New(input)
	par := New(lex)
	program := par.Parse()
	checkParseErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("Program should contain only one statement. Got: %d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not a ExpressionStatement. Got: %T",
			program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Statement does not contain an identifier. Got: %T",
			statement.Expression)
	}

	if ident.Value != "add" {
		t.Fatalf("Identifier does not contain a value 'add'. Got: %q",
			ident.Value)
	}

	if ident.TokenLiteral() != "add" {
		t.Fatalf("Identifier's TokenLiteral does not contain a value 'add'. Got: %q",
			ident.TokenLiteral())
	}
}

func TestIntegerIdentifierExpression(t *testing.T) {
	input := "10;"

	lex := lexer.New(input)
	par := New(lex)
	program := par.Parse()
	checkParseErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("Program should contain only one statement. Got: %d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not a ExpressionStatement. Got: %T",
			program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Statement does not contain an identifier. Got: %T",
			statement.Expression)
	}

	if ident.Value != 10 {
		t.Fatalf("Identifier does not contain a value '10'. Got: %q",
			ident.Value)
	}

	if ident.TokenLiteral() != "10" {
		t.Fatalf("Identifier's TokenLiteral does not contain a value '10'. Got: %q",
			ident.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.Parse()
		checkParseErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected only one ExpressionStatement. Got: %d",
				len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement to be an ExpressionStatement. Got: %T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expected expression to be a PrefixExpression. Got: %T",
				statement.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("Operator doesn't match. Expected: %q. Got: %q",
				tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(
	t *testing.T,
	il ast.Expression,
	value int64) bool {

	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expression not an IntegerLiteral. Got: %T",
			il)
		return false
	}

	if integer.Value != value {
		t.Errorf("Values don't match. Expected: %d. Got: %d",
			value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("TokenLiteral not %d. Got: %s",
			value, integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(
	t *testing.T,
	exp ast.Expression,
	value string) bool {

	identifier, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression not an Identifier. Got: %T",
			exp)
		return false
	}

	if identifier.Value != value {
		t.Errorf("Identifier's value %q doesn't match the expected %q.",
			identifier.Value, value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("Identifier's TokenLiteral %q doesn't match the expected %q.",
			identifier.TokenLiteral(), value)
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{}) bool {

	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("Could not handle an expression %T", exp)
	return false
}

func testBooleanLiteral(
	t *testing.T,
	exp ast.Expression,
	value bool) bool {

	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("Expression not a Boolean. Got: %T",
			exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("Boolean value %t not an expected value %t",
			bo.Value, value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("Boolean TokenLiteral %s not an expected %t.",
			bo.TokenLiteral(), value)
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expression not an InfixExpression. Got: %T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("Operator %q doesn't match the expected %q.",
			opExp.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"false == false", false, "==", false},
		{"true != false", true, "!=", false},
	}

	for _, tt := range infixTests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.Parse()
		checkParseErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected only %d expression statement. Got: %d",
				1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected an ExpressionStatement. Got: %T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Expected an InfixExpression. Got: %T",
				statement.Expression)
		}

		if !testLiteralExpression(t, expression.Left, tt.leftValue) {
			return
		}

		if expression.Operator != tt.operator {
			t.Fatalf("Expected a %q operator. Got: %q",
				tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a * b + c / d", "((a * b) + (c / d))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"(5 + 5) * 3", "((5 + 5) * 3)"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"2 / (5 * 5)", "(2 / (5 * 5))"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.Parse()
		checkParseErrors(t, par)

		if tt.expected != program.String() {
			t.Fatalf("Expected operator precedence: %q. Got: %q",
				tt.expected, program.String())
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	lex := lexer.New(input)
	par := New(lex)
	program := par.Parse()
	checkParseErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected one program statement. Got: %d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected an ExpressionStatement. Got: %T",
			program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected an IfExpression. Got: %T",
			statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not a one statement. Got: %d",
			len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence not an ExpressionStatement. Got: %T",
			expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("Expression's alternative was not nil. Got: %T",
			expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	lex := lexer.New(input)
	par := New(lex)
	program := par.Parse()
	checkParseErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected one program statement. Got: %d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected an ExpressionStatement. Got: %T",
			program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected an IfExpression. Got: %T",
			statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not a one statement. Got: %d",
			len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence not an ExpressionStatement. Got: %T",
			expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not a one statement. Got: %d",
			len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative not an ExpressionStatement. Got: %T",
			expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}
