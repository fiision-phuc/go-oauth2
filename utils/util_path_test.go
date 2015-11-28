package utils

import "testing"

func TestCleanPath(t *testing.T) {
	data := []struct {
		path, result string
	}{
		// Already clean
		{"/", "/"},
		{"/abc", "/abc"},
		{"/a/b/c", "/a/b/c"},
		{"/abc/", "/abc/"},
		{"/a/b/c/", "/a/b/c/"},

		// Missing root
		{"", "/"},
		{"abc", "/abc"},
		{"abc/def", "/abc/def"},
		{"a/b/c", "/a/b/c"},

		// Remove doubled slash
		{"//", "/"},
		{"/abc//", "/abc/"},
		{"/abc/def//", "/abc/def/"},
		{"/a/b/c//", "/a/b/c/"},
		{"/abc//def//ghi", "/abc/def/ghi"},
		{"//abc", "/abc"},
		{"///abc", "/abc"},
		{"//abc//", "/abc/"},

		// Remove . elements
		{".", "/"},
		{"./", "/"},
		{"/abc/./def", "/abc/def"},
		{"/./abc/def", "/abc/def"},
		{"/abc/.", "/abc/"},

		// Remove .. elements
		{"..", "/"},
		{"../", "/"},
		{"../../", "/"},
		{"../..", "/"},
		{"../../abc", "/abc"},
		{"/abc/def/ghi/../jkl", "/abc/def/jkl"},
		{"/abc/def/../ghi/../jkl", "/abc/jkl"},
		{"/abc/def/..", "/abc"},
		{"/abc/def/../..", "/"},
		{"/abc/def/../../..", "/"},
		{"/abc/def/../../..", "/"},
		{"/abc/def/../../../ghi/jkl/../../../mno", "/mno"},

		// Combinations
		{"abc/./../def", "/def"},
		{"abc//./../def", "/def"},
		{"abc/../../././../def", "/def"},
	}

	for _, test := range data {
		if s := FormatPath(test.path); s != test.result {
			t.Errorf("FormatPath(%q) = %q, want %q", test.path, s, test.result)
		}
		if s := FormatPath(test.result); s != test.result {
			t.Errorf("FormatPath(%q) = %q, want %q", test.result, s, test.result)
		}
	}
}
