package forge_test

import (
	"testing"

	"github.com/kyberbits/forge"
)

type PreparedStatementTestCase struct {
	Query              string
	Parameters         map[string]interface{}
	ExpectedQuery      string
	ExpectedParameters []interface{}
}

func TestPreparedStatementBasic(t *testing.T) {
	testPreparedStatement(t, PreparedStatementTestCase{
		Query: "SELECT * FROM users WHERE id = :id",
		Parameters: map[string]interface{}{
			"id": 1,
		},
		ExpectedQuery:      "SELECT * FROM users WHERE id = ?",
		ExpectedParameters: []interface{}{1},
	})
}

func TestPreparedStatementComplex(t *testing.T) {
	testPreparedStatement(t, PreparedStatementTestCase{
		Query: "SELECT * FROM users WHERE firstName = :first_name OR lastName = :first_name",
		Parameters: map[string]interface{}{
			"first_name": "Joe",
		},
		ExpectedQuery:      "SELECT * FROM users WHERE firstName = ? OR lastName = ?",
		ExpectedParameters: []interface{}{"Joe", "Joe"},
	})
}

func TestPreparedStatementQuestionDouble(t *testing.T) {
	testPreparedStatement(t, PreparedStatementTestCase{
		Query: "SELECT * FROM users WHERE firstName = \"hello ? \" OR lastName = :name",
		Parameters: map[string]interface{}{
			"name": "Joe",
		},
		ExpectedQuery:      "SELECT * FROM users WHERE firstName = \"hello ? \" OR lastName = ?",
		ExpectedParameters: []interface{}{"Joe"},
	})
}

func TestPreparedStatementQuestionSingle(t *testing.T) {
	testPreparedStatement(t, PreparedStatementTestCase{
		Query: "SELECT * FROM users WHERE firstName = 'hello ? ' OR lastName = :name",
		Parameters: map[string]interface{}{
			"name": "Joe",
		},
		ExpectedQuery:      "SELECT * FROM users WHERE firstName = 'hello ? ' OR lastName = ?",
		ExpectedParameters: []interface{}{"Joe"},
	})
}

func testPreparedStatement(t *testing.T, testCase PreparedStatementTestCase) {
	preparedStatement := forge.NewPreparedStatement(testCase.Query)
	for key, value := range testCase.Parameters {
		preparedStatement.SetValue(key, value)
	}

	actualQuery := preparedStatement.GetParsedQuery()
	actualParameters := preparedStatement.GetParsedParameters()

	if err := forge.Assert(testCase.ExpectedQuery, actualQuery); err != nil {
		t.Fatal(err)
	}
	if err := forge.Assert(testCase.ExpectedParameters, actualParameters); err != nil {
		t.Fatal(err)
	}
}
