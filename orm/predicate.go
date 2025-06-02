package orm

type op string

const (
	opEq  op = "="
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

func (o op) String() string {
	return string(o)
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

//func Eq(column string, right any) Predicate {
//	return Predicate{
//		Column: column,
//		Op:     "=",
//		Arg:    right,
//	}
//}

type Column struct {
	name string
}

func C(name string) Column {
	return Column{name: name}
}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opEq,
		right: value{
			val: arg,
		},
	}
}
func (c Column) expr() {}
func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}
func (Predicate) expr() {}

type value struct {
	val any
}

func (value) expr() {}

// Expression 标记接口代表，表达式
type Expression interface {
	expr()
}
