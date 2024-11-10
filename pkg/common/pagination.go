package common

import "github.com/ssentinull/dealls-dating-service/internal/types"

func NewPagination(totalData int64, totalCurrentData int64, limit int64, page int64) types.Pagination {
	var pagination types.Pagination

	if limit <= 0 {
		return pagination
	}

	// Update Pagination
	pagination.TotalData = totalData
	totalPage := pagination.TotalData / limit

	if pagination.TotalData%limit > 0 || pagination.TotalData == 0 {
		totalPage++
	}

	pagination.CurrentPage = page
	pagination.TotalPages = totalPage
	pagination.CurrentData = totalCurrentData

	// Defining Next Page
	if pagination.CurrentPage < pagination.TotalPages {
		nextPage := pagination.CurrentPage + 1
		pagination.NextPage = &nextPage
	}

	// Defining Prev Page
	if pagination.CurrentPage > 1 {
		prevPage := pagination.CurrentPage - 1
		pagination.PrevPage = &prevPage
	}

	return pagination
}
