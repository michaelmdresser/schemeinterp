package main

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Symbol: string
// Number: int, float
// Atom: (Symbol, Number)
// List: slice
// Exp: (Atom, List)
// Env: map

var baseEnv = map[string]interface{}{
	"+": func(args ...interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("+ requires at least two arguments")
		}
		var total float64 = 0
		var wasAFloat bool
		for _, arg := range args {
			switch a := arg.(type) {
			case int64:
				total += float64(a)
			case float64:
				wasAFloat = true
				total += a
			case interface{}:
				return nil, fmt.Errorf("non-number argument to +: %v", a)
			}
		}

		if !wasAFloat {
			return math.Round(total), nil
		} else {
			return total, nil
		}
	},
	//	"-": func(args ...interface{}) (interface{}, error) {
	//		if len(args) < 2 {
	//			return nil, fmt.Errorf("- requires at least two arguments")
	//		}
	//		total := 0
	//		isFirstProcessed := false
	//		for _, arg := range args {
	//			switch a := arg.(type) {
	//			case int64:
	//				if !isFirstProcessed {
	//					total = a
	//					isFirstProcessed = true
	//				} else {
	//					total -= a
	//				}
	//			case float64:
	//				if !isFirstProcessed {
	//					total = a
	//					isFirstProcessed = true
	//				} else {
	//					total -= a
	//				}
	//			case interface{}:
	//				return nil, fmt.Errorf("non-number argument to -: %v", a)
	//			}
	//		}
	//	},
}

func tokenize(chars string) []string {
	chars = strings.Replace(chars, "(", " ( ", -1)
	chars = strings.Replace(chars, ")", " ) ", -1)
	return strings.Fields(chars)
}

// returns an expression, number of tokens processed, and an error
func readFromTokens(tokens []string) (expr interface{}, processedTokens int, err error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("no tokens to read")
	}

	first := tokens[0]
	tokens = tokens[1:]

	if first == "(" {
		var expr []interface{}
		processed := 1 // we processed the open paren

		for len(tokens) > 0 && tokens[0] != ")" {
			subExpr, subProcessed, err := readFromTokens(tokens)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing subexpression: %s", err)
			}
			expr = append(expr, subExpr)
			tokens = tokens[subProcessed:]
			processed += subProcessed
		}
		if len(tokens) == 0 {
			return nil, 0, fmt.Errorf("syntax error: missing )")
		}
		return expr, processed + 1, nil // add one because we processed a close paren
	} else if first == ")" {
		return nil, 0, fmt.Errorf("syntax error: unexpected )")
	} else {
		return makeAtomType(first), 1, nil
	}

}

// makeAtomType takes an Atom that is a string and enforces a type on it, trying int then float and falls back on string (string is "symbol")
func makeAtomType(atom string) interface{} {
	if i, err := strconv.ParseInt(atom, 10, 64); err == nil {
		return i
	} else if f, err := strconv.ParseFloat(atom, 64); err == nil {
		return f
	} else {
		return atom
	}
}

func parse(program string) (interface{}, error) {
	ast, _, err := readFromTokens(tokenize(program))
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err)
	}
	return ast, nil
}

func eval(expr interface{}, env map[string]interface{}) (interface{}, error) {
	switch a := expr.(type) {
	case []interface{}: // procedure call
		procedure, err := eval(a[0], env)
		if err != nil {
			return nil, fmt.Errorf("could not evaluate %v: %s", a[0], err)
		}
		var arguments []interface{}
		for _, unevaluatedArg := range a[1:] {
			evaluatedArg, err := eval(unevaluatedArg, env)
			if err != nil {
				return nil, fmt.Errorf("could not evaluate %v: %s", unevaluatedArg, err)
			}
			arguments = append(arguments, evaluatedArg)
		}

		switch p := procedure.(type) {
		case func(args ...interface{}) (interface{}, error):
			ret, procedureErr := p(arguments...)
			if procedureErr != nil {
				return nil, fmt.Errorf("procedure with identifier %s called with arguments %v failed: %s", a[0], arguments, procedureErr)
			}
			return ret, nil
		case interface{}:
			return nil, fmt.Errorf("procedure obtained from map with identifier %s is not a procedure")
		}
	case int64: // constant
		return a, nil
	case float64: // constant
		return a, nil
	case string: // any number of things
		if a == "if" { // conditional

		} else if a == "define" { // definition

		} else { // variable reference
			if _, ok := env[a]; !ok {
				return nil, fmt.Errorf("symbol %s does not exist in environment", a)
			}
			return env[a], nil
		}
	}
	return nil, fmt.Errorf("tried to evaluate the type of %v, which is not a procedure call, constant, keyword, or reference", expr)
}

func printSliceWithTypes(s []interface{}) {
	fmt.Printf("[")
	for i, el := range s {
		switch e := el.(type) {
		case []interface{}:
			printSliceWithTypes(e)
		case interface{}:
			fmt.Printf("%s: %v", reflect.TypeOf(el), el)
		}
		if i < len(s)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("]")
}

func main() {
	program := "(+ 4 1)"
	ast, err := parse(program)
	if err != nil {
		panic(err)
	}
	switch a := ast.(type) {
	case []interface{}:
		printSliceWithTypes(a)
		fmt.Println()
	case interface{}:
		fmt.Printf("ast is not a slice and is: %s: %v\n", reflect.TypeOf(a), a)
	}

	result, err := eval(ast, baseEnv)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", result)
}
