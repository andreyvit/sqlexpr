package sqlexpr

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{"empty", []interface{}{}, ""},
		{"single keyword", []interface{}{Raw("SELECT")}, "SELECT"},
		{"two keywords", []interface{}{Raw("SELECT"), Raw("DISTINCT")}, "SELECT DISTINCT"},
		{"", []interface{}{Raw("SELECT"), Raw("*"), Raw("FROM")}, "SELECT * FROM"},
		{"comma-separated names", []interface{}{Raw("SELECT"), Raw("foo"), Raw(","), Raw("bar"), Raw("FROM")}, "SELECT foo, bar FROM"},
		{"comma-separated quoted names", []interface{}{Raw("SELECT"), Raw(`"foo"`), Raw(","), Raw("bar"), Raw(","), Raw(`"boz"`), Raw("FROM")}, `SELECT "foo", bar, "boz" FROM`},
		{"dollar placeholders and commas", []interface{}{Raw("SELECT"), Raw("$1"), Raw(","), Raw("$2"), Raw("FROM")}, "SELECT $1, $2 FROM"},
		{"empty between keywords", []interface{}{Raw("SELECT"), Empty, Raw("DISTINCT")}, "SELECT DISTINCT"},

		{"single arg", []interface{}{10}, "$1 [10]"},
		{"two args", []interface{}{10, "foo"}, "$1 $2 [10, foo]"},

		{"typical SELECT", []interface{}{Raw("SELECT"), Column("foo"), Raw(","), Column("bar"), Raw("FROM"), Table("foos"), Raw("WHERE"), Column("boz"), Raw("="), 42}, "SELECT foo, bar FROM foos WHERE boz = $1 [42]"},
	}

	for _, test := range tests {
		if test.name == "" {
			test.name = test.expected
		}
		t.Run(test.name, func(t *testing.T) {
			b := new(Builder)
			b.AppendAll(test.input...)
			if a, e := b.String(), test.expected; a != e {
				t.Errorf("got %q, wanted %q", a, e)
			}
		})
	}
}
