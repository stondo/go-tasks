package main

import (
	"testing"
)

func TestOpPrecedence(t *testing.T) {
	ans := opPrecedence(Plus)
	if ans != 1 {
		t.Errorf("opPrecedence(Plus) = %d; want +", ans)
	}

	ans2 := opPrecedence(Star)
	if ans2 != 2 {
		t.Errorf("opPrecedence(Star) = %d; want +", ans2)
	}

	ans3 := opPrecedence(Dot)
	if ans3 != 0 {
		t.Errorf("opPrecedence(Dot) = %d; want +", ans3)
	}
}

func TestDoMath(t *testing.T) {
	var valuesStack Stack
	var opsStack Stack

	valuesStack.Push(19.0)
	valuesStack.Push(9.0)

	opsStack.Push(Plus)

	doMath(&valuesStack, &opsStack)
	res := valuesStack.Head()
	if res == nil {
		panic("error evaluating the expression")
	}

	ans := res.(float64)
	if ans != 28 {
		t.Errorf("doMath(&valuesStack, &opsStack) = %f; want +", ans)
	}
}

func TestEval(t *testing.T) {
	expr := "1+2+3+4*1.0+5+6+7+8+9*1+10*1"
	ans := eval(expr)
	if ans != 55 {
		t.Errorf("eval(\"1*1+2+3+4*1.0+5+6+7+8+9*1+10*1+0.0\") = %f; want +", ans)
	}
}
