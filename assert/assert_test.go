package assert_test

import (
	"errors"
	"fmt"
	"iter"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

type tester struct {
	*testing.T
	interfered bool
}

func (t *tester) Fatal(args ...any) {
	t.interfered = true
}

func (t *tester) Fatalf(message string, args ...any) {
	t.interfered = true
}

var errFoo = errors.New("error")

func assertSuccess(t *testing.T, name string, fn func(t testing.TB)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()

		tt := &tester{T: t}
		fn(tt)
		if tt.interfered {
			t.Fatal("shouldn't fail")
		}
	})
}

func assertFailure(t *testing.T, name string, fn func(t testing.TB)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()

		tt := &tester{T: t}
		fn(tt)
		if !tt.interfered {
			t.Fatal("should fail")
		}
	})
}

type stubStruct[T any] struct {
	V T
}

func TestNil(t *testing.T) {
	var ptr *int
	var err error
	var fun func(v int)
	var gen *stubStruct[int]
	assertSuccess(t, "nil", func(t testing.TB) {
		assert.Nil(t, nil)
	})
	assertSuccess(t, "nil pointer", func(t testing.TB) {
		assert.Nil(t, ptr)
	})
	assertSuccess(t, "nil error", func(t testing.TB) {
		assert.Nil(t, err)
	})
	assertSuccess(t, "nil func", func(t testing.TB) {
		assert.Nil(t, gen)
	})
	ptr = new(int)
	err = errFoo
	fun = func(v int) {}
	assertFailure(t, "not nil pointer", func(t testing.TB) {
		assert.Nil(t, ptr)
	})
	assertFailure(t, "not nil error", func(t testing.TB) {
		assert.Nil(t, err)
	})
	assertFailure(t, "not nil func", func(t testing.TB) {
		assert.Nil(t, fun)
	})
}

func TestNotNil(t *testing.T) {
	var ptr *int
	var err error
	var fun func(v int)
	assertFailure(t, "nil pointer", func(t testing.TB) {
		assert.NotNil(t, ptr)
	})
	assertFailure(t, "nil error", func(t testing.TB) {
		assert.NotNil(t, err)
	})
	assertFailure(t, "nil func", func(t testing.TB) {
		assert.NotNil(t, fun)
	})
	ptr = new(int)
	err = errFoo
	fun = func(v int) {}
	assertSuccess(t, "not nil pointer", func(t testing.TB) {
		assert.NotNil(t, ptr)
	})
	assertSuccess(t, "not nil error", func(t testing.TB) {
		assert.NotNil(t, err)
	})
}

func TestZero(t *testing.T) {
	assertSuccess(t, "zero int", func(t testing.TB) {
		assert.Zero(t, 0)
	})
	assertSuccess(t, "zero string", func(t testing.TB) {
		assert.Zero(t, "")
	})
	assertSuccess(t, "zero struct", func(t testing.TB) {
		assert.Zero(t, struct{ v any }{})
	})
	assertFailure(t, "not zero int", func(t testing.TB) {
		assert.Zero(t, 1)
	})
	assertFailure(t, "not zero string", func(t testing.TB) {
		assert.Zero(t, "foo")
	})
	assertFailure(t, "not zero struct", func(t testing.TB) {
		assert.Zero(t, struct{ v any }{"a"})
	})
}

func TestNotZero(t *testing.T) {
	assertSuccess(t, "not zero int", func(t testing.TB) {
		assert.NotZero(t, 1)
	})
	assertSuccess(t, "not zero string", func(t testing.TB) {
		assert.NotZero(t, "foo")
	})
	assertSuccess(t, "not zero struct", func(t testing.TB) {
		assert.NotZero(t, struct{ v any }{"a"})
	})
	assertFailure(t, "zero int", func(t testing.TB) {
		assert.NotZero(t, 0)
	})
	assertFailure(t, "zero string", func(t testing.TB) {
		assert.NotZero(t, "")
	})
	assertFailure(t, "zero struct", func(t testing.TB) {
		assert.NotZero(t, struct{ v any }{})
	})
}

func TestEmpty(t *testing.T) {
	assertSuccess(t, "empty string", func(t testing.TB) {
		assert.Empty(t, "")
	})
	assertSuccess(t, "empty slice", func(t testing.TB) {
		assert.Empty(t, []any{})
	})
	assertSuccess(t, "empty map", func(t testing.TB) {
		assert.Empty(t, map[any]any{})
	})
	assertSuccess(t, "empty chan", func(t testing.TB) {
		assert.Empty(t, make(chan int, 1))
	})
	assertFailure(t, "not empty string", func(t testing.TB) {
		assert.Empty(t, "foo")
	})
	assertFailure(t, "not empty slice", func(t testing.TB) {
		assert.Empty(t, []any{1})
	})
	assertFailure(t, "not empty map", func(t testing.TB) {
		assert.Empty(t, map[any]any{1: 1})
	})
	assertFailure(t, "not empty chan", func(t testing.TB) {
		ch := make(chan int, 1)
		ch <- 1
		assert.Empty(t, ch)
	})
}

