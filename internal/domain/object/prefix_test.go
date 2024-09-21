// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package object

import (
	"testing"
)

type PrefixTest struct {
	Delimiter        string
	Prefix           string
	Keys             []string
	ExpectedPrefixes []string
	ExpectedKeys     []string
}

func TestPrefixIndex(t *testing.T) {
	keys := []string{
		"documents/2020/contract.pdf",
		"documents/2021/contract.pdf",
		"documents/2021/info.pdf",
		"photos/2020/avatar.jpg",
		"photos/2021/avatar.jpg",
		"photos/avatar.jpg",
		"photos/group.jpg",
		"image.jpg",
	}

	tests := []PrefixTest{
		{Delimiter: "/", Prefix: "", Keys: keys, ExpectedPrefixes: []string{"documents/", "photos/"}, ExpectedKeys: []string{"image.jpg"}},
		{Delimiter: "/", Prefix: "documents/", Keys: keys, ExpectedPrefixes: []string{"documents/2020/", "documents/2021/"}, ExpectedKeys: []string{}},
		{Delimiter: "/", Prefix: "photos/", Keys: keys, ExpectedPrefixes: []string{"photos/2020/", "photos/2021/"}, ExpectedKeys: []string{"photos/avatar.jpg", "photos/group.jpg"}},
	}

	for testid, test := range tests {
		x := NewPrefixIndex(test.Delimiter, test.Prefix)
		for _, k := range test.Keys {
			x.AddKey(k)
		}
		actualPrefixes := x.CommonPrefixes
		if len(actualPrefixes) != len(test.ExpectedPrefixes) {
			t.Errorf("Expected %d prefixes, got %d: %v", len(test.ExpectedPrefixes), len(actualPrefixes), actualPrefixes)
			continue
		}
		for i, a := range actualPrefixes {
			if a != test.ExpectedPrefixes[i] {
				t.Errorf("Expected prefix '%s' in position %d, got: '%s'", test.ExpectedPrefixes[i], i, a)
			}
		}
		actualKeys := x.keys
		if len(actualKeys) != len(test.ExpectedKeys) {
			t.Errorf("[Test %d]: Expected %d keys, got %d, %v", testid, len(test.ExpectedKeys), len(actualKeys), actualKeys)
		}
		for i, a := range actualKeys {
			if a != test.ExpectedKeys[i] {
				t.Errorf("Expected key '%s' in position %d, got: '%s'", test.ExpectedKeys[i], i, a)
			}
		}
	}
}
