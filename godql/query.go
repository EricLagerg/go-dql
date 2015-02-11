package godql

import (
	"fmt"
)

// Comparison operators and aliases for both versions of the Where() functions.
const (
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

// Formatting strings
const (
	intFormat  = "%s%s%d"
	strFormat  = "%s%s'%s'"
	boolFormat = "%s%s%b"
)

// Generic structure of a Query's parts
type Expr struct {
	S  []string // SELECT
	SD string   // SELECT DISTINCT
	SC string   // SELECT COUNT
	W  []string // WHERE
	G  string   // GROUP BY
	L  int      // LIMIT
}

// Initialize new query. Calls NewQuery under the covers, but looks
// prettier...
func Query() *Expr {
	return NewQuery()
}

// Proper Go form for initializing a new query
func NewQuery() *Expr {
	return new(Expr)
}

// Generic SELECT statement
//
// Select([]string{"firstname", "lastname", "dateofbirth"}) is
// equivalent to: SELECT firstname,lastname,dateofbirth
func (e *Expr) Select(fields []string) *Expr {
	e.S = fields
	return e
}

// SELECT DISTINCT
//
// SelectDistinct("firstname") is equivalent to:
// SELECT DISTINCT firstname
func (e *Expr) SelectDistinct(field string) *Expr {
	e.SD = field
	return e
}

// SelectCount("firstname") is equivalent to:
// SELECT COUNT(firstname)
func (e *Expr) SelectCount(field string) *Expr {
	e.SC = fmt.Sprintf("SELECT COUNT(%s)", field)
	return e
}

// CountDistinct("firstname") is equivalent to:
func (e *Expr) CountDistinct(field string) *Expr {
	e.SC = fmt.Sprintf("SELECT COUNT(DISTINCT %s)", field)
	return e
}

// enum for position of Where() inside of Group()
const (
	beg = iota
	mid
	end
	two
)

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
//
// Used to group two WHERE expressions together, e.g.:
// ... WHERE (firstname='eric' AND lastname='lagergren') OR ...
func (e *Expr) Group(exprs ...*where) *Expr {
	l := len(exprs)
	if l == 0 {
		panic("cannot Group() without Where()!")
	}

	// append first element regardless
	if l > 2 {
		e.W = append(e.W, exprs[0].whereFormat(beg))
	} else {
		e.W = append(e.W, exprs[0].whereFormat(two))
		return e // notice hidden return if l == 1
	}

	// loop over middle elements
	for i := 1; i < l-1; i++ {
		if exprs[i].std {
			e.W = append(e.W, exprs[i].whereFormat(mid))
		} else {
			// "OR"
			e.W = append(e.W, strOr)
		}
	}

	// append last regardless
	e.W = append(e.W, exprs[l-1].whereFormat(end))
	return e
}

// Where("firstname", Equals, "Eric")
//
// Used inside a Group() since it's not a method of the Expr struct
// Returns the values in a struct which is passed to whereFormat()
// which formats
func Where(field string, op string, value interface{}) *where {
	return &where{field, op, value, true}
}

// Where("firstname", Equals, "Eric") ...
// Where("age", Greater, 20) ...
// Where("active", NotEqual, true) ...
//
// Between each Where() statement an "AND" will be inserted unless
// an Or() is used.
func (e *Expr) Where(field string, op string, value interface{}) *Expr {
	e.W = append(e.W, typeFormat(field, op, value))
	return e
}

// Used outside of a Group() to denote an OR statement between two
// WHERE statements
//
// e.g. ... Where("firstname" Equals, "Eric").Or().Where("lastname", Equals, "lagergren") ...
func (e *Expr) Or() *Expr {
	e.W = append(e.W, strOr)
	return e
}

// Used inside a Group() to denote an OR statement between two
// WHERE statements
//
// e.g. Group(Where(), Or(), Where())
func Or() *where {
	return &where{"", "", nil, false}
}

// GROUP BY statement
//
// GroupBy("party") is equivalent to: GROUP BY party
func (e *Expr) GroupBy(field string) *Expr {
	e.G = field
	return e
}

// LIMIT statement
//
// Limit(200) is equivalent to: LIMIT 200
func (e *Expr) Limit(num int) *Expr {
	e.L = num
	return e
}

// Generates all DQL parameter functions and returns formatted DQL string
func (e *Expr) String() string {
	return e.toDql()
}
