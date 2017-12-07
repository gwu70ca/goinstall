package goinstall

import (
	"testing"
	"fmt"
)

func TestEquals(t *testing.T) {
	DEBUG = true
	op := Equals{}

	installVarMap := make(map[string]string)
	condition := "os_type==windows"

	//Var is not yet set, this should evaluate to false
	if op.Eval(condition, installVarMap) {
		t.Error(fmt.Sprintf("[%v] Should be false, var is not set yet", condition))
	}

	//Set var
	installVarMap["os_type"] = "windows"
	if !op.Eval(condition, installVarMap) {
		t.Error(fmt.Sprintf("[%v] Should be true", condition))
	}

	installVarMap["os_type"] = "linux"
	if op.Eval(condition, installVarMap) {
		t.Error(fmt.Sprintf("[%v] Should be false", condition))
	}
}

func TestNotEqual(t *testing.T) {
	DEBUG = true
	op := NotEqual{}

	installVarMap := make(map[string]string)
	condition := "os_type!=linux"

	if !op.Eval(condition, installVarMap) {
		t.Error(fmt.Sprintf("[%v] Should be false, var is not set yet", condition))
	}

	installVarMap["os_type"] = "windows"

	if !op.Eval(condition, installVarMap) {
		t.Error(fmt.Sprintf("[%v] Should be false", condition))
	}

	installVarMap["os_type"] = "linux"

	if op.Eval(condition, installVarMap) {
		t.Error(fmt.Sprintf("[%v] Should be true", condition))
	}
}