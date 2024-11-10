package sqltx

import "gorm.io/gorm"

func (s *sqltxImpl) BeginTX() *gorm.DB {
	return s.sql.Leader().Begin()
}

func (s *sqltxImpl) CommitTX(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (s *sqltxImpl) RollbackTX(tx *gorm.DB) error {
	return tx.Rollback().Error
}
