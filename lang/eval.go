package lang

import "fmt"

var (
	atomT   Atom = "t"
	atomNil Atom = "nil"

	env = map[Atom]*Sexp{
		"t":   &Sexp{Atom: &atomT},
		"nil": &Sexp{Atom: &atomNil},
	}

	primitives = map[Atom]func(*Sexp) *Sexp{
		"quote": primitiveQuote,
	}
)

func primitiveQuote(s *Sexp) *Sexp { return s }

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
