package internals

import (
	"bytes"
	"fmt"
	"strings"
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
	Default string   // e.g., "ky"
	Names   []string // named imports like "useState", "useEffect"
	Module  string   // e.g., "react"
}

func (i *TSImport) Render() string {
	parts := []string{}

	if i.Default != "" {
		parts = append(parts, i.Default)
	}

	if len(i.Names) > 0 {
		parts = append(parts, fmt.Sprintf("{ %s }", strings.Join(i.Names, ", ")))
	}

	return fmt.Sprintf("import %s from \"%s\";", strings.Join(parts, ", "), i.Module)
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
