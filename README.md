# Go tRPC

This project was inspired by [TRPC](https://trpc.io) and [Zod](https://zod.dev). It provides a simple and flexible way to define and validate input for HTTP procedures using [Echo](https://echo.labstack.com) as the HTTP server.

## Features

- **Validation**: Define validation rules for your input types.
- **Procedures**: Create query and mutation procedures with type-safe input and output.
- **Echo Integration**: Easily integrate with Echo to handle HTTP requests and responses.

## Installation

To install the library, run:

```bash
go get github.com/struckchure/go-trpc
```

## Usage

### Define Your Types

Define the types for your input and output:

```go
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
  Id       *int   `query:"id" json:"id"`
  AuthorId string `query:"author_id" json:"author_id"`
}
```

### Create Procedures

Create procedures for listing posts, creating a post, and getting a post:

```go
func main() {
  t := trpc.InitTRPC()

  listPostsProcedure := trpc.NewProcedure[ListPostInput, []Post]("list").
    Input(trpc.NewValidator().
      Field("Skip", trpc.Number().Min(0).Required()).
      Field("Limit", trpc.Number().Max(10)),
    ).
    Query(func(c trpc.Context[ListPostInput, []Post]) error {
      return c.Json(200, []Post{})
    })

  createPostProcedure := trpc.NewProcedure[CreatePostInput, *Post]("create").
    Input(trpc.NewValidator().
      Field("Title", trpc.String().MinLength(10)).
      Field("Content", trpc.String().MinLength(10)),
    ).
    Mutation(func(c trpc.Context[CreatePostInput, *Post]) error {
      return c.Json(201, &Post{})
    })

  getPostProcedure := trpc.NewProcedure[GetPostInput, *Post]("get").
    Input(trpc.NewValidator().
      Field("Id", trpc.Number().Required()).
      Field("AuthorId", trpc.String()),
    ).
    Query(func(c trpc.Context[GetPostInput, *Post]) error {
      return c.Json(200, &Post{})
    })

  listPostsProcedure(t)
  createPostProcedure(t)
  getPostProcedure(t)

  t.Start(9090)
}
```

### Running the HTTP Server

The `main` function initializes the TRPC instance, defines the procedures, and starts the Echo HTTP server on port 9090.

## Custom Validation Library

The custom validation library provides a fluent API for defining validation rules. The library supports validation for various types such as strings and numbers, and allows specifying custom error messages.

### Example Usage

```go
validator := trpc.NewValidator().
  Field("Name", trpc.String().MinLength(3).MaxLength(50)).
  Field("Email", trpc.String().Email()).
  Field("Age", trpc.Number().Min(18).Max(65))

err := validator.Validate(input)
if err != nil {
  // handle validation error
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [TRPC](https://trpc.io)
- [Zod](https://zod.dev)
- [Echo](https://echo.labstack.com)
