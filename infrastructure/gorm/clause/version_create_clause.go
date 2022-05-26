package clause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func NewVersionCreateClause(field *schema.Field) ClauseInterface {
	return &versionCreateClause{
		field: field,
	}
}

type versionCreateClause struct {
	field *schema.Field
}

func (v versionCreateClause) Name() string {
	return ""
}

func (v versionCreateClause) Build(clause.Builder) {
}

func (v versionCreateClause) MergeClause(*clause.Clause) {
}

func (v versionCreateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.String() == "" {
		values := callbacks.ConvertToCreateValues(stmt)

		columnIdx := 0
		found := false
		for idx := range values.Columns {
			if values.Columns[idx].Name == v.field.DBName {
				columnIdx = idx
				found = true
				break
			}
		}

		if found {
			for rowIdx := range values.Values {
				values.Values[rowIdx][columnIdx] = 1
			}
		} else {
			values.Columns = append(values.Columns, clause.Column{Name: v.field.DBName})
			for rowIdx := range values.Values {
				values.Values[rowIdx] = append(values.Values[rowIdx], 1)
			}
		}

		stmt.AddClauseIfNotExists(clause.Insert{})
		stmt.AddClause(values)
		stmt.Build("INSERT", "VALUES", "ON CONFLICT")
	}
}
