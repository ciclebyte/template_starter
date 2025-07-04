// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/ciclebyte/template_starter/internal/dao/internal"
)

// internalLanguagesDao is internal type for wrapping internal DAO implements.
type internalLanguagesDao = *internal.LanguagesDao

// languagesDao is the data access object for table languages.
// You can define custom methods on it to extend its functionality as you wish.
type languagesDao struct {
	internalLanguagesDao
}

var (
	// Languages is globally public accessible object for table languages operations.
	Languages = languagesDao{
		internal.NewLanguagesDao(),
	}
)

// Fill with you ideas below.
