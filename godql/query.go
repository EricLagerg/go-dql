package godql

import (
	"fmt"
)

// Operators and aliases
const (
	// Comparison operators
	Noop           = ""
	Equals         = "="
	NotEqual       = "!="
	DoesntEqual    = NotEqual
	GreaterThan    = ">"
	Greater        = GreaterThan
	GreaterOrEqual = ">="
	LessThan       = "<"
	Less           = LessThan
	LessOrEqual    = "<="
	Like           = "~"
	NotLike        = "!~"
	Unlike         = NotLike
)

// S: Select, SD: SelectDistinct, SC: SelectCount, W: Where,
// G: GroupBy, L: Limit
type Query struct {
	s  []string
	sd string
	sc string
	w  []string
	g  string
	l  int
}

// Select([]string{"firstname", "lastname", "dateofbirth"})
// SELECT firstname,lastname,dateofbirth
func (q *Query) Select(fields []string) *Query {
	q.s = fields
	return q
}

// SelectDistinct("firstname")
// SELECT DISTINCT firstname
func (q *Query) SelectDistinct(field string) *Query {
	q.sd = field
	return q
}

// SelectCount("firstname")
// SELECT COUNT(firstname)
func (q *Query) SelectCount(field string) *Query {
	q.sc = fmt.Sprintf("SELECT COUNT(%s)", field)
	return q
}

// CountDistinct("firstname")
func (q *Query) CountDistinct(field string) *Query {
	q.sc = fmt.Sprintf("SELECT COUNT(DISTINCT %s)", field)
	return q
}

// Where("firstname", Equals, "Eric")
// Where("age", Greater, 20)
// Where("active", NotEqual, true)
// Between each Where() statement an "AND" will be inserted unless
// an ungroup-ed OR is used --
func (q *Query) Where(field string, op string, value interface{}) *Query {
	switch value.(type) {
	case int:
		q.w = append(q.w, fmt.Sprintf("%s%s%d", field, op, value))
	case string:
		q.w = append(q.w, fmt.Sprintf("%s%s'%s'", field, op, value))
	case bool:
		q.w = append(q.w, fmt.Sprintf("%s%s%b", field, op, value))
	default:
		panic("Invalid type")
	}
	return q
}

// Where("firstname" Equals, "Eric").And().Where("lastname", Equals, "lagergren")
// Used to group two or more AND statements
func (q *Query) And() *Query {
	q.w = append(q.w, "+")
	return q
}

// Where("firstname" Equals, "Eric").Or().Where("lastname", Equals, "lagergren")
// Used to group two or more OR statements
func (q *Query) Or() *Query {
	q.w = append(q.w, "|")
	return q
}

func (q *Query) Org() *Query {
	q.w = append(q.w, "OR")
	return q
}

// GroupBy("party")
func (q *Query) GroupBy(field string) *Query {
	q.g = field
	return q
}

// Limit(200)
func (q *Query) Limit(num int) *Query {
	q.l = num
	return q
}

// Generates all DQL parameter functions and returns formatted DQL string
func (q *Query) String() string {
	return q.toDql()
}
