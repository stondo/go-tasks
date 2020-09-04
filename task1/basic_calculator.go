package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// Stack LIFO
type Stack []interface{}

// IsEmpty checks if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(str interface{}) {
	*s = append(*s, str)
}

// PushAll pushes all values from a list onto the stack
func (s *Stack) PushAll(strs []interface{}) {
	for _, str := range strs {
		*s = append(*s, str)
	}
}

// Pop remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	}

	idx := len(*s) - 1
	elem := (*s)[idx]
	*s = (*s)[:idx]
	return elem, true

}

// Head returns the top element of stack
func (s *Stack) Head() interface{} {
	if s.IsEmpty() {
		return nil
	}

	idx := len(*s) - 1
	elem := (*s)[idx]
	return elem
}

// Plus const
const Plus = '+'

// Star const
const Star = '*'

// Dot const
const Dot = '.'

func main() {
	expr := "1.5*2.1*3.5 + 375.4+9*5 * 6"

	fmt.Println(expr)

	result := eval(expr)

	fmt.Println("The result of the expression is:", result)
}

// opPrecedence returns the rank of the operator.
func opPrecedence(op rune) int {
	if op == Star {
		return 2
	}
	if op == Plus {
		return 1
	}
	return 0
}

func doMath(values *Stack, operators *Stack) {
	x, xErr := values.Pop()
	y, yErr := values.Pop()
	op, opErr := operators.Pop()

	if !xErr && !yErr && !opErr {
		fmt.Println("error:", xErr, yErr, opErr)
		panic("error popping from stack!")
	}

	xValue := x.(float64)
	yValue := y.(float64)
	opValue := op.(rune)

	switch opValue {
	case Star:
		values.Push(xValue * yValue)
		// fmt.Println(*values)
	case Plus:
		values.Push(xValue + yValue)
		// fmt.Println(*values)
	default:
		panic("Unsupported operation")
	}
}

// eval calculates the value of the expression.
func eval(expr string) float64 {
	var valuesStack Stack
	var opsStack Stack

	for i := 0; i < len(expr); i++ {

		// fmt.Println("current char:", string(expr[i]))
		// fmt.Println("current char type:", reflect.TypeOf(expr[i]))

		if unicode.IsSpace(rune(expr[i])) {
			// fmt.Println("space")
			continue
		} else if unicode.IsNumber(rune(expr[i])) || expr[i] == Dot {
			var value string = ""
			for i < len(expr) && (unicode.IsNumber(rune(expr[i])) || expr[i] == Dot) {
				value = value + (string(expr[i]))
				// fmt.Println("value: ", value)
				i++
			}

			// fmt.Println("str to be converted:", value)

			f, err := strconv.ParseFloat(value, 64)

			// fmt.Println("float64:", f)

			if err != nil {
				fmt.Println("error:", err)
				panic("Error converting string to float64")
			}

			valuesStack.Push(f)
			// fmt.Println(valuesStack)
			i--
		} else {
			// fmt.Println("else oeprator:", string(expr[i]))
			// fmt.Println("len expr[i]:", len(string(expr[i])))

			var op rune
			opHead := opsStack.Head()

			if opHead != nil {
				op = opHead.(rune)
				// fmt.Println("current op:", string(op))
			}

			topOp := opPrecedence(op)
			currentOp := opPrecedence(rune(expr[i]))
			// fmt.Println("top op prec:", topOp)
			// fmt.Println("cur op prec:", currentOp)
			// fmt.Println("opStack is empty?", opsStack.IsEmpty())
			for !opsStack.IsEmpty() && topOp > currentOp {
				doMath(&valuesStack, &opsStack)
			}

			opsStack.Push(rune(expr[i]))
			// fmt.Println(opsStack)
		}
	}

	for !opsStack.IsEmpty() {
		doMath(&valuesStack, &opsStack)
	}

	res := valuesStack.Head()
	if res == nil {
		panic("error evaluating the expression")
	}

	return res.(float64)
}
