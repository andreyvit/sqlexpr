package sqlexpr_test

import (
	"fmt"
	"github.com/andreyvit/sqlexpr"
)

func ExampleSelect() {
	const (
		accounts = sqlexpr.Table("accounts")
		id       = sqlexpr.Column("id")
		email    = sqlexpr.Column("email")
		name     = sqlexpr.Column("name")
		notes    = sqlexpr.Column("notes")
		deleted  = sqlexpr.Column("deleted")
	)

	includeNotes := true
	query := "%abc%"

	s := sqlexpr.Select{From: accounts}
	s.AddField(id, email, name)
	if includeNotes {
		s.AddField(notes)
	}
	s.AddWhere(sqlexpr.Or{sqlexpr.Like(name, query), sqlexpr.Like(notes, query)})
	s.AddWhere(sqlexpr.Not(deleted))

	sql, args := sqlexpr.Build(s)
	fmt.Println(sql)
	fmt.Printf("%#v", args)

	// Output: SELECT id, email, name, notes FROM accounts WHERE (name LIKE $1 OR notes LIKE $2) AND NOT deleted
	// []interface {}{"%abc%", "%abc%"}
}

func ExampleInsert() {
	const (
		accounts   = sqlexpr.Table("accounts")
		id         = sqlexpr.Column("id")
		email      = sqlexpr.Column("email")
		name       = sqlexpr.Column("name")
		notes      = sqlexpr.Column("notes")
		updated_at = sqlexpr.Column("updated_at")
		deleted    = sqlexpr.Column("deleted")
	)

	s := sqlexpr.Insert{Table: accounts}
	s.Set(email, "john@example.com")
	s.Set(name, "John Doe")
	s.Set(notes, "Little Johny")
	s.Set(updated_at, sqlexpr.NOW)
	s.AddReturning(id)

	sql, args := sqlexpr.Build(s)
	fmt.Println(sql)
	fmt.Printf("%#v", args)

	// Output: INSERT INTO accounts (email, name, notes, updated_at) VALUES ($1, $2, $3, NOW()) RETURNING id
	// []interface {}{"john@example.com", "John Doe", "Little Johny"}
}

func ExampleUpdate() {
	const (
		accounts   = sqlexpr.Table("accounts")
		id         = sqlexpr.Column("id")
		email      = sqlexpr.Column("email")
		name       = sqlexpr.Column("name")
		notes      = sqlexpr.Column("notes")
		updated_at = sqlexpr.Column("updated_at")
		deleted    = sqlexpr.Column("deleted")
	)

	s := sqlexpr.Update{Table: accounts}
	s.Set(email, "john@example.com")
	s.Set(name, "John Doe")
	s.Set(notes, "Little Johny")
	s.Set(updated_at, sqlexpr.NOW)
	s.AddWhere(sqlexpr.Eq(id, 42))
	s.AddReturning(updated_at)

	sql, args := sqlexpr.Build(s)
	fmt.Println(sql)
	fmt.Printf("%#v", args)

	// Output: UPDATE accounts SET email = $1, name = $2, notes = $3, updated_at = NOW() WHERE id = $4 RETURNING updated_at
	// []interface {}{"john@example.com", "John Doe", "Little Johny", 42}
}

func ExampleDelete() {
	const (
		accounts = sqlexpr.Table("accounts")
		id       = sqlexpr.Column("id")
		email    = sqlexpr.Column("email")
		name     = sqlexpr.Column("name")
		notes    = sqlexpr.Column("notes")
		deleted  = sqlexpr.Column("deleted")
	)

	s := sqlexpr.Delete{Table: accounts}
	s.AddWhere(sqlexpr.Eq(id, 42))

	sql, args := sqlexpr.Build(s)
	fmt.Println(sql)
	fmt.Printf("%#v", args)

	// Output: DELETE FROM accounts WHERE id = $1
	// []interface {}{42}
}
