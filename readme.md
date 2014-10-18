# Expect

A library to help you write tests in Go. Wraps Go's testing framework; your workflow stays the same.

## Example

```go
import (
  "testing"
  . "github.com/karlseguin/expect"
)

type CalculatorTests struct{}

func Test_Caculator(t *testing.T) {
  Expectify(new(CalculatorTests), t)
}

func (c *CalculatorTests) AddsTwoNumbers() {
  c := new(Calculator)
  Expect(c.Add(4, 8)).To.Equal(12)
  Expect(c.Add(10, 2)).Greater.Than(11)
  Expect(c.Add(10, 2)).Not.Greater.Than(9000)
  Expect(c.Add(1, 1)).Not.To.Equal(3)
}
```

## Running

Run tests as you normally would via `go run test`. However, to run specific tests, use the -m flag, which will do a case-insensitive regular expression match.

    go run test -m AddsTwo

## Expectations

* `Greater.Than(x)`
* `GreaterOrEqual.To(x)`
* `Less.Than(x)`
* `LessOrEqual.To(x)`
* `To.Equal(x)`

All expectations can be reversed by starting the chain with `Not.`

### Contains

`To.Contain` works with strings, arrays, slices and maps. For arrays and slices, only individual values are matched. For example:

    Expect([]int{1,2,3}).To.Contain([]int{1,2})

will, sadly, not work.

The exception to this is for strings and `[]byte`. These work with either a single value or an array (they use the stdlib's `strings.Contains` and `bytes.Contains`).


## Multiple Values

`Expect` throws away all but the first value. This is convenient when a method returns an error which you don't care to test:

    Expect(ioutil.ReadFile("blah")).To.Equal([]byte{1, 2, 3, 4})

However, using `To.Equal` multiple values can be provided:

    Expect(1, true, "a").To.Equal(1, true, "a")

## Optional *testing.T

Your tests can optionally accept a `*testing.T` parameter:

```go
func (c *CalculatorTests) AddsTwoNumbers(t *testing.T) {
  ...
}
```
