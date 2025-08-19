package mock

import (
	"fmt"
	"photosync/src/helper"
	"reflect"
	"testing"
)

type DatabaseMock struct {
	expectedQueryResults  helper.List[[][]any]
	expectedQueryErrors   helper.List[error]
	expectedQuerySQLs     helper.List[string]
	expectedQueryArgs     helper.List[[]any]
	expectedExecuteErrors helper.List[error]
	expectedExecuteSQLs   helper.List[string]
	expectedExecuteArgs   helper.List[[]any]
	t                     *testing.T
}

func NewDatabaseMock(t *testing.T) DatabaseMock {
	return DatabaseMock{t: t}
}

func (dbm *DatabaseMock) ExpectQuery(sql string, result [][]any, args []any, err error) {
	dbm.expectedQuerySQLs.Append(sql)
	dbm.expectedQueryResults.Append(result)
	dbm.expectedQueryErrors.Append(err)
	dbm.expectedQueryArgs.Append(args)
}

func (dbm *DatabaseMock) ExpectExecute(sql string, args []any, err error) {
	dbm.expectedExecuteSQLs.Append(sql)
	dbm.expectedExecuteErrors.Append(err)
	dbm.expectedExecuteArgs.Append(args)
}

func (dbm *DatabaseMock) AssertAllExpectionsSatisfied() {
	if dbm.expectedQuerySQLs.Length() != 0 || dbm.expectedExecuteSQLs.Length() != 0 {
		fmt.Println("Not all expects satisfied!")
		dbm.t.FailNow()
	}
}

func (dbm *DatabaseMock) Execute(sql string, args ...any) error {
	if dbm.expectedExecuteSQLs.Length() == 0 {
		fmt.Println("Unexpected call!")
		dbm.t.FailNow()
	}

	expectedSQL := dbm.expectedExecuteSQLs.PopFirst()
	if sql != expectedSQL {
		fmt.Println("Unexpected sql!")
		dbm.t.FailNow()
	}

	expectedExecuteArgs := dbm.expectedExecuteArgs.PopFirst()
	if !reflect.DeepEqual(args, expectedExecuteArgs) {
		fmt.Println("Unexpected args!")
		dbm.t.FailNow()
	}

	return dbm.expectedExecuteErrors.PopFirst()
}

func (dbm *DatabaseMock) Query(sql string, args ...any) ([][]any, error) {
	if dbm.expectedQuerySQLs.Length() == 0 {
		fmt.Println("Unexpected call!")
		dbm.t.FailNow()
	}

	expectedSQL := dbm.expectedQuerySQLs.PopFirst()
	if sql != expectedSQL {
		fmt.Println("Unexpected sql!")
		dbm.t.FailNow()
	}

	expectedQueryArgs := dbm.expectedQueryArgs.PopFirst()
	if !reflect.DeepEqual(args, expectedQueryArgs) {
		fmt.Println("Unexpected args!")
		dbm.t.FailNow()
	}

	return dbm.expectedQueryResults.PopFirst(), dbm.expectedQueryErrors.PopFirst()
}
