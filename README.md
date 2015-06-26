minima
------

minima is a minimal Lisp. At its core is a very restrictive grammar:

```
An s-expression is classically defined inductively as

1. an atom, or
2. an expression of the form (x . y) where x and y are s-expressions.
```

Atop this core are layers of syntactic sugar that, taken together, comprise a Scheme-like syntax. For example, the `.` token is made optional, and lists written as `(x y z)` expand to `(x . (y . (z . nil)))`.

The interpreter binary only implements `Eval` and `Apply`. These are used to write a meta-circular evaluator that implements the traditional Lisp primitives: `quote`, `atom`, `eq`, `car`, `cdr`, `cons`, `cond`, `label`, and `lambda`.

Finally, a standard library of higher-order functions are implemented using these primitives (`map`, `fold`, `zip`, etc.)

Installation
============

Download and install using `go get`:

`go get github.com/lukechampine/minima`

Then simply run `minima` to start the REPL.
