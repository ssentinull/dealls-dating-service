package mapper

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
)

func MapBookModelToBookType(src model.BookModel) *types.Book {
	res := &types.Book{
		ID:     src.ID,
		Title:  src.Title,
		Author: src.Author,
	}

	return res
}
