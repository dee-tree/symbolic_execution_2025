// Package symbolic содержит конкретные реализации символьных выражений
package symbolic

import "fmt"

// SymbolicExpression - базовый интерфейс для всех символьных выражений
type SymbolicExpression interface {
	// Type возвращает тип выражения
	Type() ExpressionType

	// String возвращает строковое представление выражения
	String() string

	// Accept принимает visitor для обхода дерева выражений
	Accept(visitor Visitor) interface{}
}

// SymbolicVariable представляет символьную переменную
type SymbolicVariable struct {
	Name     string
	ExprType ExpressionType
}

// NewSymbolicVariable создаёт новую символьную переменную
func NewSymbolicVariable(name string, exprType ExpressionType) *SymbolicVariable {
	return &SymbolicVariable{
		Name:     name,
		ExprType: exprType,
	}
}

// Type возвращает тип переменной
func (sv *SymbolicVariable) Type() ExpressionType {
	return sv.ExprType
}

// String возвращает строковое представление переменной
func (sv *SymbolicVariable) String() string {
	return sv.Name
}

// Accept реализует Visitor pattern
func (sv *SymbolicVariable) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariable(sv)
}

// IntConstant представляет целочисленную константу
type IntConstant struct {
	Value int64
}

// NewIntConstant создаёт новую целочисленную константу
func NewIntConstant(value int64) *IntConstant {
	return &IntConstant{Value: value}
}

// Type возвращает тип константы
func (ic *IntConstant) Type() ExpressionType {
	return ExpressionType{Kind: IntType}
}

// String возвращает строковое представление константы
func (ic *IntConstant) String() string {
	return fmt.Sprintf("%d", ic.Value)
}

// Accept реализует Visitor pattern
func (ic *IntConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitIntConstant(ic)
}

// BoolConstant представляет булеву константу
type BoolConstant struct {
	Value bool
}

// NewBoolConstant создаёт новую булеву константу
func NewBoolConstant(value bool) *BoolConstant {
	return &BoolConstant{Value: value}
}

// Type возвращает тип константы
func (bc *BoolConstant) Type() ExpressionType {
	return ExpressionType{Kind: BoolType}
}

// String возвращает строковое представление константы
func (bc *BoolConstant) String() string {
	return fmt.Sprintf("%t", bc.Value)
}

// Accept реализует Visitor pattern
func (bc *BoolConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitBoolConstant(bc)
}

