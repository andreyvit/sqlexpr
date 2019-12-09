package sqlexpr

import (
	"strconv"
)

type ArgStyle int

const (
	QuestionMarkArgs ArgStyle = iota
	DollarNumberArgs
)

type Dialect struct {
	Name     string
	ArgStyle ArgStyle
}

func (d *Dialect) FormatPlaceholder(index int) string {
	switch d.ArgStyle {
	case QuestionMarkArgs:
		return "?"
	case DollarNumberArgs:
		return "$" + strconv.Itoa(index+1)
	default:
		panic("unknown ArgStyle")
	}
}

var PostgresDialect = &Dialect{
	Name:     "PostgreSQL",
	ArgStyle: DollarNumberArgs,
}

var SQLiteDialect = &Dialect{
	Name:     "SQLite",
	ArgStyle: QuestionMarkArgs,
}

var MySQLDialect = &Dialect{
	Name:     "MySQL",
	ArgStyle: QuestionMarkArgs,
}

var dialect *Dialect = PostgresDialect

func SetDialect(d *Dialect) {
	dialect = d
}
