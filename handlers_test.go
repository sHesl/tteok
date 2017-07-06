package tteok

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_single_string(t *testing.T) {
	input := "hi"
	result := testBuildLog(input)

	if result.Message != input {
		t.Errorf("Expected log.Message to be %s. Got %s.", input, result.Message)
	}
}

func Test_two_strings(t *testing.T) {
	input1 := "hi"
	input2 := "hello"
	inputs := []string{input1, input2}

	result := testBuildLog(input1, input2)

	if !reflect.DeepEqual(result.Messages, inputs) {
		t.Errorf("Expected log.Messages to be %v. Got %v.", inputs, result.Messages)
	}
}

func Test_three_strings(t *testing.T) {
	input1 := "hi"
	input2 := "hello"
	input3 := "hey"
	inputs := []string{input1, input2, input3}

	result := testBuildLog(input1, input2, input3)

	if !reflect.DeepEqual(result.Messages, inputs) {
		t.Errorf("Expected log.Messages to be %v. Got %v.", inputs, result.Messages)
	}
}

func Test_error(t *testing.T) {
	input := fmt.Errorf("this is an error")
	result := testBuildLog(input)

	if result.Error != input.Error() {
		t.Errorf("Expected log.Error to be %s. Got %s.", input.Error(), result.Error)
	}

	if !strings.Contains(result.Stack, "Test_error(") {
		t.Errorf("Expected log.Stack to be a valid stack dump. Got %s.", result.Stack)
	}
}

func Test_pointer(t *testing.T) {
	input := "hello"
	result := testBuildLog(&input)

	if result.Message != input {
		t.Errorf("Expected log.Message to be %s. Got %s.", input, result.Message)
	}
}

func Test_pointer_to_pointer(t *testing.T) {
	input := "hello"
	inputPtr := &input
	result := testBuildLog(&inputPtr)

	if result.Message != input {
		t.Errorf("Expected log.Message to be %s. Got %s.", input, result.Message)
	}
}

func Test_slice(t *testing.T) {
	input1 := "hello"
	input2 := "hey"
	inputs := []string{input1, input2}
	result := testBuildLog(inputs)

	if !reflect.DeepEqual(result.Messages, inputs) {
		t.Errorf("Expected log.Messages to be %v. Got %v.", inputs, result.Messages)
	}
}

func Test_map(t *testing.T) {
	input := make(map[string]interface{})
	input["hi"] = 123
	input["bye"] = 321

	result := testBuildLog(input)

	if !reflect.DeepEqual(result.Data, input) {
		t.Errorf("Expected log.Data to be %v. Got %v.", input, result.Data)
	}
}

func Test_struct(t *testing.T) {
	type nestedStruct struct {
		Nested bool
	}

	type testStruct struct {
		Hi       string
		ValuePtr *int
		Nested   nestedStruct
	}

	val := 123
	testObj := testStruct{"hello message", &val, nestedStruct{true}}

	result := testBuildLog(testObj)

	if result.Data["Hi"] == nil || result.Data["Hi"] != testObj.Hi {
		t.Errorf("Expected log.Data.Hi to be %v. Got %v.", testObj.Hi, result.Data["Hi"])
	}

	if result.Data["ValuePtr"] == nil || result.Data["ValuePtr"] != val {
		t.Errorf("Expected log.Data.ValuePtr to be %v. Got %v.", val, result.Data["ValuePtr"])
	}

	if result.Data["Nested"] == nil {
		t.Errorf("Expected log.Data.Nested to be a struct. Got %v.", result.Data["Nested"])
	}
}

func testBuildLog(params ...interface{}) log {
	return buildLog("TEST", params)
}
