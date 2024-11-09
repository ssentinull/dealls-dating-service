package model

import (
	params "github.com/ssentinull/dealls-dating-service/internal/types/book"
)

type (
	BookModel struct {
		ID     int64
		Title  string
		Author string
	}

	CreateBookParams struct {
		params.CreateBookParams
	}
)

func (b BookModel) TableName() string {
	return "books"
}
