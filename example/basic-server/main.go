package main

import (
	"fmt"

	"github.com/samber/do"
	"github.com/struckchure/xrpc"
	"github.com/struckchure/xrpc/validation"
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

		return nil
	})

	t.Router("post",
		xrpc.NewProcedure[ListPostInput, []Post]("list").
			Use(
				func(c xrpc.Context[ListPostInput, []Post]) error {
					fmt.Println("Middleware 1")

					// return &xrpc.XRPCError{Code: 401, Detail: "something went wrong"}
					return nil
				},
				func(c xrpc.Context[ListPostInput, []Post]) error {
					fmt.Println("Middleware 2")

					c.Locals("m2", true)

					return nil
				},
			).
			Input(validation.NewValidator().
				Field("Skip", validation.Int().Min(0).Required()).
				Field("Limit", validation.Int().Max(10)),
			).
			Query(func(c xrpc.Context[ListPostInput, []Post]) error {
				fmt.Println(c.Locals("m2"))
				fmt.Println(c.Locals("userId"))
				carService := do.MustInvoke[*CarService](c.Injector)
				carService.Start()

				return c.Json(200, []Post{})
			}),

		xrpc.NewProcedure[CreatePostInput, *Post]("create").
			Input(validation.NewValidator().
				Field("Title", validation.String().MinLength(10)).
				Field("Content", validation.String().MinLength(10)),
			).
			Mutation(func(c xrpc.Context[CreatePostInput, *Post]) error {
				return c.Json(201, &Post{})
			}),

		xrpc.NewProcedure[GetPostInput, *Post]("get").
			Input(validation.NewValidator().
				Field("Id", validation.Int().Required()).
				Field("AuthorId", validation.String().Required()),
			).
			Query(func(c xrpc.Context[GetPostInput, *Post]) error {
				return c.Json(200, &Post{Title: c.Locals("userId").(string)})
			}),
	)
	t.Start(9090)
}
