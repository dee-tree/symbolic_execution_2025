// Package ssa предоставляет функции для построения SSA представления
package ssa

import (
	"go/token"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Builder отвечает за построение SSA из исходного кода Go
type Builder struct {
	fset *token.FileSet
}

// NewBuilder создаёт новый экземпляр Builder
func NewBuilder() *Builder {
	return &Builder{
		fset: token.NewFileSet(),
	}
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// ParseAndBuildSSA парсит исходный код Go и создаёт SSA представление
// Возвращает SSA программу и функцию по имени
func (b *Builder) ParseAndBuildSSA(source string, funcName string) (*ssa.Function, error) {
	// TODO: Реализовать
	// Шаги:
	// 1. Парсинг исходного кода с помощью go/parser
	// 2. Создание SSA программы
	// 3. Поиск нужной функции по имени

	// Подсказки:
	// - Используйте parser.ParseFile для парсинга
	// - Создайте packages.Config и загрузите пакет
	// - Используйте ssautil.CreateProgram для создания SSA
	// - Найдите функцию в SSA программе

	var cfg = packages.Config{
		Fset: b.fset,
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo,
		Overlay: map[string][]byte{
			"main.go": []byte(source),
		}}

	pkgs, err := packages.Load(&cfg)
	if err != nil {
		return nil, err
	}

	ssaprg, ssapkgs := ssautil.Packages(pkgs, ssa.SanityCheckFunctions)
	ssaprg.Build()

	for _, pkg := range ssapkgs {
		fun := pkg.Func(funcName)
		if fun != nil {
			return fun, nil
		}
	}

	return nil, nil
}
