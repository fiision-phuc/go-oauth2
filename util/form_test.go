package util

import (
	"testing"

	"github.com/phuc0302/go-oauth2/test"
)

type TestStruct struct {
	ValueString string  `field:"string" validation:"^\\w+(\\s\\w+)?$"`
	ValueBool   bool    `field:"bool"`
	ValueFloat  float64 `field:"float" validation:"^\\d+(\\.\\d+)?$"`
	ValueInt    int64   `field:"int" validation:"^-?\\d+$"`
	ValueInt8   int8    `field:"int8"`
	ValueUInt   uint64  `field:"uint"`
}

func Test_BindForm_WithEmptyData(t *testing.T) {
	values := make(map[string]string)
	testStruct := new(TestStruct)

	err := BindForm(values, testStruct)
	if err == nil {
		t.Error(test.ExpectedNotNil)
	}
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

func Test_BindForm_InternalForm(t *testing.T) {
	values := map[string]string{
		"string":      "A String",
		"bool":        "FaLSE",
		"float":       "1234.94502",
		"int":         "-100",
		"int8":        "1",
		"uint":        "200",
		"invalid_int": "abcdef",
	}

	var form struct {
		ValueString string  `field:"string" validation:"^\\w+(\\s\\w+)?$"`
		ValueBool   bool    `field:"bool"`
		ValueFloat  float64 `field:"float" validation:"^\\d+(\\.\\d+)?$"`
		ValueInt    int64   `field:"int" validation:"^-?\\d+$"`
		ValueInt8   int8    `field:"int8"`
		ValueUInt   uint64  `field:"uint"`
	}
	err := BindForm(values, &form)

	if err != nil {
		t.Error("Expected error nil when sending valid input.")
	}

	if form.ValueString != "A String" {
		t.Errorf("Expected %s but found %s.", "A String", form.ValueString)
	}

	if form.ValueBool {
		t.Errorf("Expected %t but found %t.", false, form.ValueBool)
	}

	if form.ValueFloat != 1234.94502 {
		t.Errorf("Expected %f but found %f.", 1234.94502, form.ValueFloat)
	}

	if form.ValueInt != -100 {
		t.Errorf("Expected %d but found %d.", -100, form.ValueInt)
	}
	if form.ValueInt8 != 1 {
		t.Errorf("Expected %d but found %d.", 1, form.ValueInt8)
	}
	if form.ValueUInt != 200 {
		t.Errorf("Expected %d but found %d.", 200, form.ValueUInt)
	}
}
