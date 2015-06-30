package lang

import "fmt"

var atomT = Atom("t")
var atomNil = Atom("nil")

func sexpT() *Sexp {
	return &Sexp{Atom: &atomT}
}

func sexpNil() *Sexp {
	return &Sexp{Atom: &atomNil}
}

func bool2Sexp(p bool) *Sexp {
	if p {
		return sexpT()
	}
	return sexpNil()
}

func null(x *Sexp) bool {
	return x.IsAtom() && *x.Atom == atomNil
}

// atom returns t if x is an atom, and nil otherwise.
func atom(x *Sexp) *Sexp {
	return bool2Sexp(x.IsAtom())
}

// eq returns true if x and y are both nil, or if they are the same atom, and
// false otherwise.
func eq(x, y *Sexp) bool {
	return x.IsAtom() && y.IsAtom() && *x.Atom == *y.Atom
}

// car returns x.X.
func car(x *Sexp) *Sexp {
	if x.IsAtom() {
		panic("car: not a list: " + x.String())
	}
	return x.X
}

// cdr returns x.Y
func cdr(x *Sexp) *Sexp {
	if x.IsAtom() {
		panic("cdr: not a list: " + x.String())
	}
	return x.Y
}

// helpers
func caar(x *Sexp) *Sexp   { return car(car(x)) }
func cadr(x *Sexp) *Sexp   { return car(cdr(x)) }
func cadar(x *Sexp) *Sexp  { return car(cdr(car(x))) }
func caddr(x *Sexp) *Sexp  { return car(cdr(cdr(x))) }
func caddar(x *Sexp) *Sexp { return car(cdr(cdr(car(x)))) }

// cons returns an S-expression with x as its car and y as its cdr.
func cons(x, y *Sexp) *Sexp {
	return &Sexp{X: x, Y: y}
}

// list is a cons plus a nil at the end.
func list(x, y *Sexp) *Sexp {
	return cons(x, cons(y, sexpNil()))
}

// concat returns the concatenation of x and y.
func concat(x, y *Sexp) *Sexp {
	if null(x) {
		return y
	}
	return cons(car(x), concat(cdr(x), y))
}

// assoc returns the mapping of x in y.
func assoc(x, y *Sexp) *Sexp {
	if null(y) {
		panic("undefined atom: " + x.String())
	}
	if eq(x, caar(y)) {
		return cadar(y)
	}
	return assoc(x, cdr(y))
}

// cond evalutes a list of mappings from predicates to values, returning the
// first value with a non-null predicate.
func cond(c, a *Sexp) *Sexp {
	if !null(eval(caar(c), a)) {
		return eval(cadar(c), a)
	}
	return cond(cdr(c), a)
}

// zip joins two equal-length lists into one list of x,y pairs.
func zip(x, y *Sexp) *Sexp {
	if null(x) && null(y) {
		return sexpNil()
	}
	return cons(list(car(x), car(y)), zip(cdr(x), cdr(y)))
}

// evlis evaluates the contents of a list.
func evlis(m, a *Sexp) *Sexp {
	if null(m) {
		return sexpNil()
	}
	return cons(eval(car(m), a), evlis(cdr(m), a))
}

// eval is the actual eval function, which may panic.
func eval(e, a *Sexp) *Sexp {
	if e.IsAtom() {
		return assoc(e, a)
	}
	if car(e).IsAtom() {
		switch *car(e).Atom {
		case "quote":
			return cdr(e)
		case "atom":
			return atom(eval(cadr(e), a))
		case "eq":
			return bool2Sexp(eq(eval(cadr(e), a), eval(caddr(e), a)))
		case "car":
			return car(eval(cadr(e), a))
		case "cdr":
			return cdr(eval(cadr(e), a))
		case "cons":
			return cons(eval(cadr(e), a), eval(caddr(e), a))
		case "cond":
			return cond(cdr(e), a)
		default:
			return eval(cons(assoc(car(e), a), cdr(e)), a)
		}
	}
	switch *caar(e).Atom {
	case "label":
		return eval(cons(caddar(e), cdr(e)), cons(list(cadar(e), car(e)), a))
	case "lambda":
		return eval(caddar(e), concat(zip(cadar(e), evlis(cdr(e), a)), a))
	}
	panic("could not evaluate expression: " + e.String())
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
	s = eval(sexp, sexpNil())
	return
}
