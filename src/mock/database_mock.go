package mock

import (
	"fmt"
	"reflect"
	"testing"
)

type DatabaseMock struct {
	expectedResults [][][]any
	expectedErrors  []error
	expectedSQLs    []string
	expectedArgs    [][]any
	t               *testing.T
}

func NewDatabaseMock(t *testing.T) DatabaseMock {
	return DatabaseMock{[][][]any{}, []error{}, []string{}, [][]any{}, t}
}

func (dbm *DatabaseMock) ExpectQuery(sql string, result [][]any, args []any, err error) {
	dbm.expectedSQLs = append(dbm.expectedSQLs, sql)
	dbm.expectedResults = append(dbm.expectedResults, result)
	dbm.expectedErrors = append(dbm.expectedErrors, err)
	dbm.expectedArgs = append(dbm.expectedArgs, args)
}

func (dbm *DatabaseMock) AssertAllExpectionsSatisfied() {
	if len(dbm.expectedSQLs) != 0 {
		fmt.Println("Not all expects satisfied!")
		dbm.t.FailNow()

	}
}

func (dbm *DatabaseMock) Query(sql string, args ...any) ([][]any, error) {
	if len(dbm.expectedSQLs) == 0 {
		fmt.Println("Unexpected call!")
		dbm.t.FailNow()
	}

	expectedSQL := dbm.expectedSQLs[len(dbm.expectedSQLs)-1]
	dbm.expectedSQLs = dbm.expectedSQLs[:len(dbm.expectedSQLs)-1]
	if sql != expectedSQL {
		fmt.Println("Unexpected sql!")
		dbm.t.FailNow()
	}

	expectedArgs := dbm.expectedArgs[len(dbm.expectedArgs)-1]
	dbm.expectedArgs = dbm.expectedArgs[:len(dbm.expectedArgs)-1]
	if !reflect.DeepEqual(args, expectedArgs) {
		fmt.Println("Unexpected args!")
		dbm.t.FailNow()

	}

	result := dbm.expectedResults[len(dbm.expectedResults)-1]
	dbm.expectedResults = dbm.expectedResults[:len(dbm.expectedResults)-1]
	err := dbm.expectedErrors[len(dbm.expectedErrors)-1]
	dbm.expectedErrors = dbm.expectedErrors[:len(dbm.expectedErrors)-1]

	return result, err
}
