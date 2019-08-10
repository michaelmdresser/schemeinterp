package internal

import "fmt"

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
