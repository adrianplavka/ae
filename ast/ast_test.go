package ast

import (
	"testing"

	"github.com/adrianplavka/fe/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&DeclareStatement{
				Token: token.Token{Type: token.DECLARE, Literal: "as"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "y"},
					Value: "y",
				},
			},
		},
	}

	if program.String() != "as x = y;" {
		t.Errorf("Program didn't return a testing string. Got: %q", program.String())
	}
}
