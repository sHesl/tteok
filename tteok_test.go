package tteok

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_process(t *testing.T) {
	buf := &bytes.Buffer{}
	out = buf

	Info()

	result, err := extractLogOutput(buf)

	if err != nil {
		t.Errorf(err.Error())
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.Contains(result.Process, "tteok.test.exe") {
		t.Errorf("Expected log.Process to contain \"tteok.test.exe\". Got %s.", result.Process)
	}
}

func Test_timestamp(t *testing.T) {
	startTime := time.Now().Format(time.RFC3339)

	buf := &bytes.Buffer{}
	out = buf

	Info()

	result, err := extractLogOutput(buf)

	if err != nil {
		t.Errorf(err.Error())
	}

	if result.Timestamp > startTime || result.Timestamp > time.Now().Format(time.RFC3339) {
		t.Errorf("log.timestamp does not have a correct value.")
	}
}

func Test_at(t *testing.T) {
	buf := &bytes.Buffer{}
	out = buf

	Info()

	result, err := extractLogOutput(buf)

	if err != nil {
		t.Errorf(err.Error())
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	if result.At != "tteok_test.go:57" {
		t.Errorf("Expected log.At to contain \"tteok_test.go:57\". Got %s.", result.At)
	}
}

func Test_Debug_toggle(t *testing.T) {
	os.Setenv("DEBUG", "false")
	buf := &bytes.Buffer{}
	out = buf

	Debug()

	_, err := extractLogOutput(buf)

	if err == nil {
		t.Errorf("Expected calling Debug() with debug env variable set to false to produce no log")
	}
}

func Test_Levels(t *testing.T) {
	type testCase struct {
		expectedLevel string
		logFunc       func(...interface{})
	}

	os.Setenv("DEBUG", "true")

	testCases := []testCase{
		testCase{"DEBUG", Debug},
		testCase{"INFO", Info},
		testCase{"WARN", Warn},
		testCase{"FATAL", Fatal},
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		out = buf

		testCase.logFunc()

		result, err := extractLogOutput(buf)

		if err != nil {
			t.Error(err)
		}

		if result.Level != testCase.expectedLevel {
			t.Errorf("Expect Warn() to produce log with level \"%s\". Got %s", testCase.expectedLevel, result.Level)
		}
	}
}

func extractLogOutput(buf *bytes.Buffer) (log, error) {
	b, err := ioutil.ReadAll(buf)

	if err != nil {
		return log{}, err
	}

	var output log
	err = json.Unmarshal(b, &output)

	if err != nil {
		return log{}, err
	}

	return output, nil
}
