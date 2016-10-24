package inject

import "testing"

type TestStruct struct {
	ValueString string  `field:"string" validation:"^\\w+(\\s\\w+)?$"`
	ValueBool   bool    `field:"bool"`
	ValueFloat  float64 `field:"float" validation:"^\\d+(\\.\\d+)?$"`
	ValueInt    int64   `field:"int" validation:"^-?\\d+$"`
	ValueInt8   int8    `field:"int8"`
	ValueUInt   uint64  `field:"uint"`
}

func Test_BindForm(t *testing.T) {
	values := map[string]string{
		"string":      "A String",
		"bool":        "FaLSE",
		"float":       "1234.94502",
		"int":         "-100",
		"int8":        "1",
		"uint":        "200",
		"invalid_int": "abcdef",
	}
	testStruct := new(TestStruct)
	err := BindForm(values, testStruct)

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
}
