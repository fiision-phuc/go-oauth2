package utils

import (
	"net/url"
	"testing"
)

type TestStruct struct {
	ValueString string  `inject:"string"`
	ValueBool   bool    `inject:"bool"`
	ValueFloat  float64 `inject:"float"`
	ValueInt    int64   `inject:"int"`
	ValueUInt   uint64  `inject:"uint"`

	ValueInvalid int `inject:"invalid_int"`
}

func Test_BindForm(t *testing.T) {
	err := BindForm(nil, nil)
	if err == nil {
		t.Error("Expected error when sending nil.")
	}

	err = BindForm(nil, "A string")
	if err == nil {
		t.Error("Expected error return when sending non struct type.")
	}

	values := url.Values{
		"string":      []string{"A String"},
		"bool":        []string{"FaLSE"},
		"float":       []string{"1234.94502"},
		"int":         []string{"-100"},
		"uint":        []string{"200"},
		"invalid_int": []string{"abcdef"},
	}
	testStruct := TestStruct{}
	err = BindForm(values, &testStruct)

	if err != nil {
		t.Error("Expected error nil when sending valid input.")
	}

	if testStruct.ValueString != "A String" {
		t.Errorf("Expected %s but found %s.", "A String", testStruct.ValueString)
	}

	if testStruct.ValueBool {
		t.Errorf("Expected %t but found %t.", false, testStruct.ValueBool)
	}

	if testStruct.ValueFloat != 1234.94502 {
		t.Errorf("Expected %f but found %f.", 1234.94502, testStruct.ValueFloat)
	}

	if testStruct.ValueInt != -100 {
		t.Errorf("Expected %d but found %d.", -100, testStruct.ValueInt)
	}

	if testStruct.ValueUInt != 200 {
		t.Errorf("Expected %d but found %d.", 200, testStruct.ValueUInt)
	}

	if testStruct.ValueInvalid != 0 {
		t.Errorf("Expected %d but found %d.", 0, testStruct.ValueInvalid)
	}
}