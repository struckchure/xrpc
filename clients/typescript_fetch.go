package clients

import (
	"github.com/samber/lo"
	"github.com/struckchure/xrpc"
	"github.com/struckchure/xrpc/internals"
)

var goToTSType = map[string]string{
	"string":      "string",
	"int":         "number",
	"int8":        "number",
	"int16":       "number",
	"int32":       "number",
	"int64":       "number",
	"uint":        "number",
	"uint8":       "number",
	"uint16":      "number",
	"uint32":      "number",
	"uint64":      "number",
	"float32":     "number",
	"float64":     "number",
	"bool":        "boolean",
	"interface{}": "any",
}

// Convert a Go type to TypeScript type, handling arrays and maps
func convertGoTypeToTS(goType xrpc.TypeDescriptor) string {
	if goType.Array != nil {
		return convertGoTypeToTS(*goType.Array) + "[]"
	}
	tsType, exists := goToTSType[goType.TypeName]
	if exists {
		return tsType
	}
	return goType.TypeName // Fallback to using the struct name as the TypeScript type
}

type TypeScriptClientConfig struct {
	Spec     xrpc.TRPCSpec
	Output   string
	PostHook func()
}

func GenerateTypeScriptFetchClient(cfg TypeScriptClientConfig) error {
	file := &internals.TSFile{}

	types := map[string]bool{} // Track declared types

	for _, procedure := range cfg.Spec.Procedures {
		// Convert Input Type
		inputTypeName := lo.PascalCase(procedure.Input.TypeName)
		inputFields := map[string]string{}
		for _, field := range procedure.Input.Fields {
			if field.Nillable {
				inputFields[field.Alias+"?"] = convertGoTypeToTS(xrpc.TypeDescriptor{TypeName: field.Type})
			} else {
				inputFields[field.Alias] = convertGoTypeToTS(xrpc.TypeDescriptor{TypeName: field.Type})
			}
		}

		if _, exists := types[inputTypeName]; !exists {
			file.AddNode(&internals.TSInterface{Name: inputTypeName, Fields: inputFields})
			types[inputTypeName] = true
		}

		// Convert Output Type
		outputTypeName := lo.PascalCase(procedure.Output.TypeName)
		if procedure.Output.Array != nil {
			outputTypeName = lo.PascalCase(procedure.Output.Array.TypeName) + "[]"
		}

		if _, exists := types[outputTypeName]; !exists {
			outputFields := map[string]string{}
			for _, field := range procedure.Output.Fields {
				if field.Nillable {
					outputFields[field.Alias+"?"] = convertGoTypeToTS(xrpc.TypeDescriptor{TypeName: field.Type})
				} else {
					outputFields[field.Alias] = convertGoTypeToTS(xrpc.TypeDescriptor{TypeName: field.Type})
				}
			}
			if procedure.Output.TypeName != "" {
				file.AddNode(&internals.TSInterface{Name: lo.PascalCase(procedure.Output.TypeName), Fields: outputFields})
				types[outputTypeName] = true
			}
		}

		// Define function params and body
		var params map[string]string
		var body []string

		if procedure.Type == xrpc.XRPCSpecProcedureTypeQuery {
			body = []string{
				"const queryParams = new URLSearchParams(data as unknown as Record<string, any>).toString();",
				"const response = await fetch(`" + cfg.Spec.ServerUrl + procedure.Path + "?${queryParams}`);",
				"return response.json();",
			}
			params = map[string]string{"data": inputTypeName}
		} else {
			body = []string{
				"const response = await fetch(\"" + cfg.Spec.ServerUrl + procedure.Path + "\", {",
				"  method: \"POST\",",
				"  headers: { 'Content-Type': 'application/json' },",
				"  body: JSON.stringify(data)",
				"});",
				"return response.json();",
			}
			params = map[string]string{"data": inputTypeName}
		}

		file.AddNode(&internals.TSFunction{
			Name:       lo.PascalCase(procedure.Path),
			ReturnType: "Promise<" + outputTypeName + ">",
			Params:     params,
			Body:       body,
		})
	}

	err := xrpc.WriteFile(cfg.Output, file.Render())
	if err != nil {
		return err
	}

	if cfg.PostHook != nil {
		cfg.PostHook()
	}

	return nil
}
