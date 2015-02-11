package godql

import (
	"fmt"
)

// Operators and aliases
const (
	// Comparison operators
	_              = ""
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

// enum for position of Where() inside of Group()
const (
	beg = iota
	mid
	end
	two
)

// Formatting
const (
	intFormat  = "%s%s%d"
	strFormat  = "%s%s'%s'"
	boolFormat = "%s%s%b"
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

// unexported struct used to package Where()
type where struct {
	field string
	op    string
	value interface{}
	std   bool
}

// to quote or not to quote the value
func typeFormat(field string, op string, value interface{}) string {
	switch value.(type) {
	case int:
		return fmt.Sprintf("%s%s%d", field, op, value)
	case string:
		return fmt.Sprintf("%s%s'%s'", field, op, value)
	case bool:
		return fmt.Sprintf("%s%s%b", field, op, value)
	default:
		panic("Invalid type")
	}
}

// unexported for formatting Where()s inside a Group()
func (w *where) whereFormat(pos int) string {
	v := typeFormat(w.field, w.op, w.value)
	switch pos {
	case mid: // most common
		return v
	case beg:
		return fmt.Sprintf("(%s", v)
	case end:
		return fmt.Sprintf("%s)", v)
	case two:
		return fmt.Sprintf("(%s)", v)
	default:
		panic("shouldn't be here")
	}
}

// Group(Where("firstname", Equals, "Eric"), Where("lastname", Equals, "lagergren"))
// Used to group two WHERE expressions together, e.g.:
// ... WHERE (firstname='eric' AND lastname='lagergren') OR ...
func (q *Query) Group(exprs ...*where) *Query {
	l := len(exprs)
	if l == 0 {
		panic("cannot Group() without Where()!")
	}

	// append first element regardless
	if l > 2 {
		q.w = append(q.w, exprs[0].whereFormat(beg))
	} else {
		q.w = append(q.w, exprs[0].whereFormat(two))
		return q // notice hidden return if l == 1
	}

	// loop over middle elements
	for i := 1; i < l-1; i++ {
		if exprs[i].std {
			q.w = append(q.w, exprs[i].whereFormat(mid))
		} else {
			// "OR"
			q.w = append(q.w, strOr)
		}
	}

	// append last regardless
	q.w = append(q.w, exprs[l-1].whereFormat(end))
	return q
}

// Where("firstname", Equals, "Eric")
// Used inside a Group() since it's not a method of the Query struct
// Returns the values in a struct which is passed to whereFormat()
// which formats
func Where(field string, op string, value interface{}) *where {
	return &where{field, op, value, true}
}

// Where("firstname", Equals, "Eric")
// Where("age", Greater, 20)
// Where("active", NotEqual, true)
// Between each Where() statement an "AND" will be inserted unless
// an ungroup-ed OR is used --
func (q *Query) Where(field string, op string, value interface{}) *Query {
	q.w = append(q.w, typeFormat(field, op, value))
	return q
}

// Where("firstname" Equals, "Eric").Or().Where("lastname", Equals, "lagergren")
func (q *Query) Or() *Query {
	q.w = append(q.w, strOr)
	return q
}

// Bastardized Or() for use on Group(Where(), Or(), Where())
func Or() *where {
	return &where{"", "", nil, false}
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
