# Idiomatic, simple and powerful SQL Builder for Go


Three goals:

1. Make constructing SQL queries less tedious (e.g. avoid matching up INSERT field names to values, construct IN expressions easily).

2. Allow to reuse SQL-related code (e.g. share code that sets fields across INSERTs and UPDATEs, share conditionals across SELECTs, UPDATEs and DELETEs).

3. Allow to use branches and loops when constructing SQL (e.g. add a WHERE condition inside an `if` statement)


Features:

1. Uses plain Go types (e.g. `Or` is a slice of expressions, `Select` is a struct with a bunch of fields) that are easy to manipulate without any layers of abstraction.

2. Allows to insert arbitrary SQL whenever needed (just pass `Raw`, which is just a `type Raw string`).

3. Allows to pick SQL dialect (this is a global setting, because working with multiple different database types in a single app is extremely rate and is not worth polluting the entire code with dialect selection).


## Usage

Import:

    "github.com/andreyvit/sqlexpr"

Define tables and columns:

```go
const (
    accounts   = sqlexpr.Table("accounts")
    id         = sqlexpr.Column("id")
    email      = sqlexpr.Column("email")
    name       = sqlexpr.Column("name")
    notes      = sqlexpr.Column("notes")
    updated_at = sqlexpr.Column("updated_at")
    deleted    = sqlexpr.Column("deleted")
)
```

### Select

```go
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
```

Result:

```sql
SELECT id, email, name, notes FROM accounts WHERE (name LIKE $1 OR notes LIKE $2) AND NOT deleted
```

### Insert

```go
s := sqlexpr.Insert{Table: accounts}
s.Set(email, "john@example.com")
s.Set(name, "John Doe")
s.Set(notes, "Little Johny")
s.Set(updated_at, sqlexpr.NOW)
s.AddReturning(id)

sql, args := sqlexpr.Build(s)
```

Result:

```sql
INSERT INTO accounts (email, name, notes, updated_at) VALUES ($1, $2, $3, NOW()) RETURNING id
```

### Update

```go
s := sqlexpr.Update{Table: accounts}
s.Set(email, "john@example.com")
s.Set(name, "John Doe")
s.Set(notes, "Little Johny")
s.Set(updated_at, sqlexpr.NOW)
s.AddWhere(sqlexpr.Eq(id, 42))
s.AddReturning(updated_at)

sql, args := sqlexpr.Build(s)
```

Result:

```sql
UPDATE accounts SET email = $1, name = $2, notes = $3, updated_at = NOW() WHERE id = $4 RETURNING updated_at
```

### Delete

```go
s := sqlexpr.Delete{Table: accounts}
s.AddWhere(sqlexpr.Eq(id, 42))

sql, args := sqlexpr.Build(s)
```

Result:

```sql
DELETE FROM accounts WHERE id = $1
```


## Principles

1. Everything that this package produces is an `sqlexpr.Expr`. You can turn an `Expr` into SQL (plus arguments slice) using `sqlexpr.Build(expr)`.

2. All names are passed as `sqlexpr.Table` and `sqlexpr.Column`, all SQL keywords, punctuation and raw SQL code as `sqlexpr.Raw`. These are all just strings, defined as new types.

3. Any value that is not `Expr` becomes an argument (i.e. adds a placeholders like `$1` or `?` into the SQL statement).

4. There are four top-level types of Exprs: `sqlexpr.Select`, `sqlexpr.Insert`, `sqlexpr.Update` and `sqlexpr.Delete`. These are structs with very simple fields that you can fill.

5. These four structs define a few simple helper methods, e.g. `AddWhere` appends the given condition to `.Where` slice. The helpers both simplify your code and allow code reuse via interfaces:

    ```go
    type Settable interface {
        Set(field Expr, value interface{})
    }
    type Whereable interface {
        AddWhere(conds ...Expr)
    }
    type Returnable interface {
        AddReturning(fields ...Expr)
    }
    ```

6. All other expressions can be produced via types and functions in this package. See the tests for examples.


## Development

Use [modd](https://github.com/cortesi/modd) (`go get github.com/cortesi/modd/cmd/modd`) to re-run tests automatically during development by running `modd` (recommended).


## License

MIT.
