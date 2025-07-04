// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"github.com/ciclebyte/template_starter/internal/dao/internal"
)

// internalTemplateVariablesDao is internal type for wrapping internal DAO implements.
type internalTemplateVariablesDao = *internal.TemplateVariablesDao

// templateVariablesDao is the data access object for table template_variables.
// You can define custom methods on it to extend its functionality as you wish.
type templateVariablesDao struct {
	internalTemplateVariablesDao
}

var (
	// TemplateVariables is globally public accessible object for table template_variables operations.
	TemplateVariables = templateVariablesDao{
		internal.NewTemplateVariablesDao(),
	}
)

// Fill with you ideas below.
