package expect

import (
	"fmt"
	"testing"
)

var originalErrorf = Errorf

type ExpectTests struct{}

func Test_Expect(t *testing.T) {
	Expectify(new(ExpectTests), t)
}

func (e *ExpectTests) Equal_Success() {
	Expect("abc").To.Equal("abc")
	Expect(true).To.Equal(true)
	Expect(false).To.Equal(false)
	Expect([]int{1, 2, 3}).To.Equal([]int{1, 2, 3})
	Expect([]bool{true, false, false}).To.Equal([]bool{true, false, false})
	Expect(map[string]int{"a": 1, "b": 2}).To.Equal(map[string]int{"a": 1, "b": 2})
}

func (e *ExpectTests) Equal_Failures() {
	failing("expected abc to be equal to 123", 1, func() {
		Expect("abc").To.Equal("123")
	})

	failing("expected true to be equal to false", 1, func() {
		Expect(true).To.Equal(false)
	})

	failing("expected [1 2] to be equal to [1]", 1, func() {
		Expect([]int{1, 2}).To.Equal([]int{1})
	})

	// hard to test, the print out is indeterministic, maps aren't sorted
	// failing("expected map[a:1 b:2] to be equal to map[b:1 a:2]", 1, func(){
	// 	Expect(map[string]int{"a":1,"b":2}).Equal.To(map[string]int{"b":1,"a":2})
	// })
}

func (e *ExpectTests) NotEqual_Success() {
	Expect("abc").Not.To.Equal("123")
}

func (e *ExpectTests) NoEqual_Failures() {
	failing("expected 123 not to equal 123", 1, func() {
		Expect("123").Not.To.Equal("123")
	})

	failing("expected true not to equal true", 1, func() {
		Expect(true).Not.To.Equal(true)
	})

	failing("expected [1 2] not to equal [1 2]", 1, func() {
		Expect([]int{1, 2}).Not.To.Equal([]int{1, 2})
	})
}

func (e *ExpectTests) Contain_Success() {
	Expect("abc").To.Contain("bc")
	Expect("abc").Not.To.Contain("99")

	Expect([]int{1, 3, 4}).To.Contain(4)
	Expect([]int{2, 3, 4}).Not.To.Contain(5)

	Expect(map[string]int{"a": 3}).To.Contain("a")
	Expect(map[string]int{"a": 3}).Not.To.Contain("b")

	Expect([]byte{1, 2, 3}).To.Contain(byte(1))
	Expect([]byte{1, 2, 3}).Not.To.Contain(byte(5))
	Expect([]byte{1, 2, 3}).Not.To.Contain([]byte{2, 3, 4})
}

func (e *ExpectTests) Contain_Failures() {
	failing("abc does not contain zz", 1, func() {
		Expect("abc").To.Contain("zz")
	})
	failing("aazzbb contains zz", 1, func() {
		Expect("aazzbb").Not.To.Contain("zz")
	})

	failing("[1 2] does not contain 3", 1, func() {
		Expect([]int{1,2}).To.Contain(3)
	})
	failing("[1 2] contains 2", 1, func() {
		Expect([]int{1,2}).Not.To.Contain(2)
	})

	failing("map[a:3] does not contain b", 1, func() {
		Expect(map[string]int{"a": 3}).To.Contain("b")
	})
	failing("map[a:3] contains a", 1, func() {
		Expect(map[string]int{"a": 3}).Not.To.Contain("a")
	})

	failing("[1 2 3] does not contain [3 4]", 1, func() {
		Expect([]byte{1, 2, 3}).To.Contain([]byte{3,4})
	})
	failing("[1 2 3] contains [2 3]", 1, func() {
		Expect([]byte{1, 2, 3}).Not.To.Contain([]byte{2,3})
	})
}

func (e *ExpectTests) ExpectMulipleValues_Success() {
	Expect(1, true, "over 9000").To.Equal(1, true, "over 9000")
}

func (e *ExpectTests) ExpectMulipleValues_Failure() {
	failing("mismatch number of values and expectations 3 != 2", 1, func() {
		Expect(1, true, "over 9000").To.Equal(1, true)
	})
	failing("mismatch number of values and expectations 2 != 3", 1, func() {
		Expect(1, true).To.Equal(1, true, "a")
	})
	failing("expected 1 to be equal to 2", 1, func() {
		Expect(1, "a").To.Equal(2, "a")
	})
	failing("expected a to be equal to b", 1, func() {
		Expect(1, "a").To.Equal(1, "b")
	})
}

func (e *ExpectTests) ExpectNil_Success() {
	var t *ExpectTests
	Expect(nil).To.Equal(nil)
	Expect(t).To.Equal(nil)
	Expect(t).To.Equal(t)
	Expect(nil).Not.To.Equal(32)
	Expect(32).Not.To.Equal(t)
}

func (e *ExpectTests) ExpectNil_Failure() {
	var t *ExpectTests
	failing("expected <nil> to be equal to 44", 2, func(){
		Expect(nil).To.Equal(44)
		Expect(t).To.Equal(44)
	})
	failing("expected <nil> not to equal <nil>", 3, func(){
		Expect(nil).Not.To.Equal(nil)
		Expect(nil).Not.To.Equal(t)
		Expect(t).Not.To.Equal(t)
	})
}

func (e *ExpectTests) Fail() {
	failing("the failure reason 123", 1, func() {
		Fail("the failure reason %d", 123)
	})
}

func (e *ExpectTests) Skip_After() {
	Fail("failed")
	Skip("the skip reason %d", 9001)
}

func (e *ExpectTests) Skip_Before() {
	Skip("the skip reason")
	Fail("failed")
}

func failing(expected string, count int, f func()) {
	actuals := make([]string, 0, 5)
	Errorf = func(format string, args ...interface{}) {
		actuals = append(actuals, fmt.Sprintf(format, args...))
	}
	f()
	Errorf = originalErrorf
	if len(actuals) != count {
		Errorf("expected %d failures got %d", count, len(actuals))
		return
	}
	for i := 0; i < count; i++ {
		if actuals[i] != expected {
			if count == 1 {
				Errorf("expected failure '%s' got '%s'", expected, actuals[i])
			} else {
				Errorf("expected failure '%s' got '%s' index %d", expected, actuals[i], i+1)
			}
		}
	}
}
