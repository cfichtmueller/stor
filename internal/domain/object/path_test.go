package object

import "testing"

type splitPathTest struct {
	key   string
	parts []string
}

func TestSplitPath(t *testing.T) {
	tests := []splitPathTest{
		{key: "", parts: []string{}},
		{key: "file.txt", parts: []string{"file.txt"}},
		{key: "a/b.txt", parts: []string{"a", "b.txt"}},
		{key: "a/b/c/", parts: []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		expectedLength := len(test.parts)
		parts := SplitPath(test.key, "/")
		actualLength := len(parts)
		if actualLength != expectedLength {
			t.Errorf("Expected %d parts, got %d", expectedLength, actualLength)
		}
		for i, a := range parts {
			if expectedLength <= i {
				t.Errorf("Didn't expect part '%s' on position %d", a, i)
			} else if a != test.parts[i] {
				t.Errorf("Expected part %d to be '%s', got '%s'", i, test.parts[i], a)
			}
		}
	}
}

type joinPathTest struct {
	parts    []string
	expected string
}

func TestJoinPath(t *testing.T) {
	tests := []joinPathTest{
		{parts: []string{}, expected: ""},
		{parts: []string{"a"}, expected: "a/"},
		{parts: []string{"a", "b"}, expected: "a/b/"},
		{parts: []string{"a", "b", "c"}, expected: "a/b/c/"},
	}

	for _, test := range tests {
		actual := JoinPath(test.parts, "/")
		if actual != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, actual)
		}
	}
}

type pathPrefixTest struct {
	key      string
	expected string
}

func TestPathPrefix(t *testing.T) {
	tests := []pathPrefixTest{
		{key: "", expected: ""},
		{key: "a.pdf", expected: ""},
		{key: "a/b.pdf", expected: "a/"},
		{key: "a/b/c.pdf", expected: "a/b/"},
	}

	for _, test := range tests {
		actual := PathPrefix(test.key, "/")
		if actual != test.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", test.expected, test.key, actual)
		}
	}
}
