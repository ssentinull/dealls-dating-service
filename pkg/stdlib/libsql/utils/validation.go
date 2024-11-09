package utils

import (
	"fmt"
	"strings"

	"github.com/ssentinull/golang-boilerplate/pkg/common"
)

type SortOption map[string]struct {
	Transform   string
	Insensitive bool
}

func ValidateLimit(limit *int64) *int64 {
	min := int64(10)
	max := int64(100)
	if limit == nil {
		return &min
	}

	if *limit <= 0 {
		return &min
	}

	if *limit > max {
		return &max
	}

	return limit
}

func ValidatePage(page *int64) *int64 {
	initPage := int64(1)
	if page == nil {
		return &initPage
	}

	if *page <= 0 {
		return &initPage
	}

	return page
}

func ValidateSort(sort *string, transformer ...map[string]string) *string {
	order := "created_at ASC"
	if sort == nil {
		return &order
	}

	if *sort == "" {
		return &order
	}

	sortArr := strings.Split(*sort, " ")
	if len(sortArr) == 2 {
		if len(transformer) > 0 {
			if transformed, ok := transformer[0][sortArr[0]]; ok {
				sortArr[0] = transformed

				*sort = strings.Join(sortArr, " ")
			}
		}

		if strings.ToUpper(sortArr[1]) == "ASC" || strings.ToUpper(sortArr[1]) == "DESC" {
			return sort
		}

		return &order
	}

	return &order
}

func ValidateSortV2(rawSort *string, opts ...SortOption) string {
	defaultSort := "created_at ASC"

	if rawSort == nil {
		return defaultSort
	}

	if *rawSort == "" {
		return defaultSort
	}

	opt := common.FirstOrZero(opts)
	sortArr := strings.Split(*rawSort, ",")
	resArr := []string{}

	for _, tmpSA := range sortArr {
		sort := strings.Split(tmpSA, " ")
		if len(sort) < 1 || len(sort) > 2 {
			continue
		}

		key := sort[0]
		if len(key) == 0 {
			continue
		}

		order := "ASC" // default
		if len(sort) == 2 && (strings.ToUpper(sort[1]) == "ASC" || strings.ToUpper(sort[1]) == "DESC") {
			order = sort[1]
		}

		specOpt := opt[key]
		key = common.Fallback(specOpt.Transform, key)

		if specOpt.Insensitive {
			key = fmt.Sprintf("UPPER(%s)", key)
		}

		resArr = append(resArr, key+" "+order)
	}

	if len(resArr) == 0 {
		return defaultSort
	}

	return strings.Join(resArr, ",")
}
