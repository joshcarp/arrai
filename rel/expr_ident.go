package rel

import (
	"fmt"
)

// IdentLookupFailed represents failure to lookup a scope variable.
type IdentLookupFailed struct {
	scope Scope
	expr  IdentExpr
}

func (e IdentLookupFailed) Error() string {
	return fmt.Sprintf("Name %q not found in %v", e.expr.ident, e.scope.m.Keys())
}

// IdentExpr returns the variable referenced by ident.
type IdentExpr struct {
	ident string
}

// DotIdent represents the special identifier '.'.
var DotIdent = IdentExpr{"."}

// NewIdentExpr returns a new identifier.
func NewIdentExpr(ident string) IdentExpr {
	if ident == "." {
		return DotIdent
	}
	return IdentExpr{ident}
}

// Ident returns the ident for the IdentExpr.
func (e IdentExpr) Ident() string {
	return e.ident
}

// String returns a string representation of the expression.
func (e IdentExpr) String() string {
	if e.ident == "." {
		return "(" + e.ident + ")"
	}
	return e.ident
}

// Eval returns the value from scope corresponding to the ident.
func (e IdentExpr) Eval(local Scope) (Value, error) {
	if a, found := local.Get(e.ident); found && a != nil {
		return a.Eval(local)
	}
	return nil, IdentLookupFailed{local, e}
}
