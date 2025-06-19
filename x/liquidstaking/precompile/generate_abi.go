//go:build ignore
// +build ignore

// This file generates the Go bindings from the ABI JSON file
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

const abiTemplate = `// Code generated - DO NOT EDIT.
// This file is generated from the ABI JSON file.

package precompile

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ABI method names
const (
{{- range .Methods }}
	Method{{ .Name }} = "{{ .RawName }}"
{{- end }}
)

// ABI event names
const (
{{- range .Events }}
	Event{{ .Name }} = "{{ .RawName }}"
{{- end }}
)

// ABI JSON string
const ABIJSON = ` + "`" + `{{ .JSON }}` + "`" + `

// Method IDs
var (
{{- range .Methods }}
	{{ .Name }}ID = common.Hex2Bytes("{{ .ID }}")
{{- end }}
)

// Initialize ABI
func init() {
	var err error
	LiquidStakingABI, err = abi.JSON(strings.NewReader(ABIJSON))
	if err != nil {
		panic(fmt.Errorf("failed to parse liquid staking ABI: %w", err))
	}
}

// Generated method wrappers

{{- range .Methods }}
// {{ .Name }} calls the {{ .RawName }} method
func {{ .Name }}(abi abi.ABI, {{ .GoArgs }}) ([]byte, error) {
	return abi.Pack("{{ .RawName }}"{{ if .PackArgs }}, {{ .PackArgs }}{{ end }})
}

// Unpack{{ .Name }}Input unpacks the input for {{ .RawName }}
func Unpack{{ .Name }}Input(abi abi.ABI, input []byte) ({{ .GoArgs }}, error) {
	args, err := abi.Methods["{{ .RawName }}"].Inputs.Unpack(input)
	if err != nil {
		return {{ .ZeroReturns }}, err
	}
	{{ .UnpackLogic }}
	return {{ .ReturnVars }}, nil
}

// Unpack{{ .Name }}Output unpacks the output for {{ .RawName }}
func Unpack{{ .Name }}Output(abi abi.ABI, output []byte) ({{ .GoReturns }}, error) {
	returns, err := abi.Methods["{{ .RawName }}"].Outputs.Unpack(output)
	if err != nil {
		return {{ .ZeroOutputs }}, err
	}
	{{ .UnpackOutputLogic }}
	return {{ .OutputVars }}, nil
}
{{- end }}
`

type Method struct {
	Name              string
	RawName           string
	ID                string
	GoArgs            string
	PackArgs          string
	GoReturns         string
	ZeroReturns       string
	UnpackLogic       string
	ReturnVars        string
	ZeroOutputs       string
	UnpackOutputLogic string
	OutputVars        string
}

type Event struct {
	Name    string
	RawName string
	ID      string
}

type TemplateData struct {
	Methods []Method
	Events  []Event
	JSON    string
}

