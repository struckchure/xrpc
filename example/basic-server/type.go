package main

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
	Id       int    `query:"id" json:"id"`
	AuthorId string `query:"author_id" json:"author_id"`
}
