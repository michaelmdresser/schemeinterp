package internal

type Procedure = func(args ...interface{}) (interface{}, error)

var Null []interface{} = []interface{}{}

func isNull(i []interface{}) bool {
	return len(i) == 0
}
