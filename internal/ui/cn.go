// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "strings"

const (
	btn                 = "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 h-8 rounded-md px-3 text-xs"
	btnPrimary          = "bg-primary text-primary-foreground shadow hover:bg-primary/80"
	btnSecondary        = "bg-secondary text-secondary-foreground shadow-sm hover:bg-secondary/80"
	btnDanger           = "bg-danger text-danger-foreground shadow-sm hover:bg-danger/80"
	cnInput             = "flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
	cnLabel             = "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
	cnDIalogCloseButton = "absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground"
)

func cn(cn ...string) string {
	index := make(map[string]interface{})
	c := make([]string, 0)

	for _, cur := range cn {
		for _, e := range strings.Split(cur, " ") {
			if _, ok := index[e]; ok {
				continue
			}
			index[e] = nil
			c = append(c, e)
		}
	}

	return strings.Join(c, " ")
}
