package lang

import (
	"strconv"

	"github.com/lukechampine/minima/prim"
)

func sexpToSlice(s *prim.Sexp) []*prim.Sexp {
	if s.IsAtom() {
		return []*prim.Sexp{s}
	}
	return append([]*prim.Sexp{s.X}, sexpToSlice(s.Y)...)
}

// Eval evaluates an Sexp using minimal Lisp primitives.
func Eval(expr *prim.Sexp) *prim.Sexp {
	// constants
	var atomT = prim.Atom("t")
	var atomNil = prim.Atom("nil")
	var sexpT = &prim.Sexp{Atom: &atomT}
	var sexpNil = &prim.Sexp{Atom: &atomNil}

	// environment (for lookup)
	env := map[prim.Atom]*prim.Sexp{
		atomT:   sexpT,
		atomNil: sexpNil,
	}

	// labels (for lambdas)
	var labelCount int

	// primitives
	lisp := prim.Primitives{}

	// inner eval function, closing around lisp
	eval := func(s *prim.Sexp) *prim.Sexp {
		return prim.Eval(s, lisp)
	}

	// helper for creating primitives from lambda exprs
	newLambda := func(s *prim.Sexp) func(*prim.Sexp) *prim.Sexp {
		return func(args *prim.Sexp) *prim.Sexp {
			// bind param names to eval'd args
			paramSlice := sexpToSlice(s.X)
			argSlice := sexpToSlice(args)
			for i := range paramSlice {
				env[*paramSlice[i].Atom] = eval(argSlice[i])
			}
			// eval lambda body
			return eval(s.Y)
		}
	}

	// return argument unevaluated
	lisp["quote"] = func(s *prim.Sexp) *prim.Sexp {
		return s
	}

	// return t if s is an atom, nil otherwise
	lisp["atom"] = func(s *prim.Sexp) *prim.Sexp {
		if eval(s).IsAtom() {
			return sexpT
		}
		return sexpNil
	}

	// return t if s.X and s.Y are both atoms and are equal
	lisp["eq"] = func(s *prim.Sexp) *prim.Sexp {
		ex, ey := eval(s.X), eval(s.Y)
		if ex.IsAtom() && ey.IsAtom() && *ex.Atom == *ey.Atom {
			return sexpT
		}
		return sexpNil
	}

	// return s.X
	lisp["car"] = func(s *prim.Sexp) *prim.Sexp {
		return eval(s.X)
	}

	// return s.Y
	lisp["cdr"] = func(s *prim.Sexp) *prim.Sexp {
		return eval(s.Y)
	}

	// return s.Y with s.X.Atom prepended
	lisp["cons"] = func(s *prim.Sexp) *prim.Sexp {
		return &prim.Sexp{X: eval(s.X), Y: eval(s.Y)}
	}

	// return associated value in env
	lisp["lookup"] = func(s *prim.Sexp) *prim.Sexp {
		if s.IsAtom() {
			if v, ok := env[*s.Atom]; ok {
				return v
			}
		}
		return s
	}

	// return first expr whose predicate evaluates to t
	lisp["cond"] = func(s *prim.Sexp) *prim.Sexp {
		if e := eval(s.X.X); e.IsAtom() && *e.Atom == atomT {
			return eval(s.X.Y)
		}
		return lisp["cond"](s.Y)
	}

	// return an atom bound to the new lambda function
	lisp["lambda"] = func(s *prim.Sexp) *prim.Sexp {
		labelCount++
		label := "__func-" + prim.Atom(strconv.Itoa(labelCount))
		lisp[label] = newLambda(s)
		return &prim.Sexp{Atom: &label}
	}

	// return the supplied atom, bound to the supplied lambda
	lisp["label"] = func(s *prim.Sexp) *prim.Sexp {
		lisp[*s.X.Atom] = newLambda(s.Y.X.Y)
		return s.X
	}

	return eval(expr)
}
