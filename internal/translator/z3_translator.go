// Package translator содержит реализацию транслятора в Z3
package translator

import (
	"fmt"
	"math/big"
	"symbolic-execution-course/internal/symbolic"

	"github.com/ebukreev/go-z3/z3"
)

// Z3Translator транслирует символьные выражения в Z3 формулы
type Z3Translator struct {
	ctx    *z3.Context
	config *z3.Config
	vars   map[string]z3.Value // Кэш переменных
}

// NewZ3Translator создаёт новый экземпляр Z3 транслятора
func NewZ3Translator() *Z3Translator {
	config := &z3.Config{}
	ctx := z3.NewContext(config)

	return &Z3Translator{
		ctx:    ctx,
		config: config,
		vars:   make(map[string]z3.Value),
	}
}

// GetContext возвращает Z3 контекст
func (zt *Z3Translator) GetContext() interface{} {
	return zt.ctx
}

// Reset сбрасывает состояние транслятора
func (zt *Z3Translator) Reset() {
	zt.vars = make(map[string]z3.Value)
}

// Close освобождает ресурсы
func (zt *Z3Translator) Close() {
	// Z3 контекст закрывается автоматически
}

// TranslateExpression транслирует символьное выражение в Z3
func (zt *Z3Translator) TranslateExpression(expr symbolic.SymbolicExpression) (interface{}, error) {
	return expr.Accept(zt), nil
}

