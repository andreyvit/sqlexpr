package sqlexpr

import (
	"testing"
)

func TestStatements(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{"where 0", Where{}, ""},
		{"where 1", Where{Eq(Column("foo"), 42)}, "WHERE foo = $1 [42]"},
		{"where 2", Where{Eq(Column("foo"), 42), Column("is_bar")}, "WHERE foo = $1 AND is_bar [42]"},

		{"simple select", Select{
			Fields: List{Star},
			From:   Table("foos"),
		}, "SELECT * FROM foos"},

		{"complex select", Select{
			Leading: Raw("DISTINCT"),
			Fields:  List{Column("foo"), Column("bar"), Column("boz")},
			From:    InnerJoin(Table("foos"), "widget_id", Table("widgets"), "id"),
			Where:   Where{Eq(Column("foo"), 42), IsNotNull(Column("bar"))},
			OrderBy: OrderBy{Column("foo"), Desc(Column("bar"))},
			Limit:   1,
		}, "SELECT DISTINCT foo, bar, boz FROM foos INNER JOIN widgets ON foos.widget_id = widgets.id WHERE foo = $1 AND bar IS NOT NULL ORDER BY foo, bar DESC LIMIT 1 [42]"},

		{"simple insert", Insert{
			Table: Table("foos"),
			Setters: []Setter{
				Setter{Column("foo"), 42},
				Setter{Column("bar"), "test"},
			},
		}, "INSERT INTO foos (foo, bar) VALUES ($1, $2) [42, test]"},

		{"complex insert", Insert{
			Table: Table("foos"),
			Setters: []Setter{
				Setter{Column("foo"), 42},
				Setter{Column("bar"), "test"},
			},
			Returning: Returning{Column("boz")},
		}, "INSERT INTO foos (foo, bar) VALUES ($1, $2) RETURNING boz [42, test]"},

		{"simple update", Update{
			Table: Table("foos"),
			Setters: []Setter{
				Setter{Column("foo"), 42},
				Setter{Column("bar"), "test"},
			},
		}, "UPDATE foos SET foo = $1, bar = $2 [42, test]"},

		{"complex update", Update{
			Table: Table("foos"),
			Setters: []Setter{
				Setter{Column("foo"), 42},
				Setter{Column("bar"), "test"},
			},
			Where:     Where{Op(Column("boz"), "<", 123)},
			Returning: Returning{Column("boz")},
		}, "UPDATE foos SET foo = $1, bar = $2 WHERE boz < $3 RETURNING boz [42, test, 123]"},

		{"complex delete", Delete{
			Table:     Table("foos"),
			Where:     Where{Op(Column("boz"), "=", 123)},
			Returning: Returning{Column("foo"), Column("bar")},
		}, "DELETE FROM foos WHERE boz = $1 RETURNING foo, bar [123]"},
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
