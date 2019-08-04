package internal

import (
	"fmt"
	"math"
)

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
