package internal

import "fmt"

func IsBool(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return Null, fmt.Errorf("boolean? requires exactly 1 argument")
	}

	switch a := args[0].(type) {
	case bool:
		return true, nil
	default:
		return false, nil
	}
}

func And(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return Null, fmt.Errorf("and requires at least 2 arguments")
	}

	for _, arg := range args {
		switch a := arg.(type) {
		case bool:
			if !a {
				return false, nil
			}
		default:
			return Null, fmt.Errorf("provided non-bool argument to and")
		}
	}

	return true, nil
}

func Or(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return Null, fmt.Errorf("or requires at least 2 arguments")
	}

	for _, arg := range args {
		switch a := arg.(type) {
		case bool:
			if a {
				return true, nil
			}
		default:
			return Null, fmt.Errorf("provided non-bool argument to and")
		}
	}

	return false, nil
}

func Eq(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return Null, fmt.Errorf("= requires exactly 2 arguments")
	}

	switch a := args[0].(type) {
	case int64:
		switch b := args[1].(type) {
		case int64:
			return a == b, nil
		case float64:
			return float64(a) == b, nil
		default:
			return Null, fmt.Errorf("second argument to = was not a number")
		}
	case float64:
		switch b := args[1].(type) {
		case int64:
			return a == float64(b), nil
		case float64:
			return a == b, nil
		default:
			return Null, fmt.Errorf("second argument to = was not a number")
		}
	default:
		return Null, fmt.Errorf("first argument to = was not a number")
	}
}

func Gt(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return Null, fmt.Errorf("> requires exactly 2 arguments")
	}

	switch a := args[0].(type) {
	case int64:
		switch b := args[1].(type) {
		case int64:
			return a > b, nil
		case float64:
			return float64(a) > b, nil
		default:
			return Null, fmt.Errorf("second argument to > was not a number")
		}
	case float64:
		switch b := args[1].(type) {
		case int64:
			return a > float64(b), nil
		case float64:
			return a > b, nil
		default:
			return Null, fmt.Errorf("second argument to > was not a number")
		}
	default:
		return Null, fmt.Errorf("first argument to > was not a number")
	}
}
