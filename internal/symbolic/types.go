// Package symbolic определяет базовые типы символьных выражений
package symbolic

// ExpressionType представляет тип символьного выражения
type ExpressionType struct {
	Kind  TypeKind
	Inner *ExpressionType
}
type TypeKind int

const (
	IntType TypeKind = iota
	BoolType
	ArrayType
	// Добавьте другие типы по необходимости
)

// String возвращает строковое представление типа
func (et ExpressionType) String() string {
	switch et.Kind {
	case IntType:
		return "int"
	case BoolType:
		return "bool"
	case ArrayType:
		return "array[" + et.Inner.String() + "]"
	default:
		return "unknown"
	}
}
