package sqlexpr

import (
	"strconv"
)

type Setter struct {
	Field Expr
	Value interface{}
}

type Clause struct {
	Start     string
	Separator string
	Items     []Expr
}

func (v Clause) AppendToSQLBuilder(b *Builder) {
	if len(v.Items) == 0 {
		return
	}
	for i, item := range v.Items {
		if i == 0 {
			b.AppendRaw(v.Start)
		} else {
			b.AppendRaw(v.Separator)
		}
		b.Append(item)
	}
}

type Where []Expr

func (v Where) AppendToSQLBuilder(b *Builder) {
	Clause{"WHERE", "AND", v}.AppendToSQLBuilder(b)
}

type OrderBy []Expr

func (v OrderBy) AppendToSQLBuilder(b *Builder) {
	Clause{"ORDER BY", ",", v}.AppendToSQLBuilder(b)
}

type Returning []Expr

func (v Returning) AppendToSQLBuilder(b *Builder) {
	Clause{"RETURNING", ",", v}.AppendToSQLBuilder(b)
}

type LimitExpr int

func (v LimitExpr) AppendToSQLBuilder(b *Builder) {
	if v > 0 {
		b.AppendRaw("LIMIT")
		b.AppendRaw(strconv.Itoa(int(v)))
	}
}

func InnerJoin(a Expr, aCol Column, b Expr, bCol Column) Expr {
	return Fragment{a, Raw("INNER JOIN"), b, Raw("ON"), Dot(a, aCol), Raw("="), Dot(b, bCol)}
}

type Settable interface {
	Set(field Expr, value interface{})
}
type Whereable interface {
	AddWhere(conds ...Expr)
}
type Returnable interface {
	AddReturning(fields ...Expr)
}

type Select struct {
	Leading  Expr
	From     Expr
	Fields   List
	Where    Where
	Grouping Expr
	OrderBy  OrderBy
	Limit    int
	Trailing Expr
}

func (s *Select) AddField(fields ...Expr) {
	s.Fields = append(s.Fields, fields...)
}

func (s *Select) AddWhere(conds ...Expr) {
	s.Where = append(s.Where, conds...)
}

func (s Select) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw("SELECT")
	b.AppendExpr(s.Leading)
	b.AppendExpr(s.Fields)
	b.AppendRaw("FROM")
	b.AppendExpr(s.From)
	b.AppendExpr(s.Where)
	b.AppendExpr(s.Grouping)
	b.AppendExpr(s.OrderBy)
	b.AppendExpr(LimitExpr(s.Limit))
	b.AppendExpr(s.Trailing)
}

type Insert struct {
	Leading   Expr
	Table     Expr
	Setters   []Setter
	Trailing  Expr
	Returning Returning
}

func (s *Insert) Set(field Expr, value interface{}) {
	s.Setters = append(s.Setters, Setter{field, value})
}

func (s *Insert) AddReturning(fields ...Expr) {
	s.Returning = append(s.Returning, fields...)
}

func (s Insert) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw("INSERT INTO")
	b.AppendExpr(s.Table)
	b.AppendExpr(s.Leading)
	b.AppendRaw("(")
	for i, setter := range s.Setters {
		if i > 0 {
			b.AppendRaw(",")
		}
		b.AppendExpr(setter.Field)
	}
	b.AppendRaw(")")
	b.AppendRaw("VALUES")
	b.AppendRaw("(")
	for i, setter := range s.Setters {
		if i > 0 {
			b.AppendRaw(",")
		}
		b.Append(setter.Value)
	}
	b.AppendRaw(")")
	b.AppendExpr(s.Trailing)
	b.AppendExpr(s.Returning)
}

type Update struct {
	Table     Expr
	Leading   Expr
	Setters   []Setter
	Where     Where
	Trailing  Expr
	Returning Returning
}

func (s *Update) Set(field Expr, value interface{}) {
	s.Setters = append(s.Setters, Setter{field, value})
}

func (s *Update) AddWhere(conds ...Expr) {
	s.Where = append(s.Where, conds...)
}

func (s *Update) AddReturning(fields ...Expr) {
	s.Returning = append(s.Returning, fields...)
}

func (s Update) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw("UPDATE")
	b.AppendExpr(s.Table)
	b.AppendExpr(s.Leading)
	b.AppendRaw("SET")
	for i, setter := range s.Setters {
		if i > 0 {
			b.AppendRaw(",")
		}
		b.AppendExpr(setter.Field)
		b.AppendRaw("=")
		b.Append(setter.Value)
	}
	b.AppendExpr(s.Where)
	b.AppendExpr(s.Trailing)
	b.AppendExpr(s.Returning)
}

type Delete struct {
	Table     Expr
	Leading   Expr
	Where     Where
	Trailing  Expr
	Returning Returning
}

func (s *Delete) AddWhere(conds ...Expr) {
	s.Where = append(s.Where, conds...)
}

func (s *Delete) AddReturning(fields ...Expr) {
	s.Returning = append(s.Returning, fields...)
}

func (s Delete) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw("DELETE FROM")
	b.AppendExpr(s.Table)
	b.AppendExpr(s.Leading)
	b.AppendExpr(s.Where)
	b.AppendExpr(s.Trailing)
	b.AppendExpr(s.Returning)
}
