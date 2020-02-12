package sqlexpr

import (
	"testing"
)

func TestExprs(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{"column eq", Eq(Column("foo"), Column("bar")), "foo = bar"},
		{"column value eq", Eq(Column("foo"), 42), "foo = $1 [42]"},
		{"less than op", Op(Column("foo"), "<", 42), "foo < $1 [42]"},

		{"IsNull", IsNull(Column("foo")), "foo IS NULL"},
		{"IsNotNull", IsNotNull(Column("foo")), "foo IS NOT NULL"},

		{"List", List{Column("foo"), Value(42), Column("bar")}, "foo, $1, bar [42]"},

		{"In empty array", In(Column("foo"), Array{}), "FALSE"},
		{"In single-element array", In(Column("foo"), Array{42}), "foo = $1 [42]"},
		{"In array", In(Column("foo"), Array{10, 20, 30}), "foo IN ($1, $2, $3) [10, 20, 30]"},

		{"ArrayOfInt64s", ArrayOfInt64s([]int64{10, 20, 30}), "($1, $2, $3) [10, 20, 30]"},
		{"ArrayOfInts", ArrayOfInts([]int{10, 20, 30}), "($1, $2, $3) [10, 20, 30]"},
		{"ArrayOfStrings", ArrayOfStrings([]string{"foo", "bar", "boz"}), "($1, $2, $3) [foo, bar, boz]"},

		{"And", And{Column("foo"), Column("bar")}, "(foo AND bar)"},
		{"And 1", And{Column("foo")}, "foo"},
		{"And 0", And{}, "TRUE"},

		{"Or", Or{Column("foo"), Column("bar")}, "(foo OR bar)"},
		{"Or 1", Or{Column("foo")}, "foo"},
		{"Or 0", Or{}, "FALSE"},

		{"Not", Not(Column("foo")), "NOT foo"},

		{"Dot", Dot(Table("foo"), Column("bar")), "foo.bar"},

		{"Func0", Func("foo"), "foo ()"},
		{"Func1", Func("foo", Column("boz")), "foo (boz)"},
		{"Func2", Func("foo", "bar", Column("boz")), "foo ($1, boz) [bar]"},
		{"Max", Max(Column("boz")), "MAX (boz)"},
		{"Min", Min(Column("boz")), "MIN (boz)"},
		{"Count", Count(Column("boz")), "COUNT (boz)"},

		{"As", As(Column("foo"), "bar"), "foo AS bar"},

		{"Parens(func)", Parens(NOW), "(NOW())"},
		{"Parens(value)", Parens(Value(123)), "($1) [123]"},
		{"Parens(subquery)", Parens(Select{From: Table("foo"), Fields: List{Column("x")}}), "(SELECT x FROM foo)"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sql, args := Build(test.expr)
			a := FormatSQLArgs(sql, args)
			if a != test.expected {
				t.Errorf("got %q, wanted %q", a, test.expected)
			}
		})
	}
}
