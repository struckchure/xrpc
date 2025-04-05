package main

import (
	"fmt"
	"os"

	"github.com/struckchure/xrpc"
	"github.com/struckchure/xrpc/clients"
	"gopkg.in/yaml.v3"
)

func main() {
	spec := xrpc.TRPCSpec{}

	data, err := os.ReadFile("../basic-server/xrpc.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = clients.GenerateTypeScriptKyClient(clients.TypeScriptClientConfig{
		Spec:   spec,
		Output: "./post_service.ts",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
