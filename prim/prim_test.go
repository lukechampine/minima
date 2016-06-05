package prim

import "testing"

// TestReadEval tests the Read and Eval functions. (Testing either by itself
// would be unnecessarily inconvenient.)
func TestReadEval(t *testing.T) {
	// lisp primatives, sans label/lambda
	var atomT = newAtom("t")
	var atomNil = newAtom("nil")
	var lisp Primitives
	eval := func(s *Sexp) *Sexp { return Eval(s, lisp) }
	lisp = Primitives{
		// return argument unevaluated
		"quote": func(s *Sexp) *Sexp {
			return s
		},
		// return t if s is an atom, nil otherwise
		"atom": func(s *Sexp) *Sexp {
			if eval(s).IsAtom() {
				return atomT
			}
			return atomNil
		},
		// return t if s.X and s.Y are both atoms and are equal
		"eq": func(s *Sexp) *Sexp {
			ex, ey := eval(s.X), eval(s.Y)
			if ex.IsAtom() && ey.IsAtom() && *ex.Atom == *ey.Atom {
				return atomT
			}
			return atomNil
		},
		// return s.X
		"car": func(s *Sexp) *Sexp {
			return eval(s.X)
		},
		// return s.Y
		"cdr": func(s *Sexp) *Sexp {
			return eval(s.Y)
		},
		// return s.Y with s.X.Atom prepended
		"cons": func(s *Sexp) *Sexp {
			return &Sexp{X: eval(s.X), Y: eval(s.Y)}
		},
		// return first expr whose predicate evaluates to t
		"cond": func(s *Sexp) *Sexp {
			if e := eval(s.X.X); e.IsAtom() && *e.Atom == *atomT.Atom {
				return eval(s.X.Y)
			}
			return lisp["cond"](s.Y)
		},
	}

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
	}

	for _, test := range tests {
		expr, err := ReadString(test.expr)
		if err != nil {
			t.Error("failed to parse expr:", err)
		} else if res := Eval(expr, lisp); res.String() != test.res {
			t.Errorf("expected %v, got %v", test.res, res.String())
		}
	}
}
