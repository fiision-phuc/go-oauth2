package util

import "testing"

func Test_EncryptPassword(t *testing.T) {
	hash, _ := EncryptPassword("12345678")
	if len(hash) == 0 {
		t.Error("Expected password is not nil but found nil.")
	}

	if !ComparePassword(hash, "12345678") {
		t.Error("Expected true but found false.")
	}
}
