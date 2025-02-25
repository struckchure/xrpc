package internals

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type TSNode interface {
	Render() string
}

type TSFile struct {
	nodes []TSNode
}

func (f *TSFile) AddNode(node TSNode) {
	f.nodes = append(f.nodes, node)
}

func (f *TSFile) Render() string {
	var buf bytes.Buffer
	for _, node := range f.nodes {
		buf.WriteString(node.Render() + "\n\n")
	}
	return buf.String()
}

type TSImport struct {
	Module string
	Names  []string
}

func (i *TSImport) Render() string {
	return fmt.Sprintf("import { %s } from \"%s\";", strings.Join(i.Names, ", "), i.Module)
}

type TSInterface struct {
	Name   string
	Fields map[string]string
}

func (iface *TSInterface) Render() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("interface %s {\n", iface.Name))
	for field, typ := range iface.Fields {
		sb.WriteString(fmt.Sprintf("  %s: %s;\n", field, typ))
	}
	sb.WriteString("}")
	return sb.String()
}

type TSFunction struct {
	Name       string
	ReturnType string
	Params     map[string]string
	Body       []string
}

func (fn *TSFunction) Render() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("async function %s(", fn.Name))
	paramList := []string{}
	for param, typ := range fn.Params {
		paramList = append(paramList, fmt.Sprintf("%s: %s", param, typ))
	}
	sb.WriteString(strings.Join(paramList, ", ") + "): " + fn.ReturnType + " {\n")
	for _, stmt := range fn.Body {
		sb.WriteString("  " + stmt + "\n")
	}
	sb.WriteString("}")
	return sb.String()
}

type TRPCSpecProcedureType string

const (
	TRPCSpecProcedureTypeQuery    TRPCSpecProcedureType = "Query"
	TRPCSpecProcedureTypeMutation TRPCSpecProcedureType = "Mutation"
)

type TRPCSpecProcedure struct {
	Path   string                `yaml:"path"`
	Type   TRPCSpecProcedureType `yaml:"type"`
	Input  map[string]string     `yaml:"input"`
	Output map[string]string     `yaml:"output"`
}

type TRPCSpec struct {
	Name       string              `yaml:"name"`
	ServerUrl  string              `yaml:"server_url"`
	Procedures []TRPCSpecProcedure `yaml:"procedures"`
}

func GenerateTSClient(spec TRPCSpec) string {
	file := &TSFile{}
	file.AddNode(&TSImport{Module: "", Names: []string{""}}) // Remove Axios import

	for _, proc := range spec.Procedures {
		ifaceName := fmt.Sprintf("%sInput", lo.PascalCase(proc.Path))
		file.AddNode(&TSInterface{Name: ifaceName, Fields: proc.Input})

		funcName := lo.PascalCase(proc.Path)
		params := map[string]string{"data": ifaceName}
		returnType := "Promise<" + lo.PascalCase(proc.Path) + "Output>"
		body := []string{
			"const response = await fetch(\"" + spec.ServerUrl + "/" + proc.Path + "\", {",
			"  method: \"POST\",",
			"  headers: { 'Content-Type': 'application/json' },",
			"  body: JSON.stringify(data)",
			"});",
			"return response.json();",
		}
		file.AddNode(&TSFunction{Name: funcName, ReturnType: returnType, Params: params, Body: body})
	}

	return file.Render()
}

func main() {
	spec := TRPCSpec{
		Name:      "ExampleAPI",
		ServerUrl: "https://api.example.com",
		Procedures: []TRPCSpecProcedure{
			{
				Path:  "getUser",
				Type:  TRPCSpecProcedureTypeQuery,
				Input: map[string]string{"id": "number"},
				Output: map[string]string{
					"id":   "number",
					"name": "string",
				},
			},
		},
	}

	fmt.Println(GenerateTSClient(spec))
}
