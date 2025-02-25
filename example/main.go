package main

import (
	"github.com/struckchure/go-trpc"
)

type Post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ListPostInput struct {
	Skip  *int `query:"skip" json:"skip"`
	Limit *int `query:"limit" json:"limit"`
}

type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetPostInput struct {
	Id       *int    `query:"id" json:"id"`
	AuthorId *string `query:"author_id" json:"author_id"`
}

func main() {
	t := trpc.InitTRPC(trpc.InitTRPCConfig{
		Name:            "TRPC Example",
		ServerUrl:       "localhost:9090",
		SpecPath:        "./example/trpc.yaml",
		AutoGenTRPCSpec: true,
	})

	t.Use(func(c trpc.Context[any, any]) error {
		c.Locals("userId", "1290")

		return c.Next()
	})

	t.Router("post",
		trpc.NewProcedure[ListPostInput, []Post]("list").
			Input(trpc.NewValidator().
				Field("Skip", trpc.Number().Min(0).Required()).
				Field("Limit", trpc.Number().Max(10)),
			).
			Query(func(c trpc.Context[ListPostInput, []Post]) error {
				return c.Json(200, []Post{})
			}),

		trpc.NewProcedure[CreatePostInput, *Post]("create").
			Input(trpc.NewValidator().
				Field("Title", trpc.String().MinLength(10)).
				Field("Content", trpc.String().MinLength(10)),
			).
			Mutation(func(c trpc.Context[CreatePostInput, *Post]) error {
				return c.Json(201, &Post{})
			}),

		trpc.NewProcedure[GetPostInput, *Post]("get").
			Input(trpc.NewValidator().
				Field("Id", trpc.Number().Required()).
				Field("AuthorId", trpc.String().Required()),
			).
			Query(func(c trpc.Context[GetPostInput, *Post]) error {
				return c.Json(200, &Post{Title: c.Locals("userId").(string)})
			}),
	)
	t.Start(9090)
}
