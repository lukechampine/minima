package lang

import (
	"testing"

	"github.com/lukechampine/minima/prim"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		res  string
	}{
		{"(quote foo)", "foo"},
		{"(quote (atom foo))", "(atom foo)"},
		{"(atom foo)", "t"}, // no error about unknown symbols; need an environment for that
		{"(atom (foo bar))", "nil"},
		{"(atom (quote bar))", "t"},
		{"(eq (foo foo))", "t"},
		{"(eq (foo (quote foo)))", "t"},
		{"(car (foo bar))", "foo"},
		{"(cdr (foo bar))", "bar"},
		{"(cons (foo (bar baz)))", "(foo (bar baz))"},
		{"(cons (foo (cons (bar baz))))", "(foo (bar baz))"},
		{"(cond ((nil foo) ((t bar) nil)))", "bar"}, // nil at end simplifies cond implementation
		{"(lambda (x (cons (x (quote (y z))))))", "__func-1"},
		{"((lambda ((x nil) (cons ((lookup x) (quote (y z)))))) ((quote x) nil))", "(x (y z))"},
		{"(label (f ((lambda (nil (quote x))) f)))", "f"},
		{"((label (f ((lambda (nil (quote x))) f))) (f nil))", "x"},
	}

	for _, test := range tests {
		expr, err := prim.ReadString(test.expr)
		if err != nil {
			t.Error("failed to parse expr:", err)
		} else if res := Eval(expr); res.String() != test.res {
			t.Errorf("expected %v, got %v", test.res, res.String())
		}
	}
}
