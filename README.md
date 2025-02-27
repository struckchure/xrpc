<p align="center">
  <img src="https://github.com/user-attachments/assets/f491cd2b-f9c7-47ef-9ed5-03d9634fa928" />
</p>

# xRPC

This project was inspired by [TRPC](https://trpc.io) and [Zod](https://zod.dev). It provides a simple and flexible way to define and validate input for HTTP procedures using [Echo](https://echo.labstack.com) as the HTTP server.

## Features

- **Validation**: Define validation rules for your input types.
- **Procedures**: Create query and mutation procedures with type-safe input and output.
- **Echo Integration**: Easily integrate with Echo to handle HTTP requests and responses.

## Installation

To install the library, run:

```bash
go get github.com/struckchure/xrpc
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
  t := xrpc.InitTRPC()

  listPostsProcedure := xrpc.NewProcedure[ListPostInput, []Post]("list").
    Input(xrpc.NewValidator().
      Field("Skip", xrpc.Number().Min(0).Required()).
      Field("Limit", xrpc.Number().Max(10)),
    ).
    Query(func(c xrpc.Context[ListPostInput, []Post]) error {
      return c.Json(200, []Post{})
    })

  createPostProcedure := xrpc.NewProcedure[CreatePostInput, *Post]("create").
    Input(xrpc.NewValidator().
      Field("Title", xrpc.String().MinLength(10)).
      Field("Content", xrpc.String().MinLength(10)),
    ).
    Mutation(func(c xrpc.Context[CreatePostInput, *Post]) error {
      return c.Json(201, &Post{})
    })

  getPostProcedure := xrpc.NewProcedure[GetPostInput, *Post]("get").
    Input(xrpc.NewValidator().
      Field("Id", xrpc.Number().Required()).
      Field("AuthorId", xrpc.String()),
    ).
    Query(func(c xrpc.Context[GetPostInput, *Post]) error {
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
validator := xrpc.NewValidator().
  Field("Name", xrpc.String().MinLength(3).MaxLength(50)).
  Field("Email", xrpc.String().Email()).
  Field("Age", xrpc.Number().Min(18).Max(65))

err := validator.Validate(input)
if err != nil {
  // handle validation error
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [TRPC](https://xrpc.io)
- [Zod](https://zod.dev)
- [Echo](https://echo.labstack.com)
