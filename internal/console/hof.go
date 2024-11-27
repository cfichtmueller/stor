package console

import "github.com/cfichtmueller/jug"

func withPrincipal(f func(c jug.Context, principal string)) func(c jug.Context) {
	return func(c jug.Context) {
		principal := contextMustGetPrincipal(c)
		f(c, principal)
	}
}
