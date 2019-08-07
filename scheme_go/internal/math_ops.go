package internal

import (
	"fmt"
	"math"
)

type Procedure = func(args ...interface{}) (interface{}, error)

var Null []interface{} = []interface{}{}

func isNull(i []interface{}) bool {
	return len(i) == 0
}

// I'm not going to support pairs for now
func Cons(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return Null, fmt.Errorf("cons requires exactly 2 arguments")
	}

	// forcing second arg to be a list or nil
	switch a := args[1].(type) {
	case []interface{}:
		if isNull(a) {
			return []interface{}{args[0]}, nil
		}
		return append([]interface{}{args[0]}, a...), nil
	default:
		return Null, fmt.Errorf("I don't support cons with a non-list (or nil) as the second arg")
	}
}

func IsEmpty(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return Null, fmt.Errorf("empty? requires exactly 1 argument")
	}

	switch a := args[0].(type) {
	case []interface{}:
		return isNull(a), nil
	default:
		return Null, fmt.Errorf("empty? must be called on a list")
	}
}

func Car(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return Null, fmt.Errorf("car requires exactly 1 argument")
	}

	switch a := args[0].(type) {
	case []interface{}:
		if isNull(a) {
			return Null, fmt.Errorf("cannot use car on empty list")
		}
		return a[0], nil
	default:
		return Null, fmt.Errorf("cannot call car on non-list")
	}
}

func Cdr(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return Null, fmt.Errorf("cdr requires exactly 1 argument")
	}

	switch a := args[0].(type) {
	case []interface{}:
		if isNull(a) {
			return Null, fmt.Errorf("cannot use cdr on empty list")
		}

		//		if a[1:] == Null {
		//			return Null, nil
		//		}

		return a[1:], nil
	default:
		return nil, fmt.Errorf("cannot call cdr on non-list")
	}
}

// Plus is a function that expects 2 or more numbers (of any type)
func Plus(args ...interface{}) (interface{}, error) {
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
}

// Minus is a function that expects 2 or more numbers (of any type)
func Minus(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("- requires at least two arguments")
	}
	var total float64 = 0
	var wasAFloat bool
	var isFirst bool = true
	for _, arg := range args {
		switch a := arg.(type) {
		case int64:
			if isFirst {
				isFirst = false
				total = float64(a)
			} else {
				total -= float64(a)
			}
		case float64:
			wasAFloat = true
			if isFirst {
				isFirst = false
				total = a
			} else {
				total -= a
			}
		case interface{}:
			return nil, fmt.Errorf("non-number argument to -: %v", a)
		}
	}

	if !wasAFloat {
		return math.Round(total), nil
	} else {
		return total, nil
	}
}

// Mult is a function that expects 2 or more numbers (of any type)
func Mult(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("* requires at least two arguments")
	}
	var total float64 = 0
	var wasAFloat bool
	var isFirst bool = true
	for _, arg := range args {
		switch a := arg.(type) {
		case int64:
			if isFirst {
				isFirst = false
				total = float64(a)
			} else {
				total *= float64(a)
			}
		case float64:
			wasAFloat = true
			if isFirst {
				isFirst = false
				total = a
			} else {
				total *= a
			}
		case interface{}:
			return nil, fmt.Errorf("non-number argument to *: %v", a)
		}
	}

	if !wasAFloat {
		return math.Round(total), nil
	} else {
		return total, nil
	}
}

// Div is a function that expects 2 or more numbers (of any type)
func Div(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("/ requires exactly two arguments")
	}
	var total float64 = 0
	var wasAFloat bool
	var isFirst bool = true
	for _, arg := range args {
		switch a := arg.(type) {
		case int64:
			if isFirst {
				isFirst = false
				total = float64(a)
			} else {
				total /= float64(a)
			}
		case float64:
			wasAFloat = true
			if isFirst {
				isFirst = false
				total = a
			} else {
				total /= a
			}
		case interface{}:
			return nil, fmt.Errorf("non-number argument to /: %v", a)
		}
	}

	if !wasAFloat {
		// having even integer division is possible, this is a check if it was even(ish)
		// there is probably a better way, but I want to just move on
		rounded := math.Round(total)
		if math.Abs(rounded-total) < 0.00001 {
			return rounded, nil
		}
		return total, nil
	} else {
		return total, nil
	}
}
