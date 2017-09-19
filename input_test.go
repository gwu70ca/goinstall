package goinstall

import (
	"testing"
	"fmt"
)

func TestValidateBool(t *testing.T) {
	testValid("y", "", "Bool input '%s' should be valid", t)
	testValid("Y", "", "Bool input '%s' should be valid", t)
	testValid("n", "", "Bool input '%s' should be valid", t)
	testValid("N", "", "Bool input '%s' should be valid", t)

	testYes("y", "", "Bool input '%s' should be yes", t)
	testYes("Y", "", "Bool input '%s' should be yes", t)
	testNo("n", "", "Bool input '%s' should be no", t)
	testNo("N", "", "Bool input '%s' should be no", t)

	testInvalid("", "", "Bool input '%s' should be invalid", t)
	testInvalid("Q", "", "Bool input '%s' should be invalid", t)
	testInvalid("aaa", "", "Bool input '%s' should be invalid", t)

	//Test with default value
	testValid("", "y", "Bool input '%s' should be valid (y)", t)
	testValid("", "Y", "Bool input '%s' should be valid (Y)", t)
	testValid("", "n", "Bool input '%s' should be invalid (n)", t)
	testValid("", "N", "Bool input '%s' should be invalid (N)", t)

	testInvalid("", "", "Bool input '%s' should be invalid, no default value", t)
}

func testValid(v string, defaltValue string, msg string, t *testing.T) {
	boolInput := BoolInput{Value: v, DefaultValue: defaltValue}
	if !boolInput.IsValid() {
		t.Error(fmt.Sprintf(msg, v))
	}
}

func testInvalid(v string, defaltValue string, msg string, t *testing.T) {
	boolInput := BoolInput{Value: v, DefaultValue: defaltValue}
	if boolInput.IsValid() {
		t.Error(fmt.Sprintf(msg, v))
	}
}

func testYes(v string, defaltValue string, msg string, t *testing.T) {
	boolInput := BoolInput{Value: v, DefaultValue: defaltValue}
	if !boolInput.yes() {
		t.Error(fmt.Sprintf(msg, v))
	}
}

func testNo(v string, defaltValue string, msg string, t *testing.T) {
	boolInput := BoolInput{Value: v, DefaultValue: defaltValue}

	if boolInput.yes() {
		t.Error(fmt.Sprintf(msg, v))
	}
}