func TestNotEmpty(t *testing.T) {
	assertSuccess(t, "not empty string", func(t testing.TB) {
		assert.NotEmpty(t, "foo")
	})
	assertSuccess(t, "not empty slice", func(t testing.TB) {
		assert.NotEmpty(t, []any{1})
	})
	assertSuccess(t, "not empty map", func(t testing.TB) {
		assert.NotEmpty(t, map[any]any{1: 1})
	})
	assertSuccess(t, "not empty chan", func(t testing.TB) {
		ch := make(chan int, 1)
		ch <- 1
		assert.NotEmpty(t, ch)
	})
	assertFailure(t, "empty string", func(t testing.TB) {
		assert.NotEmpty(t, "")
	})
	assertFailure(t, "empty slice", func(t testing.TB) {
		assert.NotEmpty(t, []any{})
	})
	assertFailure(t, "empty map", func(t testing.TB) {
		assert.NotEmpty(t, map[any]any{})
	})
	assertFailure(t, "empty chan", func(t testing.TB) {
		assert.NotEmpty(t, make(chan int, 1))
	})
}

func TestTrue(t *testing.T) {
	assertSuccess(t, "true", func(t testing.TB) {
		assert.True(t, true)
	})
	assertFailure(t, "not true", func(t testing.TB) {
		assert.True(t, false)
	})
}

func TestFalse(t *testing.T) {
	assertSuccess(t, "false", func(t testing.TB) {
		assert.False(t, false)
	})
	assertFailure(t, "not false", func(t testing.TB) {
		assert.False(t, true)
	})
}

type pointer struct {
	X, Y float64
}

func TestEqual(t *testing.T) {
	assertSuccess(t, "equal numbers", func(t testing.TB) {
		assert.Equal(t, 1, 1)
	})
	assertSuccess(t, "equal strings", func(t testing.TB) {
		assert.Equal(t, "foo", "foo")
	})
	assertSuccess(t, "equal bytes slice", func(t testing.TB) {
		assert.Equal(t, []byte{'a', 'b'}, []byte{'a', 'b'})
	})
	assertSuccess(t, "equal empty bytes slice", func(t testing.TB) {
		assert.Equal(t, []byte(""), []byte{})
	})
	assertSuccess(t, "equal structs", func(t testing.TB) {
		assert.Equal(t, pointer{0.55, 0.75}, pointer{0.55, 0.75})
	})
	assertFailure(t, "nonequal numbers", func(t testing.TB) {
		assert.Equal(t, 1, 2)
	})
	assertFailure(t, "nonequal strings", func(t testing.TB) {
		assert.Equal(t, "foo", "bar")
	})
	assertFailure(t, "nonequal bytes slice", func(t testing.TB) {
		assert.Equal(t, []byte{'a', 'b'}, []byte{'a', 'c'})
	})
	assertFailure(t, "nonequal structs", func(t testing.TB) {
		assert.Equal(t, pointer{0.55, 0.75}, pointer{1, 1})
	})
}

func TestNotEqual(t *testing.T) {
	assertSuccess(t, "nonequal numbers", func(t testing.TB) {
		assert.NotEqual(t, 1, 2)
	})
	assertSuccess(t, "nonequal strings", func(t testing.TB) {
		assert.NotEqual(t, "foo", "bar")
	})
	assertSuccess(t, "nonequal bytes slice", func(t testing.TB) {
		assert.NotEqual(t, []byte{'a', 'b'}, []byte{'a', 'c'})
	})
	assertSuccess(t, "nonequal structs", func(t testing.TB) {
		assert.NotEqual(t, pointer{0.55, 0.75}, pointer{1, 1})
	})
	assertFailure(t, "equal numbers", func(t testing.TB) {
		assert.NotEqual(t, 1, 1)
	})
	assertFailure(t, "equal strings", func(t testing.TB) {
		assert.NotEqual(t, "foo", "foo")
	})
	assertFailure(t, "equal bytes slice", func(t testing.TB) {
		assert.NotEqual(t, []byte{'a', 'b'}, []byte{'a', 'b'})
	})
	assertFailure(t, "equal empty bytes slice", func(t testing.TB) {
		assert.NotEqual(t, []byte(""), []byte{})
	})
	assertFailure(t, "equal structs", func(t testing.TB) {
		assert.NotEqual(t, pointer{0.55, 0.75}, pointer{0.55, 0.75})
	})
}

