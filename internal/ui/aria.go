// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func ariaCurrent(current string) e.Node {
	return e.Attr("aria-current", current)
}

func ariaDisabled() e.Node {
	return e.Attr("aria-disabled", "true")
}

func ariaLabel(label string) e.Node {
	return e.Attr("aria-label", label)
}

func ariaRole(role string) e.Node {
	return e.Attr("role", role)
}

// returns a class="sr-only" node
func srOnly() e.Node {
	return e.Class("sr-only")
}
