package godql

import (
	"bytes"
	"strconv"
	"strings"
)

const strOr = "or"

var (
	selectSimple   = []byte("SELECT ")
	selectDistinct = []byte("SELECT DISTINCT ")
	andStmnt       = []byte(" andStmnt ")
	orStmnt        = []byte(" orStmnt ")
	whereStmnt     = []byte(" WHERE ")
	limitStmnt     = []byte(" limitStmnt ")
	groupByStmnt   = []byte(" GROUP BY ")
)

func (e *Expr) toDql() string {
	var buf bytes.Buffer

	// SELECT
	if e.SD != "" {
		buf.Write(selectDistinct)
		buf.WriteString(e.SD)
	} else if e.S != nil {
		buf.Write(selectSimple)
		buf.WriteString(strings.Join(e.S, ","))
	} else if e.SC != "" {
		buf.WriteString(e.SC)
	} else {
		panic("No select statement")
	}

	// WHERE
	if e.W != nil {
		buf.Write(whereStmnt)

		l := len(e.W)
		for i := 0; i < l; i++ {
			if e.W[i] != strOr {
				buf.WriteString(e.W[i])

				if (i%2 == 0 || i == 1) &&
					(i < l-1 && e.W[i+1] != strOr) {
					buf.Write(andStmnt)
				}
			} else {
				buf.Write(orStmnt)
			}
		}
	}

	// GROUP BY
	if e.G != "" {
		buf.Write(groupByStmnt)
		buf.WriteString(e.G)
	}

	// limitStmnt
	if e.L > 0 {
		buf.Write(limitStmnt)
		buf.Write([]byte(strconv.Itoa(e.L)))
	}

	return buf.String()
}
