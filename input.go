package goinstall

import (
	"strings"
	"fmt"
	"os"
	"strconv"
)

type Input interface {
	IsValid() bool
	GetValue() string
}

//Integer input
type IntegerInput struct {
	Prompt       string
	Value        string
	DefaultValue string
}

func (input IntegerInput) IsValid() bool {
	Output(fmt.Sprintf("Validate integer: %s, default value: %s", input.Value, input.DefaultValue))

	if len(input.Value) == 0 {
		return input.DefaultValue != ""
	}

	if _, err := strconv.Atoi(input.Value); err != nil {
		Output(fmt.Sprintf("Invalid integer: %s", input.Value))
		return false
	}
	return true
}

func (input IntegerInput) GetValue() string {
	return input.Value
}

//String input
type StringInput struct {
	Prompt       string
	Value        string
	DefaultValue string
}

func (input StringInput) IsValid() bool {
	return len(input.Value) > 0 || input.DefaultValue != ""
}

func (input StringInput) GetValue() string {
	return input.Value
}

//Select input, choose from group of values
type SelectInput struct {
	Prompt string
	Value  string
	Values []string
}

func (input SelectInput) IsValid() bool {
	value := input.Value
	if len(value) == 0 {
		return false
	}
	value = strings.ToLower(value)
	Output(fmt.Sprintf("Validate select: %s", value))
	index, _ := strconv.Atoi(value)
	if index > 0 && index <= len(input.Values) {
		input.Value = input.Values[index-1]
		return true
	}

	/*for _, num := range input.Values {
		//fmt.Println("\tnum:",num)
		if num == value {
			valid = true
			break
		}
	}*/

	return false
}

func (input SelectInput) GetValue() string {
	index, _ := strconv.Atoi(input.Value)
	return input.Values[index-1]
}

//Boolean input
type BoolInput struct {
	Prompt       string
	Value        string
	DefaultValue string
}

func (input BoolInput) IsValid() bool {
	value := input.Value
	Output(fmt.Sprintf("Validate bool: %s, default: %s", value, input.DefaultValue))

	if len(value) == 0 && input.DefaultValue != "" {
		return true
	}
	value = strings.ToLower(value)
	Output(fmt.Sprintf("\tconverted: %s: ", value))
	return value == "y" || value == "n"
}

func (input BoolInput) yes() bool {
	return strings.ToLower(input.Value) == "y"
}

func (input BoolInput) GetValue() string {
	return strings.ToLower(input.Value)
}

//Directory input
type DirInput struct {
	Prompt       string
	Value        string
	DefaultValue string
}

func (input DirInput) IsValid() bool {
	value := input.Value
	if len(value) == 0 {
		return false
	}
	value = strings.ToLower(value)
	Output(fmt.Sprintf("Validate dir: %s", value))

	f, err := os.Stat(input.Value)
	switch {
	case err != nil:
		if os.IsNotExist(err) {
			Output(fmt.Sprintf("Dir %s does not exist", input.Value))
			return false
		}
	case !f.IsDir():
		fmt.Println("Dir ", input.Value, " is not a directory")
		return false
		//case: test permission
	}

	return true
}

func (input DirInput) GetValue() string {
	return input.Value
}

//File input
type FileInput struct {
	Prompt       string
	Value        string
	DefaultValue string
}

func (input FileInput) IsValid() bool {
	value := input.Value
	if len(value) == 0 {
		return false
	}
	value = strings.ToLower(value)
	Output(fmt.Sprintf("Validate file: [%v]", value))

	f, err := os.Stat(input.Value)
	switch {
	case err != nil:
		if os.IsNotExist(err) {
			Output(fmt.Sprintf("File [%v] does not exist", input.Value))
			return false
		}
	case f.IsDir():
		Output(fmt.Sprintf("File [%v] is not a file", input.Value))
		return false
		//TODO case: test permission
	}

	return true
}

func (input FileInput) GetValue() string {
	return input.Value
}