// BinaryOperation представляет бинарную операцию
type BinaryOperation struct {
	Left     SymbolicExpression
	Right    SymbolicExpression
	Operator BinaryOperator
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// NewBinaryOperation создаёт новую бинарную операцию
func NewBinaryOperation(left, right SymbolicExpression, op BinaryOperator) *BinaryOperation {
	// Создать новую бинарную операцию и проверить совместимость типов
	if left.Type() != right.Type() {
		panic("Types mismatch for binary operation")
	}

	switch op {
	case ADD, SUB, MUL, DIV, MOD:
		if left.Type().Kind != IntType {
			panic("Types mismatch for binary arithmetic operation")
		}
	}

	return &BinaryOperation{
		Left:     left,
		Right:    right,
		Operator: op,
	}
}

// Type возвращает результирующий тип операции
func (bo *BinaryOperation) Type() ExpressionType {
	// Определить результирующий тип на основе операции и типов операндов
	// Например: int + int = int, int < int = bool

	switch bo.Operator {
	case ADD, SUB, MUL, DIV, MOD:
		return bo.Left.Type()
	case EQ, NE, LT, LE, GT, GE:
		return ExpressionType{Kind: BoolType}
	}

	panic("Bad operation")
}

// String возвращает строковое представление операции
func (bo *BinaryOperation) String() string {
	// Формат: "(left operator right)"

	return fmt.Sprintf("(%v %v %v)", bo.Left.String(), bo.Operator.String(), bo.Right.String())
}

// Accept реализует Visitor pattern
func (bo *BinaryOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryOperation(bo)
}

// LogicalOperation представляет логическую операцию
type LogicalOperation struct {
	Operands []SymbolicExpression
	Operator LogicalOperator
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// NewLogicalOperation создаёт новую логическую операцию
func NewLogicalOperation(operands []SymbolicExpression, op LogicalOperator) *LogicalOperation {
	// Создать логическую операцию и проверить типы операндов
	for _, operand := range operands {
		if operand.Type().Kind != BoolType {
			panic("Bad logical operand type")
		}
	}

	if len(operands) == 0 {
		panic("Empty logical operands")
	}

	if op == NOT && len(operands) != 1 {
		panic("NOT operation must have exactly one operand")
	}

	if op == IMPLIES && len(operands) != 2 {
		panic("IMPLIES operation must have exactly two operands")
	}

	return &LogicalOperation{Operands: operands, Operator: op}
}

// Type возвращает тип логической операции (всегда bool)
func (lo *LogicalOperation) Type() ExpressionType {
	return ExpressionType{Kind: BoolType}
}

// String возвращает строковое представление логической операции
func (lo *LogicalOperation) String() string {
	// Для NOT: "!operand"
	// Для AND/OR: "(operand1 && operand2 && ...)"
	// Для IMPLIES: "(operand1 => operand2)"

	if lo.Operator == NOT {
		return fmt.Sprintf("!%v", lo.Operands[0].String())
	}

	sexpr := ""
	sop := "&&"

	switch lo.Operator {
	case OR:
		sop = "||"
	case IMPLIES:
		sop = "=>"
	}

	for i, operand := range lo.Operands {
		switch i {
		case 0:
			sexpr += operand.String()
		default:
			sexpr += fmt.Sprintf(" %v %v", sop, operand.String())
		}
	}

	return fmt.Sprintf("(%v)", sexpr)
}

// Accept реализует Visitor pattern
func (lo *LogicalOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitLogicalOperation(lo)
}

// Операторы для бинарных выражений
type BinaryOperator int

const (
	// Арифметические операторы
	ADD BinaryOperator = iota
	SUB
	MUL
	DIV
	MOD

	// Операторы сравнения
	EQ // равно
	NE // не равно
	LT // меньше
	LE // меньше или равно
	GT // больше
	GE // больше или равно

	// Array operations
	AGET // array indexing
)

// String возвращает строковое представление оператора
func (op BinaryOperator) String() string {
	switch op {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	case MOD:
		return "%"
	case EQ:
		return "=="
	case NE:
		return "!="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	case AGET:
		return "[]"
	default:
		return "unknown"
	}
}

// Логические операторы
type LogicalOperator int

const (
	AND LogicalOperator = iota
	OR
	NOT
	IMPLIES
)

// String возвращает строковое представление логического оператора
func (op LogicalOperator) String() string {
	switch op {
	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"
	case IMPLIES:
		return "=>"
	default:
		return "unknown"
	}
}

type UnaryOperator int

const (
	// NEG - arithmetic unary minus
	NEG UnaryOperator = iota
)

func (op UnaryOperator) String() string {
	switch op {
	case NEG:
		return "-"
	default:
		return "unknown"
	}
}

// UnaryOperation представляет унарную операцию
type UnaryOperation struct {
	Expr     SymbolicExpression
	Operator UnaryOperator
}

// NewUnaryOperation создаёт новую унарную операцию
func NewUnaryOperation(expr SymbolicExpression, op UnaryOperator) *UnaryOperation {
	// Создать новую унарную операцию и проверить совместимость типов
	if expr.Type().Kind != IntType {
		panic("Types mismatch for unary operation")
	}

	return &UnaryOperation{
		Expr:     expr,
		Operator: op,
	}
}

// Type возвращает результирующий тип операции
func (uo *UnaryOperation) Type() ExpressionType {
	// Определить результирующий тип на основе операции и типа операнда

	switch uo.Operator {
	case NEG:
		return ExpressionType{Kind: IntType}
	}

	panic("Bad operation")
}

// String возвращает строковое представление операции
func (uo *UnaryOperation) String() string {
	// Формат: "(operator expr)"

	return fmt.Sprintf("(%v %v)", uo.Expr.String(), uo.Operator.String())
}

// Accept реализует Visitor pattern
func (uo *UnaryOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryOperation(uo)
}

type TernaryOperator int

const (
	IFELSE TernaryOperator = iota
)

func (op TernaryOperator) String() string {
	switch op {
	case IFELSE:
		return "?"
	default:
		panic("unknown ternary operator")
	}
}

type TernaryOperation struct {
	// Better to call like "expr1", "expr2", "expr3", but since we don't have other ternary operators - this naming appears
	Condition SymbolicExpression
	Then      SymbolicExpression
	Else      SymbolicExpression
	Operator  TernaryOperator
}

func NewTernaryOperation(condition SymbolicExpression, then SymbolicExpression, els SymbolicExpression, operator TernaryOperator) *TernaryOperation {
	if condition.Type().Kind != BoolType {
		panic("Types mismatch for if-condition")
	}
	if then.Type().Kind != els.Type().Kind {
		panic("Types mismatch for then/else branches")
	}

	return &TernaryOperation{
		Condition: condition,
		Then:      then,
		Else:      els,
		Operator:  operator,
	}
}

func (op *TernaryOperation) Type() ExpressionType {
	switch op.Operator {
	case IFELSE:
		return op.Then.Type()
	default:
		panic("Bad operation")
	}
}

func (op *TernaryOperation) String() string {
	switch op.Operator {
	case IFELSE:
		return fmt.Sprintf("(if (%v) %v else %v)", op.Condition.String(), op.Then.String(), op.Else.String())
	default:
		panic("Bad operation")
	}
}

func (op *TernaryOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitTernaryOperation(op)
}

// TODO: Добавьте дополнительные типы выражений по необходимости:
// + UnaryOperation (унарные операции: -x, !x)
// - ArrayAccess (доступ к элементам массива: arr[index])
// - FunctionCall (вызовы функций: f(x, y))
// - ConditionalExpression (тернарный оператор: condition ? true_expr : false_expr)
