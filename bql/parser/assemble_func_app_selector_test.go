package parser

import (
	"testing"
)

func TestAssembleFuncAppSelector(t *testing.T) {
	ps := parseStack{}
	t.Run("correct items", func(t *testing.T) {
		ps.PushComponent(0, 6, Raw{"PRE"})
		ps.PushComponent(6, 7, FuncAppAST{
			Function: FuncName("add"),
		})
		ps.PushComponent(7, 8, Raw{"[0]"})
		ps.AssembleFuncAppSelector()

		if ps.Len() != 2 {
			t.Fatalf("assembled components should be 2 but %d", ps.Len())
		}
		top := ps.Peek()
		if top == nil {
			t.Fatal("picked component should not be nil")
		}
		if top.begin != 6 || top.end != 8 {
			t.Errorf("the component should start from 6 to 8, but %d to %d", top.begin, top.end)
		}
		comp, ok := top.comp.(FuncAppSelectorAST)
		if !ok {
			t.Fatal("the component type should be 'FuncAppSelectorAST'")
		}
		if comp.Function != "add" {
			t.Errorf("the function name should be 'add' but %s", comp.Function)
		}
		if comp.Selector.Expr != "[0]" {
			t.Errorf("the function's selector should be '[0]' but %s", comp.Selector)
		}
	})
}

func TestAssembleFuncAppSelectorWrongItem(t *testing.T) {
	ps := parseStack{}
	ps.PushComponent(0, 6, Raw{"PRE"})
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("assembling logic should cause panic")
		}
	}()
	ps.AssembleFuncAppSelector()
}

func TestAssembleFuncAppSelectorEmptyItem(t *testing.T) {
	ps := parseStack{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("assembling logic should cause panic")
		}
	}()
	ps.AssembleFuncAppSelector()
}
