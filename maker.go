package sql_query_maker

import (
	"fmt"
	"strings"
)

var prefixes = [...]string{"or", "OR", "and", "AND"}

type insertParams struct {
	isValuesAdded bool
}

type SqlQueryMaker struct {
	symbol       rune
	query        *strings.Builder
	args         []interface{}
	fieldsCount  int
	insertParams insertParams
}

func NewQueryMaker(argsCount int) *SqlQueryMaker {
	q := &SqlQueryMaker{
		query:       &strings.Builder{},
		args:        make([]interface{}, 0, argsCount),
		fieldsCount: 1,
		symbol:      '?',
	}
	q.query.Grow(1000)

	return q
}

func (q *SqlQueryMaker) ChangeSymbol(symbol rune) {
	q.symbol = symbol
}

func (q *SqlQueryMaker) Add(query string, args ...interface{}) *SqlQueryMaker {
	defer q.query.WriteRune(' ')

	if len(args) == 0 {
		q.query.WriteString(query)
		return q
	}

	q.args = append(q.args, args...)
	runes := []rune(query)

	for i := 0; i < len(runes); i++ {
		if runes[i] == q.symbol {
			q.query.WriteString(fmt.Sprintf("$%d", q.fieldsCount))
			q.fieldsCount++
		} else {
			q.query.WriteRune(runes[i])
		}
	}

	return q
}

func (q *SqlQueryMaker) WhereOptional(modifyFunc func()) *SqlQueryMaker {
	mainBuilder := q.query

	q.query = &strings.Builder{}
	q.query.Grow(100)

	startLen := q.query.Len()

	modifyFunc()

	if q.query.Len() != startLen {
		mainBuilder.WriteString("WHERE ")
		queryStr := strings.TrimSpace(q.query.String())

		for i := 0; i < 4; i++ {
			queryStr = strings.TrimPrefix(queryStr, prefixes[i])
		}

		mainBuilder.WriteString(queryStr)
		q.query = mainBuilder
		q.query.WriteRune(' ')
	} else {
		q.query = mainBuilder
	}

	return q
}

func (q *SqlQueryMaker) Where(query string, args ...interface{}) *SqlQueryMaker {
	curQuery := strings.TrimSuffix(strings.TrimSpace(q.query.String()), ",")

	q.query.Reset()

	_, _ = q.query.WriteString(curQuery)

	q.query.WriteRune(' ')

	q.Add("WHERE")

	return q.Add(query, args...)
}

func (q *SqlQueryMaker) Clear() *SqlQueryMaker {
	q.query.Reset()
	q.fieldsCount = 1
	q.args = q.args[:0]
	q.insertParams.isValuesAdded = false

	return q
}

func (q *SqlQueryMaker) AND() *SqlQueryMaker {
	q.query.WriteString("AND ")
	return q
}

func (q *SqlQueryMaker) OR() *SqlQueryMaker {
	q.query.WriteString("OR ")
	return q
}

func (q *SqlQueryMaker) Query() string {
	return q.query.String()
}

func (q *SqlQueryMaker) Args() []interface{} {
	return q.args
}

// Make return query and args
func (q *SqlQueryMaker) Make() (string, []interface{}) {
	return q.query.String(), q.args
}

func (q *SqlQueryMaker) Values(args ...interface{}) *SqlQueryMaker {
	strBuilder := strings.Builder{}

	if q.insertParams.isValuesAdded {
		strBuilder.WriteRune(',')
	} else {
		strBuilder.WriteString("VALUES ")
	}

	strBuilder.WriteRune('(')
	for i := 0; i < len(args); i++ {
		if i != 0 {
			strBuilder.WriteRune(',')
		}

		strBuilder.WriteRune('?')
	}
	strBuilder.WriteRune(')')

	if !q.insertParams.isValuesAdded {
		q.insertParams.isValuesAdded = true
	}

	q.Add(strBuilder.String(), args...)

	return q
}
