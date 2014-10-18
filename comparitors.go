package expect

import (
	"reflect"
)

const (
	Equals ComparitorType = iota
	NotEqual
	LessThan
	LessThanOrEqual
	GreaterThan
	GreaterThanOrEqual
)

type ComparitorType int
type Comparitor func(k reflect.Kind, a, b interface{}) bool
type Resolver func(a, b interface{}) bool

var (
	Comparitors = map[reflect.Kind]map[ComparitorType]Resolver{
		reflect.Int64: map[ComparitorType]Resolver{
			LessThan: func(a, b interface{}) bool { return a.(int64) < b.(int64) },
		},
		reflect.Uint64: map[ComparitorType]Resolver{
			LessThan: func(a, b interface{}) bool { return a.(uint64) < b.(uint64) },
		},
		reflect.Float32: map[ComparitorType]Resolver{
			LessThan: func(a, b interface{}) bool { return a.(float32) < b.(float32) },
		},
		reflect.Float64: map[ComparitorType]Resolver{
			LessThan: func(a, b interface{}) bool { return a.(float64) < b.(float64) },
		},
	}
)

func EqualsComparitor(k reflect.Kind, a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func NotEqualsComparitor(k reflect.Kind, a, b interface{}) bool {
	return !EqualsComparitor(k, a, b)
}

func LessThanComparitor(k reflect.Kind, a, b interface{}) bool {
	return Comparitors[k][LessThan](a, b)
}

func LessThanOrEqualToComparitor(k reflect.Kind, a, b interface{}) bool {
	return EqualsComparitor(k, a, b) || Comparitors[k][LessThan](a, b)
}

func GreaterThanComparitor(k reflect.Kind, a, b interface{}) bool {
	return !EqualsComparitor(k, a, b) && !Comparitors[k][LessThan](a, b)
}

func GreaterOrEqualToComparitor(k reflect.Kind, a, b interface{}) bool {
	return !Comparitors[k][LessThan](a, b)
}
