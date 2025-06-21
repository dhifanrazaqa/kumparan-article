package models

import "time"

type Article struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Body      string        `json:"body"`
	AuthorID  string        `json:"authorId"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Author    *UserResponse `json:"author,omitempty"`
}

type CreateArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdateArticleRequest struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type ListArticlesParams struct {
	Query  string
	Author string
	Limit  int
	Offset int
}
