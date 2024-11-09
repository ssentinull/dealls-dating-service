package utils

import (
	"context"

	"gorm.io/gorm"
)

type Options struct {
	Transaction          *gorm.DB               // start query transaction
	Preloads             []Preload              // preload associations on query
	SkipAssociations     []string               // skip associations to be updated / created
	PurgeAssociations    []string               // delete associations on delete
	ClearAssociations    []string               // clear associations on delete, association is not deleted but all references are removed, usually in 1 to 1 relation
	DeleteAssociations   map[string]interface{} // delete associations on update
	SyncAssociations     bool                   // save associations on upadte / create
	Unscoped             bool                   // unscoped query, usually for hard delete
	UpdateZeroNullFields []string               // update selected fields that can have zero or null value
}

func (o *Options) Extract(ctx context.Context, db *gorm.DB) *gorm.DB {
	if o == nil {
		return db
	}

	if o.Transaction != nil {
		db = o.Transaction
	}

	db = db.WithContext(ctx)
	if o.Unscoped {
		db = db.Unscoped()
	}

	if len(o.SkipAssociations) == 0 && o.SyncAssociations {
		db = db.Session(&gorm.Session{FullSaveAssociations: true})
	}

	return db
}

func (o *Options) Preload(db *gorm.DB) *gorm.DB {
	if o == nil {
		return db
	}

	if o.Transaction != nil {
		db = o.Transaction
	}

	for _, preload := range o.Preloads {
		if preload.Argument != nil {
			db = db.Preload(preload.AssociationTable, preload.Argument)
		} else {
			db = db.Preload(preload.AssociationTable)
		}
	}

	return db
}
