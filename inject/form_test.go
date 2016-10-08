package util

import "testing"

type TestStruct struct {
	ValueString string  `string`
	ValueBool   bool    `bool`
	ValueFloat  float64 `float`
	ValueInt    int64   `int`
	ValueInt8   int8    `int8`
	ValueUInt   uint64  `uint`

	ValueInvalid int `invalid_int`
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

	values := map[string]string{
		"string":      "A String",
		"bool":        "FaLSE",
		"float":       "1234.94502",
		"int":         "-100",
		"int8":        "1",
		"uint":        "200",
		"invalid_int": "abcdef",
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
	if testStruct.ValueInt8 != 1 {
		t.Errorf("Expected %d but found %d.", 1, testStruct.ValueInt8)
	}
	if testStruct.ValueUInt != 200 {
		t.Errorf("Expected %d but found %d.", 200, testStruct.ValueUInt)
	}

	if testStruct.ValueInvalid != 0 {
		t.Errorf("Expected %d but found %d.", 0, testStruct.ValueInvalid)
	}
}
