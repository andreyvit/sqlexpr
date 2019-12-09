package sqlexpr

import (
	"fmt"
	"strings"
)

type Expr interface {
	AppendToSQLBuilder(b *Builder)
}

type Builder struct {
	buf  strings.Builder
	args []interface{}
	last rune
}

func Build(e Expr) (sql string, args []interface{}) {
	var b Builder
	b.Append(e)
	return b.SQLArgs()
}

func (b *Builder) SQLArgs() (sql string, args []interface{}) {
	return b.SQL(), b.Args()
}
func (b *Builder) SQL() string {
	return b.buf.String()
}
func (b *Builder) Args() []interface{} {
	return b.args
}

func FormatSQLArgs(sql string, args []interface{}) string {
	var buf strings.Builder
	buf.WriteString(sql)
	if len(args) > 0 {
		buf.WriteString(" [")
		for i, arg := range args {
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprint(&buf, arg)
		}
		buf.WriteString("]")
	}
	return buf.String()
}

func (b *Builder) String() string {
	return FormatSQLArgs(b.SQL(), b.Args())
}

func (b *Builder) AppendRaw(s string) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return
	}

	if b.last != 0 {
		last, first := b.last, rune(s[0])
		if (isWordChar(last) && !skipSpaceBefore(first) && !isComma(first)) || (isWordChar(first) && !skipSpaceAfter(last)) || isComma(last) {
			b.buf.WriteByte(' ')
		}
	}

	b.buf.WriteString(s)
	b.last = rune(s[len(s)-1])
}

func (b *Builder) AppendName(s string) {
	// TODO: quote name if necessary?
	b.AppendRaw(s)
}

func (b *Builder) AppendExpr(expr Expr) {
	if expr != nil {
		expr.AppendToSQLBuilder(b)
	}
}

func (b *Builder) Append(item interface{}) {
	if expr, ok := item.(Expr); ok {
		expr.AppendToSQLBuilder(b)
	} else {
		placeholder := dialect.FormatPlaceholder(len(b.args))
		b.args = append(b.args, item)
		b.AppendRaw(placeholder)
	}
}

func (b *Builder) AppendAll(items ...interface{}) {
	for _, item := range items {
		b.Append(item)
	}
}

func isWordChar(r rune) bool {
	return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '_' || r == '$'
}

func skipSpaceAfter(r rune) bool {
	return r == '(' || r == '.'
}

func skipSpaceBefore(r rune) bool {
	return r == ')' || r == '.'
}

func isComma(r rune) bool {
	return r == ','
}

type Raw string

func (v Raw) AppendToSQLBuilder(b *Builder) {
	b.AppendRaw(string(v))
}

type Table string
type Column string

func (v Table) AppendToSQLBuilder(b *Builder) {
	b.AppendName(string(v))
}

func (v Column) AppendToSQLBuilder(b *Builder) {
	b.AppendName(string(v))
}

const (
	Empty     = Raw("")
	Star      = Raw("*")
	TRUE      = Raw("TRUE")
	FALSE     = Raw("FALSE")
	NOW       = Raw("NOW()")
	ForUpdate = Raw("FOR UPDATE")
)

func MaybeForUpdate(forUpdate bool) Expr {
	if forUpdate {
		return ForUpdate
	} else {
		return Empty
	}
}