func main() {
	// Read ABI JSON file
	abiData, err := ioutil.ReadFile("abi.json")
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	// Parse ABI
	var abiDef []map[string]interface{}
	if err := json.Unmarshal(abiData, &abiDef); err != nil {
		log.Fatalf("Failed to parse ABI JSON: %v", err)
	}

	// Extract methods and events
	var methods []Method
	var events []Event

	for _, item := range abiDef {
		itemType, ok := item["type"].(string)
		if !ok {
			continue
		}

		name, _ := item["name"].(string)

		switch itemType {
		case "function":
			method := generateMethod(item)
			methods = append(methods, method)

		case "event":
			event := Event{
				Name:    toCamelCase(name),
				RawName: name,
			}
			events = append(events, event)
		}
	}

	// Format JSON for embedding
	var buf bytes.Buffer
	json.Indent(&buf, abiData, "", "\t")

	// Generate code
	tmpl, err := template.New("abi").Parse(abiTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	data := TemplateData{
		Methods: methods,
		Events:  events,
		JSON:    buf.String(),
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	// Format generated code
	formatted, err := format.Source(output.Bytes())
	if err != nil {
		log.Printf("Warning: failed to format generated code: %v", err)
		formatted = output.Bytes()
	}

	// Write to file
	if err := ioutil.WriteFile("abi_generated.go", formatted, 0644); err != nil {
		log.Fatalf("Failed to write generated file: %v", err)
	}

	fmt.Println("Successfully generated abi_generated.go")
}

func generateMethod(method map[string]interface{}) Method {
	name := method["name"].(string)
	inputs, _ := method["inputs"].([]interface{})
	outputs, _ := method["outputs"].([]interface{})

	m := Method{
		Name:    toCamelCase(name),
		RawName: name,
		ID:      calculateMethodID(method),
	}

	// Generate input arguments
	var goArgs []string
	var packArgs []string
	var unpackLogic []string
	var returnVars []string
	var zeroReturns []string

	for i, input := range inputs {
		inputMap := input.(map[string]interface{})
		argName := inputMap["name"].(string)
		if argName == "" {
			argName = fmt.Sprintf("arg%d", i)
		}
		argType := solTypeToGo(inputMap["type"].(string))

		goArgs = append(goArgs, fmt.Sprintf("%s %s", argName, argType))
		packArgs = append(packArgs, argName)
		
		unpackLogic = append(unpackLogic, fmt.Sprintf("%s := args[%d].(%s)", argName, i, argType))
		returnVars = append(returnVars, argName)
		zeroReturns = append(zeroReturns, zeroValue(argType))
	}

	m.GoArgs = strings.Join(goArgs, ", ")
	m.PackArgs = strings.Join(packArgs, ", ")
	m.UnpackLogic = strings.Join(unpackLogic, "\n\t")
	m.ReturnVars = strings.Join(returnVars, ", ")
	m.ZeroReturns = strings.Join(zeroReturns, ", ")

	// Generate output returns
	var goReturns []string
	var outputVars []string
	var zeroOutputs []string
	var unpackOutputLogic []string

	for i, output := range outputs {
		outputMap := output.(map[string]interface{})
		retName := outputMap["name"].(string)
		if retName == "" {
			retName = fmt.Sprintf("ret%d", i)
		}
		retType := solTypeToGo(outputMap["type"].(string))

		goReturns = append(goReturns, fmt.Sprintf("%s %s", retName, retType))
		outputVars = append(outputVars, retName)
		zeroOutputs = append(zeroOutputs, zeroValue(retType))
		
		unpackOutputLogic = append(unpackOutputLogic, fmt.Sprintf("%s := returns[%d].(%s)", retName, i, retType))
	}

	m.GoReturns = strings.Join(goReturns, ", ")
	m.OutputVars = strings.Join(outputVars, ", ")
	m.ZeroOutputs = strings.Join(zeroOutputs, ", ")
	m.UnpackOutputLogic = strings.Join(unpackOutputLogic, "\n\t")

	return m
}

func calculateMethodID(method map[string]interface{}) string {
	// This is a simplified version - in reality, you'd use the actual ABI encoder
	name := method["name"].(string)
	return fmt.Sprintf("%x", name)[:8]
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func solTypeToGo(solType string) string {
	switch solType {
	case "uint256":
		return "*big.Int"
	case "address":
		return "common.Address"
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "uint8":
		return "uint8"
	case "bytes":
		return "[]byte"
	default:
		if strings.HasSuffix(solType, "[]") {
			baseType := solType[:len(solType)-2]
			return "[]" + solTypeToGo(baseType)
		}
		if strings.HasPrefix(solType, "tuple") {
			// Handle struct types
			return "interface{}"
		}
		return solType
	}
}

func zeroValue(goType string) string {
	switch goType {
	case "*big.Int":
		return "nil"
	case "common.Address":
		return "common.Address{}"
	case "string":
		return `""`
	case "bool":
		return "false"
	case "uint8":
		return "0"
	case "[]byte":
		return "nil"
	default:
		if strings.HasPrefix(goType, "[]") {
			return "nil"
		}
		return "nil"
	}
}