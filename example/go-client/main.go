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

	err = clients.GenerateGolangClient(clients.GolangClientConfig{
		Spec:    spec,
		Output:  "./post_service.go",
		PkgName: "main",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// client := NewPostServiceClient()
	// postList, err := client.PostList(ListPostInput{
	// 	Skip:  lo.ToPtr(2),
	// 	Limit: lo.ToPtr(10),
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("%#v\n", postList)

	// postCreate, err := client.PostCreate(CreatePostInput{
	// 	Title:   "OneTwoThreeFourFiveSix",
	// 	Content: "OneTwoThreeFourFiveSix",
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("%#v\n", postCreate)

	// postGet, err := client.PostGet(GetPostInput{Id: lo.ToPtr(12), AuthorId: lo.ToPtr("id-1")})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("%#v\n", postGet)
}
