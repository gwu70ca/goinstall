package goinstall

import (
	"fmt"
	"strings"
	//"io/ioutil"
	"encoding/json"
	"io/ioutil"
)

type Steps struct {
	Steps []Step
	//other properties
}

type Step struct {
	Seq       int             `json:"seq"`
	Category  string          `json:"category"`
	Name      string          `json:"name"`
	Var       string          `json:"var"`
	Condition string          `json:"condition"`
	Enabled   int             `json:"enabled"`
	Input     InstallInput    `json:"input"`
	Action    json.RawMessage `json:"action"`
}

//An input represents the data user provided. It can be bool, integer, string, directory/file name, selection from a
//pre-defined value set
type InstallInput struct {
	Type         string `json:"type"`
	Prompt       string `json:"prompt"`
	Values       string `json:"values"`
	DefaultValue string `json:"default_value"`
}

//An action is something user wants to execute, for example, copy/move file
type InstallAction interface {
	//Type   string `json:"type"`
	Run() (bool, error)
}

//Copy action
type CopyAction struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (copy CopyAction) Run() (bool, error) {
	Output(fmt.Sprintf("copying from \n\t%v \n\tto\n\t%v", copy.Source, copy.Target))
	data, err := ioutil.ReadFile(copy.Source)
	if err != nil {
		Output(fmt.Sprintf("Failed to read file %v", copy.Source))
		return false, err
	}

	if err = ioutil.WriteFile(copy.Target, data, 0774); err != nil {
		Output(fmt.Sprintf("Failed to read file %v", copy.Target))
		return false, err
	}
	return true, nil
}

//Installer runs one step a time
func (s Step) String() string {
	return fmt.Sprintf("#%d, category: %s, type: %s, values: %s\nprompt: %s", s.Seq, s.Category, s.Input.Type, s.Input.Values, s.Input.Prompt)
}

//Evaluate if a step should be executed
type Operator interface {
	Eval(condition string, installVarMap map[string]string) bool
}

type Equals struct {
}

type NotEqual struct {
}

func (eq Equals) Eval(condition string, installVarMap map[string]string) bool {
	ss := strings.Split(condition, "==")
	varName := ss[0]
	varValue := installVarMap[varName]
	varExpectedValue := ss[1]

	Output(fmt.Sprintf("EQ: varName: %v, varValue:%v, expected: %s", varName, varValue, varExpectedValue))
	return varValue == "" || varValue == varExpectedValue
}

func (ne NotEqual) Eval(condition string, installVarMap map[string]string) bool {
	ss := strings.Split(condition, "!=")
	varName := ss[0]
	varValue := installVarMap[varName]
	varExpectedValue := ss[1]

	Output(fmt.Sprintf("NE: varName: %v, varValue:%v, expected: %s", varName, varValue, varExpectedValue))
	return varValue != varExpectedValue
}
