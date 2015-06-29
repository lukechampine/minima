package lang

import "fmt"

var (
	// variable bindings
	env map[Atom]*Sexp
	// primitive functions
	primitives map[Atom]func(*Sexp) *Sexp
)

func init() {
	atomT := Atom("t")
	atomNil := Atom("nil")

	env = map[Atom]*Sexp{
		"t":   &Sexp{Atom: &atomT},
		"nil": &Sexp{Atom: &atomNil},
	}

	primitives = map[Atom]func(*Sexp) *Sexp{
		"quote": primitiveQuote,
		"atom":  primitiveAtom,
		"eq":    primitiveEq,
		"car":   primitiveCar,
		"cdr":   primitiveCdr,
		"cons":  primitiveCons,
		"cond":  primitiveCond,
		"label": primitiveLabel,
	}
}

// quote returns s without evaluating it.
func primitiveQuote(s *Sexp) *Sexp {
	return s
}

// atom returns t if s is an atom, and nil otherwise.
func primitiveAtom(s *Sexp) *Sexp {
	if eval(s).IsAtom() {
		return env["t"]
	}
	return env["nil"]
}

// eq returns t if s.X == s.Y, and nil otherwise.
func primitiveEq(s *Sexp) *Sexp {
	if eval(s.X) == eval(s.Y) {
		return env["t"]
	}
	return env["nil"]
}

// car returns s.X, which must not be nil.
func primitiveCar(s *Sexp) *Sexp {
	s = eval(s)
	if s.X == nil {
		panic("not a list: " + s.String())
	}
	return s.X
}

// cdr returns s.Y, which must not be nil.
func primitiveCdr(s *Sexp) *Sexp {
	s = eval(s)
	if s.Y == nil {
		panic("not a list: " + s.String())
	}
	return s.Y
}

// cons returns a list formed by prepending s.Y with s.X. Note that this is
// not equivalent to returning s; the contents of s are evaluated.
func primitiveCons(s *Sexp) *Sexp {
	// TODO: is it safe to evaluate directly, as below?
	//s.X, s.Y = eval(s.X), eval(s.Y)
	return &Sexp{X: eval(s.X), Y: eval(s.Y)}
}

// cond evaluates a set of predicates and their value mappings. Predicates are
// evaluated sequentially, and the first predicate that does not evalute to
// nil is the "winner." The value mapping of this predicate is returned. If
// all predicates evaluate to nil, cond returns nil.
func primitiveCond(s *Sexp) *Sexp {
	if s.X.IsAtom() {
		panic("clause is not a list: " + s.X.String())
	}
	if eval(s.X.X) != env["nil"] {
		return eval(s.X.Y)
	}
	// recurse to next clause
	return primitiveCond(s.Y)
}

// label adds a new atom->Sexp mapping to the environment. The evaluated Sexp
// is returned.
func primitiveLabel(s *Sexp) *Sexp {
	if !s.X.IsAtom() {
		panic("cannot use non-atom as label: " + s.X.String())
	}
	env[*s.X.Atom] = s.Y
	return eval(s.Y)
}

// apply applies a function to a list of arguments. Only primitives can be
// applied. Attempting to apply anything else triggers a panic.
func apply(car, cdr *Sexp) *Sexp {
	if !car.IsAtom() {
		panic("not a function: " + car.String())
	}
	prim, ok := primitives[*car.Atom]
	if !ok {
		panic("not a primitive: " + car.String())
	}
	return prim(cdr)
}

// eval is the actual eval function, which may panic.
func eval(sexp *Sexp) *Sexp {
	if sexp.IsAtom() {
		val, ok := env[*sexp.Atom]
		if !ok {
			panic("undefined atom: " + sexp.String())
		}
		return val
	}
	if sexp.X == nil || sexp.Y == nil {
		panic("invalid S-expression")
	}
	if sexp.X.IsAtom() {
		return apply(sexp.X, sexp.Y)
	}
	return apply(eval(sexp.X), sexp.Y)
}

// Eval evaluates an S-expression. The procedure for evaluation is:
//
//  - If the sexp is an atom, look up its mapping in the environment and
//  return it.
//  - If the sexp is a pair, Eval the car and then apply it to the cdr.
func Eval(sexp *Sexp) (s *Sexp, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	s = eval(sexp)
	return
}
