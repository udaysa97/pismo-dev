package utils

import (
	"fmt"
	types "pismo-dev/internal/repository/errors"

	sql "github.com/Masterminds/squirrel"
)

type IQuery interface {
	ToSql() (string, interface{}, error)
}

func QueryFormatter(query interface{}) string {
	var qs string
	var qd interface{}
	var qerr error

	if qi, ok := query.(sql.SelectBuilder); ok {
		qs, qd, qerr = qi.ToSql()
	} else if qi, ok := query.(sql.InsertBuilder); ok {
		qs, qd, qerr = qi.ToSql()
	} else if qi, ok := query.(sql.UpdateBuilder); ok {
		qs, qd, qerr = qi.ToSql()
	} else if qi, ok := query.(sql.DeleteBuilder); ok {
		qs, qd, qerr = qi.ToSql()
	} else if qi, ok := query.(sql.CaseBuilder); ok {
		qs, qd, qerr = qi.ToSql()
	} else if qi, ok := query.(string); ok {
		qs = qi
	} else {
		return types.ErrToSQLStatement.Error()
	}

	if qerr != nil {
		return types.ErrToSQLStatement.Error()
	}

	return fmt.Sprintf(types.ErrSQLStatement.Error(), qs, qd)
}
