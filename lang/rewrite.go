package lang

import "regexp"

var (
	regexpLeading  = regexp.MustCompile(`^\s+`)
	regexpTrailing = regexp.MustCompile(`\s+$`)
	regexpLParen   = regexp.MustCompile(`\(\s+`)
	regexpRParen   = regexp.MustCompile(`\s+\)`)
	regexpNil      = regexp.MustCompile(`\(\s*\)`)
	regexpDot      = regexp.MustCompile(`\s+`)
)

// Desugar applies the following substitutions to its input, in order:
//
//  1. Remove leading and trailing whitespace.
//  2. Remove whitespace around parens
//  3. Replace () with nil.
//  4. Replace any remaining whitespace with .
//
// This allows for nicer-looking code. Eventually this function will be
// written in minima itself.
func Desugar(exp string) string {
	exp = regexpLeading.ReplaceAllString(exp, "")
	exp = regexpTrailing.ReplaceAllString(exp, "")
	exp = regexpLParen.ReplaceAllString(exp, "(")
	exp = regexpRParen.ReplaceAllString(exp, ")")
	exp = regexpNil.ReplaceAllString(exp, "nil")
	exp = regexpDot.ReplaceAllString(exp, ".")
	return exp
}
