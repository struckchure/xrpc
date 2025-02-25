package trpc

type TRPCSpecProcedureType string

const (
	TRPCSpecProcedureTypeQuery    TRPCSpecProcedureType = "Query"
	TRPCSpecProcedureTypeMutation TRPCSpecProcedureType = "Mutation"
)

type TRPCSpecProcedure struct {
	Path   string                `yaml:"path"`
	Type   TRPCSpecProcedureType `yaml:"type"`
	Input  TypeDescriptor        `yaml:"input"`
	Output TypeDescriptor        `yaml:"output"`
}

type TRPCSpec struct {
	Name       string              `yaml:"name"`
	ServerUrl  string              `yaml:"server_url"`
	Procedures []TRPCSpecProcedure `yaml:"procedures"`
}
