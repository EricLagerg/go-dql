package godql

import (
	"bytes"
	"strconv"
	"strings"
)

// Same as []byte(" ")
const (
	strOR      = "OR"
	space      = 32 // " "
	leftParen  = 40 // (
	rightParen = 41 // )
)

var (
	groupAnd = "+"
	groupOr  = "|"

	selec         = []byte("SELECT ")
	selecDistinct = []byte("SELECT DISTINCT ")
	and           = []byte(" AND ")
	or            = []byte(" OR ")
	where         = []byte(" WHERE ")
	limit         = []byte(" LIMIT ")
	groupBy       = []byte(" GROUP BY ")
)

func (q *Query) toDql() string {
	var buf bytes.Buffer

	// SELECT
	if q.sd != "" {
		buf.Write(selecDistinct)
		buf.WriteString(q.sd)
	} else if q.s != nil {
		buf.Write(selec)
		buf.WriteString(strings.Join(q.s, ","))
	} else if q.sc != "" {
		buf.WriteString(q.sc)
	} else {
		panic("No select statement")
	}

	// WHERE
	if q.w != nil {
		buf.Write(where)

		var tb bytes.Buffer
		var pc, wc int
		l := len(q.w)
		for i := 0; i < l; i++ {
			if q.w[i] == groupAnd || q.w[i] == groupOr {
				buf.WriteByte(leftParen)
				pc++
			} else {
				tb.WriteString(q.w[i])
				if q.w[i] == strOR {
					tb.Write(or)
				} else if i != l-1 {
					tb.Write(and)
				}
			}

			if pc == wc {
				buf.Write(tb.Bytes())
				for i := 0; i < pc; i++ {
					buf.WriteByte(rightParen)
				}
				pc = 0
				wc = 0
			}
		}
	}

	// GROUP BY
	if q.g != "" {
		buf.Write(groupBy)
		buf.WriteString(q.g)
	}

	// LIMIT
	if q.l > 0 {
		buf.Write(limit)
		buf.Write([]byte(strconv.Itoa(q.l)))
	}

	return buf.String()
}
