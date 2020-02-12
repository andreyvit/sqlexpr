package sqlexpr

type value struct {
	v interface{}
}

func (v value) AppendToSQLBuilder(b *Builder) {
	b.Append(v.v)
}

func Value(v interface{}) Expr {
	if e, ok := v.(Expr); ok {
		return e
	} else {
		return value{v}
	}
}

type Fragment []interface{}

func (v Fragment) AppendToSQLBuilder(b *Builder) {
	b.AppendAll(v...)
}

func Op(lhs interface{}, op string, rhs interface{}) Expr {
	return Fragment{lhs, Raw(op), rhs}
}

func Eq(lhs interface{}, rhs interface{}) Expr {
	return Op(lhs, "=", rhs)
}

func Like(lhs interface{}, rhs interface{}) Expr {
	return Op(lhs, "LIKE", rhs)
}

func NotLike(lhs interface{}, rhs interface{}) Expr {
	return Op(lhs, "NOT LIKE", rhs)
}

func LikeCaseInsensitive(lhs interface{}, rhs interface{}) Expr {
	return Op(lhs, "ILIKE", rhs)
}

func NotLikeCaseInsensitive(lhs interface{}, rhs interface{}) Expr {
	return Op(lhs, "NOT ILIKE", rhs)
}

func IsNull(v interface{}) Expr {
	return Fragment{v, Raw("IS NULL")}
}

func IsNotNull(v interface{}) Expr {
	return Fragment{v, Raw("IS NOT NULL")}
}

type List []Expr

func (v List) AppendToSQLBuilder(b *Builder) {
	for i, item := range v {
		if i > 0 {
			b.AppendRaw(",")
		}
		b.AppendExpr(item)
	}
}

type Array []interface{}

func (v Array) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw("(")
	for i, item := range v {
		if i > 0 {
			b.AppendRaw(",")
		}
		b.Append(item)
	}
	b.AppendRaw(")")
}

func ArrayOfInt64s(items []int64) Array {
	arr := make(Array, len(items))
	for i, v := range items {
		arr[i] = v
	}
	return arr
}

func ArrayOfInts(items []int) Array {
	arr := make(Array, len(items))
	for i, v := range items {
		arr[i] = v
	}
	return arr
}

func ArrayOfStrings(items []string) Array {
	arr := make(Array, len(items))
	for i, v := range items {
		arr[i] = v
	}
	return arr
}

func In(lhs interface{}, items Array) Expr {
	if len(items) == 0 {
		return FALSE
	} else if len(items) == 1 {
		return Eq(lhs, items[0])
	} else {
		return Fragment{lhs, Raw("IN"), items}
	}
}

type And []interface{}

func (v And) AppendToSQLBuilder(b *Builder) {
	switch len(v) {
	case 0:
		b.AppendRaw("TRUE")
	case 1:
		b.Append(v[0])
	default:
		b.AppendRaw("(")
		for i, item := range v {
			if i > 0 {
				b.AppendRaw(" AND ")
			}
			b.Append(item)
		}
		b.AppendRaw(")")
	}
}

type Or []interface{}

func (v Or) AppendToSQLBuilder(b *Builder) {
	switch len(v) {
	case 0:
		b.AppendRaw("FALSE")
	case 1:
		b.Append(v[0])
	default:
		b.AppendRaw("(")
		for i, item := range v {
			if i > 0 {
				b.AppendRaw(" OR ")
			}
			b.Append(item)
		}
		b.AppendRaw(")")
	}
}

func Not(v interface{}) Expr {
	return Fragment{Raw("NOT"), v}
}

func Asc(v Expr) Expr {
	return Fragment{v, Raw("ASC")}
}

func Desc(v Expr) Expr {
	return Fragment{v, Raw("DESC")}
}

func Dot(a, b Expr) Expr {
	return Fragment{a, Raw("."), b}
}

type funcExpr struct {
	name string
	args []interface{}
}

func Func(name string, args ...interface{}) Expr {
	return &funcExpr{name, args}
}

func (v funcExpr) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw(v.name)
	b.AppendRaw("(")
	for i, arg := range v.args {
		if i > 0 {
			b.AppendRaw(",")
		}
		b.Append(arg)
	}
	b.AppendRaw(")")
}

func Max(v Expr) Expr {
	return Func("MAX", v)
}

func Min(v Expr) Expr {
	return Func("MIN", v)
}

func Count(v Expr) Expr {
	return Func("COUNT", v)
}

func As(v Expr, name Column) Expr {
	return Fragment{v, Raw("AS"), name}
}

func Parens(v Expr) Expr {
	return parenthesized{v}
}

type parenthesized struct {
	v Expr
}

func (v parenthesized) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw("(")
	v.v.AppendToSQLBuilder(b)
	b.AppendRaw(")")
}

func Qualified(exprs ...Expr) Expr {
	return qualified{exprs}
}

type qualified struct {
	items []Expr
}

func (v qualified) AppendToSQLBuilder(b *Builder) {
	for i, item := range v.items {
		if i > 0 {
			b.AppendRaw(".")
		}
		item.AppendToSQLBuilder(b)
	}
}