func TestEqualFunc(t *testing.T) {
	cmpFn := func(a, b int) bool {
		return a == b
	}
	assertSuccess(t, "equal", func(t testing.TB) {
		assert.EqualFunc(t, 1, 1, cmpFn)
	})
	assertFailure(t, "nonequal", func(t testing.TB) {
		assert.EqualFunc(t, 1, 0, cmpFn)
	})
}

func TestNotEqualFunc(t *testing.T) {
	cmpFn := func(a, b int) bool {
		return a == b
	}
	assertSuccess(t, "nonequal", func(t testing.TB) {
		assert.NotEqualFunc(t, 1, 0, cmpFn)
	})
	assertFailure(t, "equal", func(t testing.TB) {
		assert.NotEqualFunc(t, 1, 1, cmpFn)
	})
}

func TestError(t *testing.T) {
	var err error = errors.New("error")

	assertSuccess(t, "same error", func(t testing.TB) {
		assert.Error(t, err, err)
	})
	assertSuccess(t, "both nil errors", func(t testing.TB) {
		assert.Error(t, nil, nil)
	})
	assertFailure(t, "different errors", func(t testing.TB) {
		assert.Error(t, err, errors.New("other error"))
	})
	assertFailure(t, "error and nil error", func(t testing.TB) {
		assert.Error(t, err, nil)
	})
}

func TestNotError(t *testing.T) {
	var err error = errors.New("error")

	assertSuccess(t, "different error", func(t testing.TB) {
		assert.NotError(t, err, errors.New("other error"))
	})
	assertSuccess(t, "error and nil error", func(t testing.TB) {
		assert.NotError(t, err, nil)
	})
	assertFailure(t, "same error", func(t testing.TB) {
		assert.NotError(t, err, err)
	})
	assertFailure(t, "both nil errors", func(t testing.TB) {
		assert.NotError(t, nil, nil)
	})
}

func TestContains(t *testing.T) {
	assertSuccess(t, "array containing the element", func(t testing.TB) {
		assert.Contains(t, []int{1, 2, 3, 4, 5}, 3)
	})
	assertFailure(t, "array not containing the element", func(t testing.TB) {
		assert.Contains(t, []int{1, 2, 3, 4, 5}, 0)
	})
	assertSuccess(t, "string containing the substring", func(t testing.TB) {
		assert.Contains(t, "abcdef", "cd")
	})
	assertFailure(t, "string not containing the substring", func(t testing.TB) {
		assert.Contains(t, "abcdef", "ce")
	})
	iter := func() iter.Seq[stubStruct[int]] {
		return func(yield func(stubStruct[int]) bool) {
			for i := range 10 {
				if !yield(stubStruct[int]{i}) {
					return
				}
			}
		}
	}
	assertSuccess(t, "iterable containing the element", func(t testing.TB) {
		assert.Contains(t, iter(), stubStruct[int]{5})
	})
	assertFailure(t, "iterable not containing the element", func(t testing.TB) {
		assert.Contains(t, iter(), stubStruct[int]{-1})
	})
}

func TestNotContains(t *testing.T) {
	assertSuccess(t, "array not containing the element", func(t testing.TB) {
		assert.NotContains(t, []int{1, 2, 3, 4, 5}, 0)
	})
	assertFailure(t, "array containing the element", func(t testing.TB) {
		assert.NotContains(t, []int{1, 2, 3, 4, 5}, 3)
	})
	assertSuccess(t, "string not containing the substring", func(t testing.TB) {
		assert.NotContains(t, "abcdef", "ce")
	})
	assertFailure(t, "string containing the substring", func(t testing.TB) {
		assert.NotContains(t, "abcdef", "cd")
	})
	iter := func() iter.Seq[stubStruct[int]] {
		return func(yield func(stubStruct[int]) bool) {
			for i := range 10 {
				if !yield(stubStruct[int]{i}) {
					return
				}
			}
		}
	}
	assertSuccess(t, "iterable not containing the element", func(t testing.TB) {
		assert.NotContains(t, iter(), stubStruct[int]{-1})
	})
	assertFailure(t, "iterable containing the element", func(t testing.TB) {
		assert.NotContains(t, iter(), stubStruct[int]{5})
	})
}

