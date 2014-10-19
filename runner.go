package expect

import (
	"flag"
	"reflect"
	"regexp"
	"testing"
	"fmt"
	"time"
	"runtime"
	"strings"
	"github.com/wsxiaoys/terminal/color"
	"os"
)

var (
	showStdout = flag.Bool("vv", false, "Regular expression selecting which tests to run")
	matchFlag = flag.String("m", "", "Regular expression selecting which tests to run")
	pattern *regexp.Regexp
	runner *Runner
	stdout = os.Stdout
	silentOut *os.File
)

func init() {
	flag.Parse()
	pattern = regexp.MustCompile("(?i)" + *matchFlag)
	if *showStdout == true {
		silentOut = stdout
	}
	os.Stdout = silentOut
}

func Expectify(suite interface{}, t *testing.T) {
	tp := reflect.TypeOf(suite)
	sv := reflect.ValueOf(suite)
	count := tp.NumMethod()

	runner = &Runner{
		results: make([]*Result, 0, 10),
	}
	announced := false
	for i := 0; i < count; i++ {
		method := tp.Method(i)
		name := method.Name
		if pattern.MatchString(name) == false || method.Type.NumIn() != 1{
			continue
		}
		os.Stdout = stdout
		result := runner.Start(name)
		method.Func.Call([]reflect.Value{sv})
		if runner.End() == false || testing.Verbose() {
			if announced == false {
				color.Printf("\n@!%s@|\n", sv.Elem().Type().String())
				announced = true
			}
			result.Report()
		}
		os.Stdout = silentOut
	}
	if runner.Passed() == false {
		os.Stdout = stdout
		fmt.Println("")
		os.Stdout = silentOut
		t.Fail()
	}
}

type Runner struct {
	results []*Result
	current *Result
}

func (r *Runner) Start(name string) *Result {
	r.current = &Result{
		method: name,
		start: time.Now(),
		failures: make([]*Failure, 0, 3),
	}
	r.results = append(r.results, r.current)
	return r.current
}

func (r *Runner) End() bool {
	r.current.end = time.Now()
	passed := r.current.Passed()
	r.current = nil
	return passed
}

func (r *Runner) Passed() bool {
	for _, result := range r.results {
		if result.Passed() == false {
			return false
		}
	}
	return true
}

func (r *Runner) Skip(format string, args ...interface{}) {
	if r.current != nil {
		r.current.Skip(format, args...)
	}
}

func (r *Runner) Errorf(format string, args ...interface{}) {
	file := "???"
	line := 1
	ok := false
	for i := 3; i < 10; i++ {
		_, file, line, ok = runtime.Caller(i)
		if ok == false || strings.HasSuffix(file, "_test.go") {
			break
		}
	}

	if ok {
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	}

	failure := &Failure{
		message: fmt.Sprintf(format, args...),
		location: fmt.Sprintf("%s:%d", file, line),
	}
	r.current.failures = append(r.current.failures, failure)
}

type Result struct {
	method string
	failures []*Failure
	start time.Time
	end time.Time
	skipMessage string
	skip bool
}

type Failure struct {
	message string
	location string
}

func (r *Result) Skip(format string, args ...interface{}) {
	r.skip = true
	r.skipMessage = fmt.Sprintf(format, args...)
}

func (r *Result) Passed() bool {
	return r.skip || len(r.failures) == 0
}

func (r *Result) Report() {
	info := fmt.Sprintf(" %-50s\t%fs", r.method, r.end.Sub(r.start).Seconds())
	if r.skip {
		color.Println(" @y⸚" + info)
		color.Println("   @." + r.skipMessage)
	} else if r.Passed() {
		color.Println(" @g✓" + info)
	} else {
		color.Println(" @r×" + info)
		for _, failure := range r.failures {
			color.Printf("   @.%-50s\t%-30s\n", failure.message, failure.location)
		}
	}
}
