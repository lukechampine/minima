package prim

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

func expect(r *bufio.Reader, exp rune) error {
	c, _, err := r.ReadRune()
	if err != nil {
		return err
	} else if c != exp {
		return errors.New("expected " + string(exp) + ", got " + string(c))
	}
	return nil
}

func readAtom(r *bufio.Reader) (*Sexp, error) {
	// read until non-atom rune encountered
	atom := ""
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		} else if c == '(' || c == ')' || c == ' ' {
			r.UnreadRune()
			break
		}
		atom += string(c)
	}
	return newAtom(atom), nil
}

func readSexp(r *bufio.Reader) (*Sexp, error) {
	// read ( or atom
	c, _, err := r.ReadRune()
	if err != nil {
		return nil, err
	} else if c != '(' {
		r.UnreadRune()
		return readAtom(r)
	}

	// read two space-separated Sexps
	x, err := readSexp(r)
	if err != nil {
		return nil, err
	}
	if err = expect(r, ' '); err != nil {
		return nil, err
	}
	y, err := readSexp(r)
	if err != nil {
		return nil, err
	}

	// read )
	if err = expect(r, ')'); err != nil {
		return nil, err
	}

	return &Sexp{X: x, Y: y}, nil
}

// Read parses a Sexp from an io.Reader.
func Read(r io.Reader) (*Sexp, error) {
	return readSexp(bufio.NewReader(r))
}

// Read parses a Sexp from a string.
func ReadString(expr string) (*Sexp, error) {
	return Read(strings.NewReader(expr))
}
