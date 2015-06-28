package lang

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tError tokenType = iota
	tAtom
	tLParen
	tRParen
	tDot
)

type token struct {
	typ tokenType
	val string
}

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
	}
	if s.IsAtom() {
		return string(*s.Atom)
	} else if s.X == nil && s.Y == nil {
		return "invalid S-expression"
	}
	return fmt.Sprintf("(%s.%s)", s.X, s.Y)
}

// newAtom creates an atom Sexp from a string.
func newAtom(a string) *Sexp {
	at := Atom(a)
	return &Sexp{Atom: &at}
}

// parseSexp parses a token stream into a Sexp, either as an atom or a pair.
func parseSexp(tokens <-chan token) (*Sexp, error) {
	switch t := <-tokens; t.typ {
	case tError:
		return nil, errors.New(t.val)
	case tAtom:
		return newAtom(t.val), nil
	case tLParen:
		return parsePair(tokens)
	case tRParen:
		return nil, errors.New("unexpected )")
	case tDot:
		return nil, errors.New("unexpected .")
	}
	panic("unknown token")
}

// parsePair parses a Sexp pair.
func parsePair(tokens <-chan token) (*Sexp, error) {
	x, err := parseSexp(tokens)
	if err != nil {
		return nil, err
	}
	if t := <-tokens; t.typ != tDot {
		return nil, errors.New("unexpected " + t.val)
	}
	y, err := parseSexp(tokens)
	if err != nil {
		return nil, err
	}
	if t := <-tokens; t.typ != tRParen {
		return nil, errors.New("unexpected " + t.val)
	}
	return &Sexp{X: x, Y: y}, nil
}

// scanSexp is a bufio.SplitFunc, used to tokenize input.
func scanSexp(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// need at least one byte
	if len(data) == 0 {
		return 0, nil, nil
	}

	// dot or paren
	r, _ := utf8.DecodeRune(data)
	switch r {
	case '(', '.', ')':
		return 1, data[:1], nil
	}
	// otherwise, must be a letter
	if !unicode.IsLetter(r) {
		return 0, nil, fmt.Errorf("illegal rune '%c'", r)
	}

	// atom
	for width, i := 0, 0; i < len(data); i += width {
		r, width = utf8.DecodeRune(data[i:])
		if !unicode.IsLetter(r) {
			return i, data[:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}

	// request more data
	return 0, nil, nil
}

// tokenize reads tokens from r and sends them down a channel.
func tokenize(r io.Reader, tokens chan token) {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanSexp)
	for scanner.Scan() {
		switch t := scanner.Text(); t {
		case "(":
			tokens <- token{tLParen, t}
		case ".":
			tokens <- token{tDot, t}
		case ")":
			tokens <- token{tRParen, t}
		default:
			tokens <- token{tAtom, t}
		}
	}
	if err := scanner.Err(); err != nil {
		tokens <- token{tError, err.Error()}
	}
}

// Read parses a Sexp from an io.Reader.
func Read(r io.Reader) (*Sexp, error) {
	tokens := make(chan token)
	go tokenize(r, tokens)
	return parseSexp(tokens)
}

// Read parses a Sexp from a string.
func ReadString(exp string) (*Sexp, error) {
	return Read(strings.NewReader(exp))
}
