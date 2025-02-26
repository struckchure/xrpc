package xrpc

type XRPCSpecProcedureType string

const (
	XRPCSpecProcedureTypeQuery    XRPCSpecProcedureType = "Query"
	XRPCSpecProcedureTypeMutation XRPCSpecProcedureType = "Mutation"
)

type XRPCSpecProcedure struct {
	Path   string                `yaml:"path"`
	Type   XRPCSpecProcedureType `yaml:"type"`
	Input  TypeDescriptor        `yaml:"input"`
	Output TypeDescriptor        `yaml:"output"`
}

type TRPCSpec struct {
	Name       string              `yaml:"name"`
	ServerUrl  string              `yaml:"server_url"`
	Procedures []XRPCSpecProcedure `yaml:"procedures"`
}
