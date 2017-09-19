package goinstall

import (
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"bufio"
)

var DEBUG bool

func Output(msg string) {
	if DEBUG {
		fmt.Println(msg)
	}
}

var VarNameMatcher = regexp.MustCompile("@.*?@")

//Ask user for input
func Ask(step *Step, reader *bufio.Reader) Input {
	for {
		//setDefaultValue(step)

		prompt := step.Input.Prompt
		if len(step.Input.DefaultValue) != 0 {
			prompt = prompt + " " + step.Input.DefaultValue + " "
		}

		fmt.Print("\n", prompt)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		input := CreateInput(step, text)
		if input != nil && input.IsValid() {
			return input
		}
	}

	return nil
}

//Evaluate if a step should be executed
func Eval(conditionStr string, installVarMap map[string]string) bool {
	Output(fmt.Sprintf("condition:%v", conditionStr))

	//TODO eval '||' and '()'
	yes := true
	for _, condition := range strings.Split(conditionStr, "&&") {
		var op Operator
		if strings.Contains(condition, "==") {
			op = Equals{}
		} else if strings.Contains(condition, "!=") {
			op = NotEqual{}
		}

		yes = yes && op.Eval(condition, installVarMap)
		if !yes {
			Output("condition does not meet")
			break
		}
	}

	return yes
}

func Backup(source string, target string) {
	Output(fmt.Sprintf("Backup to: %s", target))
	cp := CopyAction{Source: source, Target: target}
	cp.Run()
}

func Save(filename string, fileString string) {
	Output(fmt.Sprintf("Write to file: %s", filename))
	ioutil.WriteFile(filename, []byte(fileString), 774)
}

func SubVar(source string, installVarMap map[string]string) string {
	resolved := source
	Output(fmt.Sprintf("source: %v", source))

	var varName string
	for varName = VarNameMatcher.FindString(resolved); varName != ""; varName = VarNameMatcher.FindString(resolved) {
		Output(fmt.Sprintf("var name : %v", varName))

		varValue := installVarMap[varName[1:len(varName)-1]]
		Output(fmt.Sprintf("var value: %v", varValue))

		resolved = strings.Replace(resolved, varName, varValue, 1)
		Output(fmt.Sprintf("resolved: %v", resolved))
	}

	return resolved
}

func GetActionType(step *Step) (string, map[string]*json.RawMessage) {
	var objMap map[string]*json.RawMessage
	var actionType string

	if actionBytes, _ := json.Marshal(step.Action); actionBytes != nil {
		var err error
		if err = json.Unmarshal(actionBytes, &objMap); err == nil && len(objMap) != 0 {
			if err = json.Unmarshal(*objMap["type"], &actionType); err == nil {
				Output(fmt.Sprintf("action type: %v", actionType))
			}
		} else {
			//Output(err.Error())
		}
	}

	return actionType, objMap
}

func HandleCopyAction(objMap map[string]*json.RawMessage, installVarMap map[string]string) {
	var source, target string
	json.Unmarshal(*objMap["source"], &source)
	json.Unmarshal(*objMap["target"], &target)

	//Replace the variables
	source = SubVar(source, installVarMap)
	target = SubVar(target, installVarMap)

	if target[len(target)-1:] == "/" {
		//Get the filename
		filename := source[strings.LastIndex(source, "/")+1:]
		Output(fmt.Sprintf("filename: %v", filename))

		target = target + filename
	}

	cp := CopyAction{Source: source, Target: target}
	cp.Run()
}

func GetVar(str string) string {
	return VarNameMatcher.FindString(str)
}

func CreateInput(step *Step, inputString string) Input {
	if step.Input.Type == "dir" {
		return DirInput{Value: strings.Replace(inputString, "\\", "/", -1)}
	} else if step.Input.Type == "select" {
		values := strings.Split(step.Input.Values, ",")
		return SelectInput{Value: inputString, Values: values}
	} else if step.Input.Type == "string" {
		return StringInput{Value: inputString, DefaultValue: step.Input.DefaultValue}
	} else if step.Input.Type == "int" {
		return IntegerInput{Value: inputString, DefaultValue: step.Input.DefaultValue}
	} else if step.Input.Type == "bool" {
		return BoolInput{Value: inputString}
	} else if step.Input.Type == "file" {
		return FileInput{Value: strings.Replace(inputString, "\\", "/", -1)}
	}
	return nil
}
