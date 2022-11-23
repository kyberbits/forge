package forge

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

type PreparedStatement struct {
	positions     map[string][]int
	parameters    []interface{}
	originalQuery string
	revisedQuery  string
}

func NewPreparedStatement(queryText string) *PreparedStatement {
	preparedStatement := &PreparedStatement{}
	preparedStatement.positions = make(map[string][]int, 8)
	preparedStatement.setQuery(queryText)

	return preparedStatement
}

func (preparedStatement *PreparedStatement) setQuery(queryText string) {
	var revisedBuilder bytes.Buffer
	var parameterBuilder bytes.Buffer
	var position []int
	var character rune
	var parameterName string
	var width int
	var positionIndex int

	preparedStatement.originalQuery = queryText
	positionIndex = 0

	for i := 0; i < len(queryText); {
		character, width = utf8.DecodeRuneInString(queryText[i:])
		i += width

		if character == ':' {

			for {

				character, width = utf8.DecodeRuneInString(queryText[i:])
				i += width

				if unicode.IsLetter(character) || unicode.IsDigit(character) || character == '_' {
					parameterBuilder.WriteString(string(character))
				} else {
					break
				}
			}

			parameterName = parameterBuilder.String()
			position = preparedStatement.positions[parameterName]
			preparedStatement.positions[parameterName] = append(position, positionIndex)
			positionIndex++

			revisedBuilder.WriteString("?")
			parameterBuilder.Reset()

			if width <= 0 {
				break
			}
		}

		revisedBuilder.WriteString(string(character))

		// Don't touch ?'s in ''
		if character == '\'' {
			for {
				character, width = utf8.DecodeRuneInString(queryText[i:])
				i += width
				revisedBuilder.WriteString(string(character))

				if character == '\'' {
					break
				}
			}
		}

		// Don't touch ?'s in ""
		if character == '"' {
			for {
				character, width = utf8.DecodeRuneInString(queryText[i:])
				i += width
				revisedBuilder.WriteString(string(character))

				if character == '"' {
					break
				}
			}
		}
	}

	preparedStatement.revisedQuery = revisedBuilder.String()
	preparedStatement.parameters = make([]interface{}, positionIndex)
}

func (preparedStatement *PreparedStatement) GetParsedQuery() string {
	return preparedStatement.revisedQuery
}

func (preparedStatement *PreparedStatement) GetParsedParameters() []interface{} {
	return preparedStatement.parameters
}

func (preparedStatement *PreparedStatement) SetValue(parameterName string, parameterValue interface{}) {
	for _, position := range preparedStatement.positions[parameterName] {
		preparedStatement.parameters[position] = parameterValue
	}
}
