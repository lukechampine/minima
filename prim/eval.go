package prim

// An Atom is a string containing only runes that satisfy unicode.IsLetter,
// e.g. 'a', 'Z', 'λ'
type Atom string

// A Sexp is either an Atom or a pair of Sexps. The other field(s) will be
// nil.
type Sexp struct {
	Atom *Atom
	X, Y *Sexp
}

// IsAtom returns true if the Sexp is an Atom.
func (s *Sexp) IsAtom() bool {
	return s.Atom != nil
}

// String formats a Sexp for printing.
func (s *Sexp) String() string {
	if s == nil {
		return "empty S-expression"
	} else if s.IsAtom() {
		return string(*s.Atom)
	} else if s.X == nil || s.Y == nil {
		return "invalid S-expression"
	}
	return "(" + s.X.String() + " " + s.Y.String() + ")"
}

// newAtom creates an atom Sexp from a string.
func newAtom(a string) *Sexp {
	at := Atom(a)
	return &Sexp{Atom: &at}
}

// Primitives is a map from Atoms to Sexp evaluation functions.
type Primitives map[Atom]func(*Sexp) *Sexp

// Eval evaluates expr using primitives p.
func Eval(expr *Sexp, p Primitives) *Sexp {
	if p == nil || expr.IsAtom() {
		// nothing to do
		return expr
	}
	fn := Eval(expr.X, p)
	if !fn.IsAtom() || p[*fn.Atom] == nil {
		// non-atom or no primitive
		return &Sexp{X: fn, Y: expr.Y}
	}
	return p[*fn.Atom](expr.Y)
}
