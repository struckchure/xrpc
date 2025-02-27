package main

import (
	"fmt"

	"github.com/samber/do"
	"github.com/struckchure/xrpc"
)

func main() {
	t := xrpc.NewXRPC(xrpc.XRPCConfig{
		Name:            "Post Service",
		ServerUrl:       "http://localhost:9090",
		SpecPath:        "./example/basic-server/xrpc.yaml",
		AutoGenTRPCSpec: true,
	})

	do.Provide(t.Injector(), NewCarService)
	do.Provide(t.Injector(), NewEngineService)

	t.Use(func(c xrpc.Context[any, any]) error {
		c.Locals("userId", "1290")

		return c.Next()
	})

	t.Router("post",
		xrpc.NewProcedure[ListPostInput, []Post]("list").
			Use(
				func(c xrpc.Context[ListPostInput, []Post]) error {
					fmt.Println("Middleware 1")

					return c.Next()
					// return &xrpc.TRPCError{Code: 401, Detail: "something went wrong"}
				},
				func(c xrpc.Context[ListPostInput, []Post]) error {
					fmt.Println("Middleware 2")

					c.Locals("m2", true)

					return c.Next()
				},
			).
			Input(xrpc.NewValidator().
				Field("Skip", xrpc.Number().Min(0).Required()).
				Field("Limit", xrpc.Number().Max(10)),
			).
			Query(func(c xrpc.Context[ListPostInput, []Post]) error {
				fmt.Println(c.Locals("m2"))
				carService := do.MustInvoke[*CarService](c.Injector)
				carService.Start()

				return c.Json(200, []Post{})
			}),

		xrpc.NewProcedure[CreatePostInput, *Post]("create").
			Input(xrpc.NewValidator().
				Field("Title", xrpc.String().MinLength(10)).
				Field("Content", xrpc.String().MinLength(10)),
			).
			Mutation(func(c xrpc.Context[CreatePostInput, *Post]) error {
				return c.Json(201, &Post{})
			}),

		xrpc.NewProcedure[GetPostInput, *Post]("get").
			Input(xrpc.NewValidator().
				Field("Id", xrpc.Number().Required()).
				Field("AuthorId", xrpc.String().Required()),
			).
			Query(func(c xrpc.Context[GetPostInput, *Post]) error {
				return c.Json(200, &Post{Title: c.Locals("userId").(string)})
			}),
	)
	t.Start(9090)
}
