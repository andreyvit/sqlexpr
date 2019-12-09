package sqlexpr

import (
	"testing"
)

func TestPostgresDialect(t *testing.T) {
	d := PostgresDialect
	t.Run("FormatPlaceholder", func(t *testing.T) {
		a := d.FormatPlaceholder(0) + " " + d.FormatPlaceholder(1)
		e := "$1 $2"
		if a != e {
			t.Errorf("got %q, wanted %q", a, e)
		}
	})
}

func TestSQLiteDialect(t *testing.T) {
	d := SQLiteDialect
	t.Run("FormatPlaceholder", func(t *testing.T) {
		a := d.FormatPlaceholder(0) + " " + d.FormatPlaceholder(1)
		e := "? ?"
		if a != e {
			t.Errorf("got %q, wanted %q", a, e)
		}
	})
}

func TestMySQLDialect(t *testing.T) {
	d := MySQLDialect
	t.Run("FormatPlaceholder", func(t *testing.T) {
		a := d.FormatPlaceholder(0) + " " + d.FormatPlaceholder(1)
		e := "? ?"
		if a != e {
			t.Errorf("got %q, wanted %q", a, e)
		}
	})
}
