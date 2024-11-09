package utils

import "gorm.io/gorm"

type Preload struct {
	AssociationTable string
	Argument         func(db *gorm.DB) *gorm.DB
}
