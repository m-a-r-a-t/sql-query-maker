package sql_query_maker

import (
	"fmt"
	"strings"
)

type SqlQueryMaker struct {
	symbol      rune
	query       strings.Builder
	args        []interface{}
	fieldsCount int
}

func NewQueryMaker(argsCount int) *SqlQueryMaker {
	q := &SqlQueryMaker{
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

func (q *SqlQueryMaker) Clear() *SqlQueryMaker {
	q.query.Reset()
	q.fieldsCount = 1
	q.args = q.args[:0]
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