func (zt *Z3Translator) translateType(ty symbolic.ExpressionType) z3.Sort {
	switch ty.Kind {
	case symbolic.IntType:
		return zt.ctx.IntSort()
	case symbolic.BoolType:
		return zt.ctx.BoolSort()
	case symbolic.ArrayType:
		return zt.ctx.ArraySort(zt.ctx.IntSort(), zt.translateType(*ty.Inner))
	default:
		panic("Unhandled type " + ty.String())
	}
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// VisitVariable транслирует символьную переменную в Z3
func (zt *Z3Translator) VisitVariable(expr *symbolic.SymbolicVariable) interface{} {
	// Проверить, есть ли переменная в кэше
	// Если нет - создать новую Z3 переменную соответствующего типа
	// Добавить в кэш и вернуть

	// Подсказки:
	// - Используйте zt.ctx.IntConst(name) для int переменных
	// - Используйте zt.ctx.BoolConst(name) для bool переменных
	// - Храните переменные в zt.vars для повторного использования

	vrbl := zt.vars[expr.Name]
	if vrbl != nil {
		return vrbl
	}

	switch expr.Type().Kind {
	case symbolic.IntType:
		vrbl = zt.ctx.IntConst(expr.Name)
	case symbolic.BoolType:
		vrbl = zt.ctx.IntConst(expr.Name)
	case symbolic.ArrayType:
		zt.ctx.Const(expr.Name, zt.translateType(expr.Type()))
	default:
		panic("unknown variable type")
	}

	zt.vars[expr.Name] = vrbl
	return vrbl
}

// VisitIntConstant транслирует целочисленную константу в Z3
func (zt *Z3Translator) VisitIntConstant(expr *symbolic.IntConstant) interface{} {
	// Создать Z3 константу с помощью zt.ctx.FromBigInt или аналогичного метода

	return zt.ctx.FromBigInt(big.NewInt(expr.Value), zt.ctx.IntSort())
}

// VisitBoolConstant транслирует булеву константу в Z3
func (zt *Z3Translator) VisitBoolConstant(expr *symbolic.BoolConstant) interface{} {
	// Использовать zt.ctx.FromBool для создания Z3 булевой константы

	return zt.ctx.FromBool(expr.Value)
}

// VisitBinaryOperation транслирует бинарную операцию в Z3
func (zt *Z3Translator) VisitBinaryOperation(expr *symbolic.BinaryOperation) interface{} {
	// TODO: Реализовать
	// 1. Транслировать левый и правый операнды
	// 2. В зависимости от оператора создать соответствующую Z3 операцию

	// Подсказки по операциям в Z3:
	// - Арифметические: left.Add(right), left.Sub(right), left.Mul(right), left.Div(right)
	// - Сравнения: left.Eq(right), left.LT(right), left.LE(right), etc.
	// - Приводите типы: left.(z3.Int), right.(z3.Int) для int операций

	switch expr.Operator {
	case symbolic.ADD:
		return expr.Left.Accept(zt).(z3.Int).Add(expr.Left.Accept(zt).(z3.Int))
	case symbolic.SUB:
		return expr.Left.Accept(zt).(z3.Int).Sub(expr.Left.Accept(zt).(z3.Int))
	case symbolic.MUL:
		return expr.Left.Accept(zt).(z3.Int).Mul(expr.Left.Accept(zt).(z3.Int))
	case symbolic.DIV:
		return expr.Left.Accept(zt).(z3.Int).Div(expr.Left.Accept(zt).(z3.Int))
	case symbolic.MOD:
		return expr.Left.Accept(zt).(z3.Int).Mod(expr.Left.Accept(zt).(z3.Int))

	case symbolic.EQ:
		switch expr.Left.Type().Kind {
		case symbolic.IntType:
			return expr.Left.Accept(zt).(z3.Int).Eq(expr.Right.Accept(zt).(z3.Int))
		case symbolic.BoolType:
			return expr.Left.Accept(zt).(z3.Bool).Eq(expr.Right.Accept(zt).(z3.Bool))
		case symbolic.ArrayType:
			return expr.Left.Accept(zt).(z3.Array).Eq(expr.Right.Accept(zt).(z3.Array))
		default:
			panic("Unhandled EQ type")
		}
	case symbolic.NE:
		switch expr.Left.Type().Kind {
		case symbolic.IntType:
			return expr.Left.Accept(zt).(z3.Int).NE(expr.Right.Accept(zt).(z3.Int))
		case symbolic.BoolType:
			return expr.Left.Accept(zt).(z3.Bool).NE(expr.Right.Accept(zt).(z3.Bool))
		case symbolic.ArrayType:
			return expr.Left.Accept(zt).(z3.Array).NE(expr.Right.Accept(zt).(z3.Array))
		default:
			panic("NE over non-bool or non-int expression")
		}
	case symbolic.LT:
		return expr.Left.Accept(zt).(z3.Int).LT(expr.Right.Accept(zt).(z3.Int))
	case symbolic.LE:
		return expr.Left.Accept(zt).(z3.Int).LE(expr.Right.Accept(zt).(z3.Int))
	case symbolic.GT:
		return expr.Left.Accept(zt).(z3.Int).GT(expr.Right.Accept(zt).(z3.Int))
	case symbolic.GE:
		return expr.Left.Accept(zt).(z3.Int).GE(expr.Right.Accept(zt).(z3.Int))
	}

	panic(fmt.Sprintf("unimplemented visitor for binary operation: %v", expr.Operator.String()))
}

// VisitLogicalOperation транслирует логическую операцию в Z3
func (zt *Z3Translator) VisitLogicalOperation(expr *symbolic.LogicalOperation) interface{} {
	// TODO: Реализовать
	// 1. Транслировать все операнды
	// 2. Применить соответствующую логическую операцию

	// Подсказки:
	// - AND: zt.ctx.And(operands...)
	// - OR: zt.ctx.Or(operands...)
	// - NOT: operand.Not() (для единственного операнда)
	// - IMPLIES: antecedent.Implies(consequent)

	var args []z3.Bool
	for _, arg := range expr.Operands {
		args = append(args, arg.Accept(zt).(z3.Bool))
	}

	switch expr.Operator {
	case symbolic.AND:
		return args[0].And(args[1:]...)
	case symbolic.OR:
		return args[0].Or(args[1:]...)
	case symbolic.NOT:
		return args[0].Not()
	case symbolic.IMPLIES:
		return args[0].Implies(args[1])

	default:
		panic(fmt.Sprintf("unimplemented visitor for logical operation: %v", expr.Operator.String()))
	}
}

func (zt *Z3Translator) VisitUnaryOperation(expr *symbolic.UnaryOperation) interface{} {
	if expr.Operator != symbolic.NEG {
		panic("unimplemented unary op visitor")
	}

	return expr.Expr.Accept(zt).(z3.Int).Neg()
}

func (zt *Z3Translator) VisitTernaryOperation(expr *symbolic.TernaryOperation) interface{} {
	if expr.Operator != symbolic.IFELSE {
		panic("unimplemented ternary op visitor")
	}

	return expr.Condition.Accept(zt).(z3.Bool).IfThenElse(expr.Then.Accept(zt).(z3.Value), expr.Else.Accept(zt).(z3.Value))
}

// Вспомогательные методы

// castToZ3Type приводит значение к нужному Z3 типу
func (zt *Z3Translator) castToZ3Type(value interface{}, targetType symbolic.ExpressionType) (z3.Value, error) {
	// TODO: Реализовать (вспомогательный метод)
	// Безопасно привести interface{} к конкретному Z3 типу
	panic("не реализовано")
}
