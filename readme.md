# Expect

A library to help you write tests in Go.

What's wrong with Go's built-in testing package? Not much, except it tends to lead to verbose code. `Expect` runs within `go test` but provides a different syntax for specifying expectations.

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
  Expect(c.Add(10, 2)).LessOrEqual.To(12)
  Expect(c.Add(10, 2)).Not.Greater.Than(9000)
  Expect(c.Add(1, 1)).Not.To.Equal(3)

  //OR

  Expect(c.Add(4.8)).ToEqual(12)
  Expect(c.Add(10, 2)).GreaterThan(11)
  Expect(c.Add(10, 2)).LessOrEqualTo(12)
  Expect(c.Add(10, 2)).Not.GreaterThan(9000)
  Expect(c.Add(1, 1)).Not.ToEqual(3)
}
```

You can also use `Skip(format, args...)` and `Fail(format, args...)` to either skip a test, or cause a test to fail.

## Running

Run tests as you normally would via `go run test`. However, to run specific tests, use the -m flag, which will do a case-insensitive regular expression match.

    go run test -m AddsTwo

## Expectations

Two similar syntaxes of expectations are supported

* `To.Equal(x)` or `ToEqual(x)`
* `Greater.Than(x)` or `GreaterThan(x)`
* `GreaterOrEqual.To(x)` or `GreaterOrEqualTo(x)`
* `Less.Than(x)` or `LessThan(x)`
* `LessOrEqual.To(x)` or `LessOrEqualTo(x)`

All expectations can be reversed by starting the chain with `Not.`:

* `Not.To.Equal(x)` or `Not.ToEqual(x)`
* `Not.Greater.Than(x)` or `Not.GreaterThan(x)`
* `Not.GreaterOrEqual.To(x)` or `Not.GreaterOrEqualTo(x)`
* `Not.Less.Than(x)` or `Not.LessThan(x)`
* `Not.LessOrEqual.To(x)` or `Not.LessOrEqualTo(x)`

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

## stdout

Go's testing package has no hooks into its reporting. `Expect` takes the drastic step of occasionally silencing `os.Stdout`, which many packages uses (such as `fmt.Print`). However, within your test, `os.Stdout` **will** work.

If you print anything outside of your test, say during `init`, it'll likely be silenced by `Expect`. You can disable this behavior with the `-vv` flag (use `-v` and `-vv` in combination to change the behavior of both `Expect` and Go's library)

## Mixing with *testing.T

Since `Expect` runs within `go test`, you can mix `Expect` style tests with traditional Go tests. To do this, you'll probably want to run your tests with `-vv` (see *stdout* section above). Even without `-vv`, you *will* get the proper return code (!= 0 on failure).

This also means that `go test` features, such as `-cover`, work with `Expect`. However, `Expect` tests cannot be run in parallel.

# Mocks

`Expect` also exposes a sub-package `mock` which is aimed at providing mock objects and buiders for common standard library components. This is a work in progress
