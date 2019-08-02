package main

import (
	"fmt"
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

//const baseEnv map[string]interface{}{
//	"+":
//}

func tokenize(chars string) []string {
	chars = strings.Replace(chars, "(", " ( ", -1)
	chars = strings.Replace(chars, ")", " ) ", -1)
	return strings.Fields(chars)
}

//func make_token_channel(tokens []string) chan string {
//	c := make(chan string, (1<<16 - 1))
//	for _, token := range tokens {
//		c <- token
//	}
//	close(c)
//	return c
//}

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

func printSliceWithTypes(s []interface{}) {
	fmt.Printf("[")
	for _, el := range s {
		switch e := el.(type) {
		case []interface{}:
			printSliceWithTypes(e)
			fmt.Printf(", ")
		case interface{}:
			fmt.Printf("%s: %v, ", reflect.TypeOf(el), el)
		}
	}
	fmt.Printf("]")
}

func parse(program string) (interface{}, error) {
	ast, _, err := readFromTokens(tokenize(program))
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err)
	}
	return ast, nil
}

func main() {
	program := "(+ (* 3 4) 1)"
	ast, err := parse(program)
	if err != nil {
		panic(err)
	}
	switch a := ast.(type) {
	case []interface{}:
		printSliceWithTypes(a)
	case interface{}:
		fmt.Printf("ast is not a slice and is: %s: %v\n", reflect.TypeOf(a), a)

	}

	//var e []interface{}
	//fmt.Printf("%#v\n", e)
}
