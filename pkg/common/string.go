package common

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToTitle(s string) string {
	c := cases.Title(language.English)

	return c.String(s)
}
