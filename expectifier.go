package expect

import (
	"flag"
	"reflect"
	"regexp"
	"testing"
)

var (
	matchString = flag.String("m", "", "Regular expression selecting which tests and/or suites to run")
	testType    = reflect.TypeOf(new(testing.T))
	t           *testing.T
)

func Expectify(suite interface{}, tt *testing.T) {
	t = tt
	tv := reflect.ValueOf(tt)
	tp := reflect.TypeOf(suite)
	sv := reflect.ValueOf(suite)
	count := tp.NumMethod()
	tests := make([]testing.InternalTest, 0, count)
	for i := 0; i < count; i++ {
		method := tp.Method(i)
		count := method.Type.NumIn()
		name := "\t" + method.Name
		var f func(t *testing.T)

		if count == 1 {
			f = func(t *testing.T) {
				method.Func.Call([]reflect.Value{sv})
			}
		} else if count == 2 && method.Type.In(1) == testType {
			f = func(t *testing.T) {
				method.Func.Call([]reflect.Value{sv, tv})
			}
		}
		if f != nil {
			tests = append(tests, testing.InternalTest{name, f})
		}
	}

	var pattern *regexp.Regexp
	if len(*matchString) != 0 {
		pattern = regexp.MustCompile("(?i)" + *matchString)
	}
	matcher := func(_, str string) (bool, error) {
		if pattern == nil {
			return true, nil
		}
		return pattern.MatchString(str), nil
	}
	passed := testing.RunTests(matcher, tests)
	if passed == false {
		t.Fail()
	}
}
