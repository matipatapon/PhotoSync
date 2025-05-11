package mock

import (
	"fmt"
	"reflect"
	"testing"
)

type DatabaseMock struct {
	expectedQueryResults  [][][]any
	expectedQueryErrors   []error
	expectedQuerySQLs     []string
	expectedQueryArgs     [][]any
	expectedExecuteErrors []error
	expectedExecuteSQLs   []string
	expectedExecuteArgs   [][]any
	t                     *testing.T
}

func NewDatabaseMock(t *testing.T) DatabaseMock {
	return DatabaseMock{[][][]any{}, []error{}, []string{}, [][]any{}, []error{}, []string{}, [][]any{}, t}
}

func (dbm *DatabaseMock) ExpectQuery(sql string, result [][]any, args []any, err error) {
	dbm.expectedQuerySQLs = append(dbm.expectedQuerySQLs, sql)
	dbm.expectedQueryResults = append(dbm.expectedQueryResults, result)
	dbm.expectedQueryErrors = append(dbm.expectedQueryErrors, err)
	dbm.expectedQueryArgs = append(dbm.expectedQueryArgs, args)
}

func (dbm *DatabaseMock) ExpectExecute(sql string, args []any, err error) {
	dbm.expectedExecuteSQLs = append(dbm.expectedExecuteSQLs, sql)
	dbm.expectedExecuteErrors = append(dbm.expectedExecuteErrors, err)
	dbm.expectedExecuteArgs = append(dbm.expectedExecuteArgs, args)
}

func (dbm *DatabaseMock) AssertAllExpectionsSatisfied() {
	if len(dbm.expectedQuerySQLs) != 0 || len(dbm.expectedExecuteSQLs) != 0 {
		fmt.Println("Not all expects satisfied!")
		dbm.t.FailNow()
	}
}

func (dbm *DatabaseMock) Execute(sql string, args ...any) error {
	if len(dbm.expectedExecuteSQLs) == 0 {
		fmt.Println("Unexpected call!")
		dbm.t.FailNow()
	}

	expectedSQL := dbm.expectedExecuteSQLs[len(dbm.expectedExecuteSQLs)-1]
	dbm.expectedExecuteSQLs = dbm.expectedExecuteSQLs[:len(dbm.expectedExecuteSQLs)-1]
	if sql != expectedSQL {
		fmt.Println("Unexpected sql!")
		dbm.t.FailNow()
	}

	expectedExecuteArgs := dbm.expectedExecuteArgs[len(dbm.expectedExecuteArgs)-1]
	dbm.expectedExecuteArgs = dbm.expectedExecuteArgs[:len(dbm.expectedExecuteArgs)-1]
	if !reflect.DeepEqual(args, expectedExecuteArgs) {
		fmt.Println("Unexpected args!")
		dbm.t.FailNow()

	}

	err := dbm.expectedExecuteErrors[len(dbm.expectedExecuteErrors)-1]
	dbm.expectedExecuteErrors = dbm.expectedExecuteErrors[:len(dbm.expectedExecuteErrors)-1]

	return err
}

func (dbm *DatabaseMock) Query(sql string, args ...any) ([][]any, error) {
	if len(dbm.expectedQuerySQLs) == 0 {
		fmt.Println("Unexpected call!")
		dbm.t.FailNow()
	}

	expectedSQL := dbm.expectedQuerySQLs[len(dbm.expectedQuerySQLs)-1]
	dbm.expectedQuerySQLs = dbm.expectedQuerySQLs[:len(dbm.expectedQuerySQLs)-1]
	if sql != expectedSQL {
		fmt.Println("Unexpected sql!")
		dbm.t.FailNow()
	}

	expectedQueryArgs := dbm.expectedQueryArgs[len(dbm.expectedQueryArgs)-1]
	dbm.expectedQueryArgs = dbm.expectedQueryArgs[:len(dbm.expectedQueryArgs)-1]
	if !reflect.DeepEqual(args, expectedQueryArgs) {
		fmt.Println("Unexpected args!")
		dbm.t.FailNow()

	}

	result := dbm.expectedQueryResults[len(dbm.expectedQueryResults)-1]
	dbm.expectedQueryResults = dbm.expectedQueryResults[:len(dbm.expectedQueryResults)-1]
	err := dbm.expectedQueryErrors[len(dbm.expectedQueryErrors)-1]
	dbm.expectedQueryErrors = dbm.expectedQueryErrors[:len(dbm.expectedQueryErrors)-1]

	return result, err
}
