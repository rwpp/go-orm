package errs

import (
	"errors"
	"fmt"
)

var (
	ErrorPointerOnly = errors.New("只支持一级指针指向的结构体")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("不支持的表达式类型%v", expr)
}
func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm:未知字段%s", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm:非法标签%s", pair)
}
