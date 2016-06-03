package utils

import "testing"

func Test_GenerateToken(t *testing.T) {
	hash := GenerateToken()
	if len(hash) != 40 {
		t.Errorf("Expected a token with the length of 40 chars but found %s", hash)
	}
}

func Benchmark_GenerateToken(b *testing.B) {
	GenerateToken()
}
