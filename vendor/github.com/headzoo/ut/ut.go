// Package ut has assert methods to make testing easier.
package ut

import (
	"reflect"
	"strings"
	"testing"
	"runtime"
	"bytes"
	"fmt"
)

// Test is the testing.T instance for the test being run.
var test *testing.T = &testing.T{}

// Run sets up an individual test.
func Run(t *testing.T) {
	test = t
}

// Test returns the testing.T passed to the Run() method.
func Test() *testing.T {
	return test
}

// AssertTrue tests whether the given value is true.
func AssertTrue(actual bool) bool {
	if !actual {
		Errorf("Failed asserting %q is true.", actual)
		return false
	}
	return true
}

// AssertFalse tests whether the given value is false.
func AssertFalse(actual bool) bool {
	if actual {
		Errorf("Failed asserting %q is false.", actual)
		return false
	}
	return true
}

// AssertNotNil tests whether the given value is not nil.
func AssertNotNil(actual interface{}) bool {
	if actual == nil {
		Errorf("Failed asserting the value is not nil.")
		return false
	}
	return true
}

// AssertNil tests whether the given value is nil.
func AssertNil(actual interface{}) bool {
	if actual != nil {
		Errorf("Failed asserting %T is nil.", actual)
		return false
	}
	return true
}

// AssertEmpty tests whether the given value is empty.
func AssertEmpty(actual interface{}) bool {
	t := reflect.ValueOf(actual)
	if !isZero(t) {
		Errorf("Failed asserting %q is empty.", actual)
		return false
	}
	return true
}

// AssertNotEmpty tests whether the given value is not empty.
func AssertNotEmpty(actual interface{}) bool {
	t := reflect.ValueOf(actual)
	if isZero(t) {
		Errorf("Failed asserting %q is not empty.", actual)
		return false
	}
	return true
}

// AssertEquals tests whether two values are equal.
func AssertEquals(expected, actual interface{}) bool {
	if !reflect.DeepEqual(expected, actual) {
		Errorf("Failed asserting %q equals %q.", expected, actual)
		return false
	}
	return true
}

// AssertNotEquals tests whether two values do not equal each other.
func AssertNotEquals(expected, actual interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		Errorf("Failed asserting %q is not equal to %q.", expected, actual)
		return false
	}
	return true
}

// AssertGreaterThan tests whether the actual value is greater than the expected value.
func AssertGreaterThan(expected, actual int) bool {
	if expected >= actual {
		Errorf("Failed asserting %q is greater than %q.", actual, expected)
		return false
	}
	return true
}

// AssertContains tests whether the expected value contains the actual value.
func AssertContains(expected, actual string) bool {
	if !strings.Contains(actual, expected) {
		Errorf("Failed asserting %q contains %q.", actual, expected)
		return false
	}
	return true
}

// isZero returns if the value is zero.
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

// Hack.
func Errorf(format string, args ...interface{}) {
	format = decorate(format)
	test.Log(fmt.Sprintf(format, args...))
	test.Fail()
}

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
func decorate(s string) string {
	_, file, line, ok := runtime.Caller(3) // decorate + log + public function.
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteByte('\t')
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}