func TestContainsFunc(t *testing.T) {
	assertSuccess(t, "array contains the element", func(t testing.TB) {
		assert.ContainsFunc(t, []int{1, 2, 3, 4, 5}, func(e int) bool {
			return e == 3
		})
	})
	assertFailure(t, "array doesn't contain the element", func(t testing.TB) {
		assert.ContainsFunc(t, []int{1, 2, 3, 4, 5}, func(e int) bool {
			return e == 0
		})
	})
	assertSuccess(t, "array of pointers contains pointer with specific x-axis value", func(t testing.TB) {
		assert.ContainsFunc(t, []pointer{{0.44, 0.22}, {0.10, 0.25}}, func(e pointer) bool {
			return e.X == 0.10
		})
	})
	assertFailure(t, "array of pointers doesn't contain pointer with specific x-axis value", func(t testing.TB) {
		assert.ContainsFunc(t, []pointer{{0.44, 0.22}, {0.10, 0.25}}, func(e pointer) bool {
			return e.X == 0.00
		})
	})
}

func TestNotContainsFunc(t *testing.T) {
	assertSuccess(t, "array doesn't contain the element", func(t testing.TB) {
		assert.NotContainsFunc(t, []int{1, 2, 3, 4, 5}, func(e int) bool {
			return e == 0
		})
	})
	assertFailure(t, "array contain the element", func(t testing.TB) {
		assert.NotContainsFunc(t, []int{1, 2, 3, 4, 5}, func(e int) bool {
			return e == 3
		})
	})
	assertSuccess(t, "array of pointers doesn't contain pointer with specific x-axis value", func(t testing.TB) {
		assert.NotContainsFunc(t, []pointer{{0.44, 0.22}, {0.10, 0.25}}, func(e pointer) bool {
			return e.X == 0.00
		})
	})
	assertFailure(t, "array of pointers contains pointer with specific x-axis value", func(t testing.TB) {
		assert.NotContainsFunc(t, []pointer{{0.44, 0.22}, {0.10, 0.25}}, func(e pointer) bool {
			return e.X == 0.10
		})
	})
}

func TestHasPrefix(t *testing.T) {
	assertSuccess(t, "string starts with the given string", func(t testing.TB) {
		assert.HasPrefix(t, "nice to meet you", "nice")
	})
	assertFailure(t, "string doesn't start with the given string", func(t testing.TB) {
		assert.HasPrefix(t, "nice to meet you", "you")
	})
}

func TestHasNoPrefix(t *testing.T) {
	assertSuccess(t, "string doesn't start with the given string", func(t testing.TB) {
		assert.HasNoPrefix(t, "nice to meet you", "you")
	})
	assertFailure(t, "string starts with the given string", func(t testing.TB) {
		assert.HasNoPrefix(t, "nice to meet you", "nice")
	})
}

func TestHasSuffix(t *testing.T) {
	assertSuccess(t, "string starts with the given string", func(t testing.TB) {
		assert.HasSuffix(t, "nice to meet you", "you")
	})
	assertFailure(t, "string doesn't start with the given string", func(t testing.TB) {
		assert.HasSuffix(t, "nice to meet you", "meet")
	})
}

func TestHasNoSuffix(t *testing.T) {
	assertSuccess(t, "string doesn't ends with the given string", func(t testing.TB) {
		assert.HasNoSuffix(t, "nice to meet you", "nice")
	})
	assertFailure(t, "string ends with the given string", func(t testing.TB) {
		assert.HasNoSuffix(t, "nice to meet you", "you")
	})
}

type writer struct {
	b []byte
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

func TestPanics(t *testing.T) {
	assertSuccess(t, "panics", func(t testing.TB) {
		assert.Panics(t, func() {
			panic("panic")
		})
	})
	assertFailure(t, "does not panic", func(t testing.TB) {
		assert.Panics(t, func() {
			fmt.Fprint(&writer{[]byte{}}, "yo!")
		})
	})
}

func TestNotPanics(t *testing.T) {
	assertSuccess(t, "does not panic", func(t testing.TB) {
		assert.NotPanics(t, func() {
			fmt.Fprint(&writer{[]byte{}}, "yo!")
		})
	})
	assertFailure(t, "panics", func(t testing.TB) {
		assert.NotPanics(t, func() {
			panic("panic")
		})
	})
}

func TestPanicIs(t *testing.T) {
	assertSuccess(t, "panics the expected panic", func(t testing.TB) {
		assert.PanicIs(t, func() {
			panic("panic")
		}, "panic")
	})
	assertFailure(t, "does not panic the expected panic", func(t testing.TB) {
		assert.PanicIs(t, func() {
			panic("panic")
		}, "other panic")
	})
}
