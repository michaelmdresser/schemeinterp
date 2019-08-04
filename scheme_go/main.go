package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/michaelmdresser/scheme_go/internal"
)

// Symbol: string
// Number: int, float
// Atom: (Symbol, Number)
// List: slice
// Exp: (Atom, List)
// Env: map

type closure struct {
	env           environment
	argumentNames []string
	expr          interface{}
}

type environment struct {
	parentEnv *environment
	env       map[string]interface{}
}

func (e environment) lookup(name string) (interface{}, error) {
	if val, ok := e.env[name]; ok {
		return val, nil
	}

	if e.parentEnv == nil {
		return nil, fmt.Errorf("name %s not in environment", name)
	}

	val, err := e.parentEnv.lookup(name)
	if err != nil {
		return nil, err
	}

	return val, nil
}

var baseEnv environment = environment{
	parentEnv: nil,
	env: map[string]interface{}{
		"parentEnv": nil,
		"true":      true,
		"false":     false,
		"+":         internal.Plus,
		"-":         internal.Minus,
		"*":         internal.Mult,
		"/":         internal.Div,
	},
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

// returns expr, new environment, error
func eval(expr interface{}, env environment) (interface{}, environment, error) {
	switch a := expr.(type) {
	case []interface{}: // procedure call or keyword
		if a[0] == "if" { // conditional
			if len(a) != 4 {
				return nil, env, fmt.Errorf("incorrect number of arguments to if: %v", a)
			}
			test := a[1]
			consequence := a[2]
			alternative := a[3]

			result, env, err := eval(test, env)
			if err != nil {
				return nil, env, fmt.Errorf("failed to evaluate test %v: %s", test, err)
			}

			switch r := result.(type) {
			case bool:
				if r {
					conseqEval, env, err := eval(consequence, env)
					if err != nil {
						return nil, env, fmt.Errorf("failed to evaluate consequence %v: %s", consequence, err)
					}
					return conseqEval, env, nil
				} else {
					altEval, env, err := eval(alternative, env)
					if err != nil {
						return nil, env, fmt.Errorf("failed to evaluate alternative %v: %s", alternative, err)
					}
					return altEval, env, nil
				}
			default:
				return nil, env, fmt.Errorf("test %v did not evaluate to a bool", test)
			}
		} else if a[0] == "define" { // definition
			if len(a) != 3 {
				return nil, env, fmt.Errorf("incorrect number of arguments to define: %v", a)
			}

			symbol := a[1]
			symbolExpr := a[2]

			switch s := symbol.(type) {
			case string:
				result, env, err := eval(symbolExpr, env)
				if err != nil {
					return nil, env, fmt.Errorf("failed to evaluate symbol expression %v for symbol %s: %s", symbolExpr, symbol, err)
				}
				env.env[s] = result
				return nil, env, nil
			default:
				return nil, env, fmt.Errorf("symbol %v in define is not a string", symbol)
			}
		} else if a[0] == "lambda" {
			if len(a) != 3 {
				return nil, env, fmt.Errorf("incorrect number of arguments to lambda: %v", a)
			}

			argList := a[1]
			innerExpr := a[2]

			var argListString []string
			switch al := argList.(type) {
			case []interface{}:
				for _, arg := range al {
					switch a := arg.(type) {
					case string:
						argListString = append(argListString, a)
					default:
						return nil, env, fmt.Errorf("argument %v was not a string", arg)
					}
				}
			default:
				return nil, env, fmt.Errorf("first argument to lambda %v was not a list", a)
			}

			return closure{
				env:           env,
				argumentNames: argListString,
				expr:          innerExpr,
			}, env, nil
		} else { // procedure call
			procedure, env, err := eval(a[0], env)
			if err != nil {
				return nil, env, fmt.Errorf("could not evaluate %v: %s", a[0], err)
			}
			var arguments []interface{}
			for _, unevaluatedArg := range a[1:] {
				evaluatedArg, env, err := eval(unevaluatedArg, env)
				if err != nil {
					return nil, env, fmt.Errorf("could not evaluate %v: %s", unevaluatedArg, err)
				}
				arguments = append(arguments, evaluatedArg)
			}

			switch p := procedure.(type) {
			case func(args ...interface{}) (interface{}, error):
				ret, procedureErr := p(arguments...)
				if procedureErr != nil {
					return nil, env, fmt.Errorf("procedure with identifier %s called with arguments %v failed: %s", a[0], arguments, procedureErr)
				}
				return ret, env, nil
			case closure:
				evalEnv := p.env
				if len(arguments) != len(p.argumentNames) {
					return nil, env, fmt.Errorf("procedure expects %d arguments, got %d", len(p.argumentNames), len(arguments))
				}

				for i, arg := range arguments {
					evalEnv.env[p.argumentNames[i]] = arg
				}
				return eval(p.expr, evalEnv)
			case interface{}:
				return nil, env, fmt.Errorf("procedure obtained from map with identifier %s is not a procedure")
			}
		}
	case int64: // constant
		return a, env, nil
	case float64: // constant
		return a, env, nil
	case string: // variable reference
		val, err := env.lookup(a)
		if err != nil {
			return nil, env, fmt.Errorf("failed to reference variable: %s", err)
		}
		return val, env, nil
	}
	return nil, env, fmt.Errorf("tried to evaluate the type of %v, which is not a procedure call, constant, keyword, or reference", expr)
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

func repl() {
	env := baseEnv
	reader := bufio.NewReader(os.Stdin)
	var result interface{}
	var err error
	var ast interface{}

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		ast, err = parse(text)
		if err != nil {
			fmt.Println("failed to parse input: %s", err)
		}
		result, env, err = eval(ast, env)
		if err != nil {
			fmt.Printf("failed to eval input: %s\n", err)
		}
		fmt.Printf("%v\n", result)
	}
}

func main() {
	repl()
}
