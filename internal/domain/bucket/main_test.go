// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bucket

import "testing"

type bucketNameTest struct {
	name  string
	valid bool
}

func TestBucketNames(t *testing.T) {
	tests := []bucketNameTest{
		{name: "", valid: false},
		{name: "a", valid: false},
		{name: "aa", valid: false},
		{name: "aaa", valid: true},
		{name: "1aa", valid: true},
		{name: "1a1", valid: true},
		{name: "111", valid: true},
		{name: "a.aa", valid: true},
		{name: "a.a.a", valid: true},
		{name: "a.a.a.", valid: false},
		{name: "a..a.a", valid: false},
		{name: "a-a-a", valid: true},
		{name: "a-a-a-", valid: false},
		{name: "a/a/a", valid: false},
		{name: ".aaa", valid: false},
		{name: "..aaa", valid: false},
		{name: "-aaa", valid: false},
		{name: "--aaa", valid: false},
		{name: "AAAAA", valid: false},
		{name: "aaaAaaa", valid: false},
	}

	for _, test := range tests {
		valid := bucketNamePattern.MatchString(test.name)
		if valid != test.valid {
			t.Errorf("Name '%s' is %v, expected %v", test.name, valid, test.valid)
		}
	}
}
