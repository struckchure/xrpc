package clients

import "github.com/struckchure/go-trpc"

type TypeScriptClientConfig struct {
	Spec     trpc.TRPCSpec
	Output   string
	PostHook func()
}

func GenerateTypeScriptClient(cfg TypeScriptClientConfig) error {
	return nil
}
