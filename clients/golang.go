package clients

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/samber/lo"
	"github.com/struckchure/xrpc"
)

type GolangClientConfig struct {
	Spec     xrpc.TRPCSpec
	PkgName  string
	Output   string
	PostHook func()
}

func GenerateGolangClient(cfg GolangClientConfig) error {
	cfg.Spec.Name = strings.ToLower(strings.Join(strings.Split(cfg.Spec.Name, " "), "_"))
	if len(cfg.PkgName) == 0 {
		cfg.PkgName = cfg.Spec.Name
	}

	clientName := lo.PascalCase(cfg.Spec.Name) + "Client"
	f := jen.NewFile(cfg.PkgName)

	f.ImportAlias("github.com/go-resty/resty/v2", "resty")

	f.Type().Id(clientName).Struct(
		jen.Id("client").Op("*").Qual("github.com/go-resty/resty/v2", "Client"),
	)

	// Define the structToQueryParams function
	f.Func().Id("structToQueryParams").Params(jen.Id("input").Any()).Params(jen.String(), jen.Error()).Block(
		jen.List(jen.Id("data"), jen.Id("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("input")),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Lit(""), jen.Qual("fmt", "Errorf").Call(jen.Lit("failed to marshal input: %w"), jen.Id("err"))),
		),

		jen.Var().Id("mapData").Map(jen.String()).Any(),
		jen.Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("mapData")),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Lit(""), jen.Qual("fmt", "Errorf").Call(jen.Lit("failed to unmarshal JSON: %w"), jen.Id("err"))),
		),

		jen.Id("query").Op(":=").Qual("net/url", "Values").Values(),
		jen.For(jen.List(jen.Id("key"), jen.Id("value")).Op(":=").Range().Id("mapData")).Block(
			jen.If(jen.Id("value").Op("==").Nil()).Block(
				jen.Continue(),
			),

			jen.Id("strVal").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%v"), jen.Id("value")),
			jen.Id("query").Dot("Add").Call(jen.Id("key"), jen.Id("strVal")),
		),

		jen.Return(jen.Id("query").Dot("Encode").Call(), jen.Nil()),
	)

	// Define MapError type
	f.Type().Id("MapError").Map(jen.String()).Any()

	// Define Error() method
	f.Func().Params(jen.Id("m").Id("MapError")).Id("Error").Params().String().Block(
		jen.List(jen.Id("data"), jen.Id("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("m")),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Qual("fmt", "Sprintf").Call(jen.Lit("failed to marshal error map: %v"), jen.Id("err"))),
		),
		jen.Return(jen.String().Call(jen.Id("data"))),
	)

	types := []string{}

	getFields := func(fields []xrpc.FieldDescriptor) []jen.Code {
		return lo.Map(
			fields,
			func(field xrpc.FieldDescriptor, _ int) jen.Code {
				var stmt *jen.Statement
				if field.Nillable {
					stmt = jen.Id(field.Name).Op("*").Id(field.Type)
				} else {
					stmt = jen.Id(field.Name).Id(field.Type)
				}

				if field.Alias != "" {
					stmt.Tag(map[string]string{"json": field.Alias})
				}

				return stmt
			},
		)
	}
	generateReturnStatement := func(output xrpc.TypeDescriptor, outputTypeName string) *jen.Statement {
		return jen.Return().List(
			jen.Id("resp").Dot("Result").Call().Assert(
				lo.
					If(output.Array != nil, jen.Op("*").Index().Id(outputTypeName)).
					Else(jen.Op("*").Id(outputTypeName)),
			),
			jen.Nil(),
		)
	}

	for _, procedure := range cfg.Spec.Procedures {
		methodName := lo.PascalCase(strings.Join(strings.Split(procedure.Path, "/"), "_"))
		input := procedure.Input
		output := procedure.Output

		outputTypeName := output.TypeName

		if output.TypeName != "" && !lo.Contains(types, output.TypeName) {
			f.Type().Id(output.TypeName).Struct(getFields(output.Array.Fields)...)
			f.Line()

			types = append(types, output.TypeName)
		}

		if output.TypeName == "" && output.Array != nil {
			f.Line()
			f.Type().Id(output.Array.TypeName).Struct(getFields(output.Array.Fields)...)
			f.Line()

			types = append(types, output.Array.TypeName)
			outputTypeName = output.Array.TypeName
		}

		if input.TypeName != "" && !lo.Contains(types, input.TypeName) {
			f.Line()
			f.Type().Id(input.TypeName).Struct(getFields(input.Fields)...)
			f.Line()

			types = append(types, input.TypeName)
		}

		_method := f.Func().Params(jen.Id("c").Op("*").Id(clientName)).Id(methodName).
			Params(jen.Id("input").Id(input.TypeName))

		if output.Array != nil {
			_method.Params(jen.Op("*").Index().Id(outputTypeName), jen.Error())
		} else {
			_method.Params(jen.Op("*").Id(outputTypeName), jen.Error())
		}

		if procedure.Type == xrpc.XRPCSpecProcedureTypeQuery {
			_method.Block(
				jen.List(
					jen.Id("queryParams"),
					jen.Err(),
				).Op(":=").Id("structToQueryParams").Call(jen.Id("input")),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.List(
					jen.Id("resp"),
					jen.Err(),
				).Op(":=").Id("c.client").
					Dot("R").Call().
					Dot("SetQueryString").Call(jen.Id("queryParams")).
					Dot("SetError").Call(jen.Op("&").Id("MapError").Values()).
					Dot("SetResult").Call(
					lo.
						If(output.Array != nil, jen.Op("&").Index().Id(outputTypeName).Values()).
						Else(jen.Op("&").Id(outputTypeName).Values()),
				).
					Dot("Get").Call(jen.Lit(procedure.Path)),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Line(),
				jen.If(jen.Id("resp").Dot("IsError").Call()).Block(
					jen.Return().List(
						jen.Nil(),
						jen.Id("resp").Dot("Error").Call().Assert(jen.Op("*").Id("MapError")),
					),
				),
				generateReturnStatement(output, outputTypeName),
			)
		} else if procedure.Type == xrpc.XRPCSpecProcedureTypeMutation {
			_method.Block(
				jen.List(
					jen.Id("resp"),
					jen.Err(),
				).Op(":=").Id("c.client").
					Dot("R").Call().
					Dot("SetBody").Call(jen.Id("input")).
					Dot("SetError").Call(jen.Op("&").Id("MapError").Values()).
					Dot("SetResult").Call(jen.Op("&").Id(outputTypeName).Values()).
					Dot("Post").Call(jen.Lit(procedure.Path)),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Line(),
				jen.If(jen.Id("resp").Dot("IsError").Call()).Block(
					jen.Return().List(
						jen.Nil(),
						jen.Id("resp").Dot("Error").Call().Assert(jen.Op("*").Id("MapError")),
					),
				),
				generateReturnStatement(output, outputTypeName),
			)
		}
	}

	f.Line()
	f.Func().Id("New"+clientName).Params().Op("*").Id(clientName).Block(
		jen.Id("client").Op(":=").Qual("github.com/go-resty/resty/v2", "New").Call(),
		jen.Id("client").Dot("SetBaseURL").Call(jen.Lit(cfg.Spec.ServerUrl)),
		jen.Line(),
		jen.Return(jen.Op("&").Id(clientName).Values(jen.Dict{
			jen.Id("client"): jen.Id("client"),
		})),
	)

	err := f.Save(cfg.Output)
	if err != nil {
		return err
	}

	if cfg.PostHook != nil {
		cfg.PostHook()
	}

	return nil
}
