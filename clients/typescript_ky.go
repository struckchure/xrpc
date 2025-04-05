package clients

import (
	"github.com/samber/lo"
	"github.com/struckchure/xrpc"
	"github.com/struckchure/xrpc/internals"
)

func GenerateTypeScriptKyClient(cfg TypeScriptClientConfig) error {
	file := &internals.TSFile{}

	// Add ky import
	file.AddNode(&internals.TSImport{
		Module:  "ky",
		Default: "ky",
	})

	types := map[string]bool{}

	for _, procedure := range cfg.Spec.Procedures {
		inputTypeName := lo.PascalCase(procedure.Input.TypeName)
		inputFields := map[string]string{}
		for _, field := range procedure.Input.Fields {
			goType := xrpc.TypeDescriptor{TypeName: field.Type}
			if field.Nillable {
				inputFields[field.Alias+"?"] = convertGoTypeToTS(goType)
			} else {
				inputFields[field.Alias] = convertGoTypeToTS(goType)
			}
		}

		if _, exists := types[inputTypeName]; !exists {
			file.AddNode(&internals.TSInterface{Name: inputTypeName, Fields: inputFields})
			types[inputTypeName] = true
		}

		outputTypeName := lo.PascalCase(procedure.Output.TypeName)
		if procedure.Output.Array != nil {
			outputTypeName = lo.PascalCase(procedure.Output.Array.TypeName) + "[]"
		}

		if _, exists := types[outputTypeName]; !exists {
			outputFields := map[string]string{}
			for _, field := range procedure.Output.Fields {
				goType := xrpc.TypeDescriptor{TypeName: field.Type}
				if field.Nillable {
					outputFields[field.Alias+"?"] = convertGoTypeToTS(goType)
				} else {
					outputFields[field.Alias] = convertGoTypeToTS(goType)
				}
			}
			if procedure.Output.TypeName != "" {
				file.AddNode(&internals.TSInterface{Name: lo.PascalCase(procedure.Output.TypeName), Fields: outputFields})
				types[outputTypeName] = true
			}
		}

		var params map[string]string
		var body []string
		if procedure.Type == xrpc.XRPCSpecProcedureTypeQuery {
			body = []string{
				"const queryParams = new URLSearchParams(data as unknown as Record<string, any>).toString();",
				"return await ky.get(`" + cfg.Spec.ServerUrl + procedure.Path + "?${queryParams}`).json<" + outputTypeName + ">();",
			}
			params = map[string]string{"data": inputTypeName}
		} else {
			body = []string{
				"return await ky.post(\"" + cfg.Spec.ServerUrl + procedure.Path + "\", {",
				"  json: data",
				"}).json<" + outputTypeName + ">();",
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
