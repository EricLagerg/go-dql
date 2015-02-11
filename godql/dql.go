package godql

import (
	"bytes"
	"strconv"
	"strings"
)

const strOr = "OR"

var (
	selec         = []byte("SELECT ")
	selecDistinct = []byte("SELECT DISTINCT ")
	and           = []byte(" AND ")
	or            = []byte(" OR ")
	wh            = []byte(" WHERE ")
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
		buf.Write(wh)

		l := len(q.w)
		for i := 0; i < l; i++ {
			if q.w[i] != strOr {
				buf.WriteString(q.w[i])

				if (i%2 == 0 || i == 1) &&
					(i < l-1 && q.w[i+1] != strOr) {
					buf.Write(and)
				}
			} else {
				buf.Write(or)
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
