package assert

import (
	"cmp"
	"fmt"
	"iter"
	"reflect"
	"strings"
	"testing"
)

func output(t testing.TB, common string, out []any) {
	t.Helper()

	if len(out) == 0 {
		t.Fatal(common)
		return
	}
	if str, ok := out[0].(string); ok {
		t.Fatalf(str, out[1:]...)
		return
	}
	panic("assert: output argument must be a format string and its args")
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Pointer, reflect.UnsafePointer,
		reflect.Map, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

func Nil(t testing.TB, v any, out ...any) {
	t.Helper()
	if !isNil(v) {
		common := fmt.Sprintf("expected nil value, got %v", v)
		output(t, common, out)
	}
}

func NotNil(t testing.TB, v any, out ...any) {
	t.Helper()
	if isNil(v) {
		output(t, "expected not nil value", out)
	}
}

func Zero[T any](t testing.TB, v T, out ...any) {
	t.Helper()

	if !reflect.ValueOf(v).IsZero() {
		common := fmt.Sprintf("expected zero value for the type %T", v)
		output(t, common, out)
	}
}

func NotZero[T any](t testing.TB, v T, out ...any) {
	t.Helper()

	if reflect.ValueOf(v).IsZero() {
		common := fmt.Sprintf("expected non-zero value for the type %T", v)
		output(t, common, out)
	}
}

func isEmpty(v any) bool {
	if isNil(v) {
		return true
	}

	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Array,
		reflect.Chan,
		reflect.Map,
		reflect.Slice,
		reflect.String:
		return value.Len() == 0
	case reflect.Ptr, reflect.UnsafePointer:
		return isEmpty(value.Elem().Interface())
	default:
		zero := reflect.Zero(value.Type())
		return reflect.DeepEqual(value, zero)
	}
}

func Empty[T any](t testing.TB, v T, out ...any) {
	t.Helper()

	if !isEmpty(v) {
		common := fmt.Sprintf("expected empty, but got %v", v)
		output(t, common, out)
	}
}

func NotEmpty[T any](t testing.TB, v T, out ...any) {
	t.Helper()

	if isEmpty(v) {
		common := fmt.Sprintf("expected not empty, but got %v", v)
		output(t, common, out)
	}
}

func True(t testing.TB, got bool, out ...any) {
	t.Helper()

	if !got {
		common := "didn't get true"
		output(t, common, out)
	}
}

func False(t testing.TB, got bool, out ...any) {
	t.Helper()

	if got {
		common := "didn't get false"
		output(t, common, out)
	}
}

func isEqual[T any](a, b T) bool {
	valA := reflect.ValueOf(a)
	if valA.Kind() == reflect.Func {
		valB := reflect.ValueOf(b)
		return valB.Kind() == reflect.Func && valA.Pointer() == valB.Pointer()
	}
	return reflect.DeepEqual(a, b)
}

func Equal[T any](t testing.TB, a, b T, out ...any) {
	t.Helper()

	if !isEqual(a, b) {
		common := fmt.Sprintf("expected equal values, but got %v and %v", a, b)
		output(t, common, out)
	}
}

func NotEqual[T any](t testing.TB, a, b T, out ...any) {
	t.Helper()

	if isEqual(a, b) {
		common := fmt.Sprintf("expected different values, but %v is equal to %v ", a, b)
		output(t, common, out)
	}
}

func EqualFunc[T any](t testing.TB, a, b T, cmp func(T, T) bool, out ...any) {
	t.Helper()

	if !cmp(a, b) {
		common := fmt.Sprintf("%v and %v can't be the same accordingly to comparator", a, b)
		output(t, common, out)
	}
}

func NotEqualFunc[T any](t testing.TB, a, b T, cmp func(T, T) bool, out ...any) {
	t.Helper()

	if cmp(a, b) {
		common := fmt.Sprintf("%v and %v are the same accordingly to comparator", a, b)
		output(t, common, out)
	}
}

func Error(t testing.TB, got, want error, out ...any) {
	t.Helper()

	if got != want {
		common := fmt.Sprintf("expected error %v, but got %v", want, got)
		output(t, common, out)
	}
}

func NotError(t testing.TB, got, nwant error, out ...any) {
	t.Helper()

	if got == nwant {
		common := fmt.Sprintf("didn't expected error %v, but got it", nwant)
		output(t, common, out)
	}
}

type StringOrSet[T comparable] interface {
	string | Set[T]
}

func contains[T comparable, S StringOrSet[T]](s S, lf T) bool {
	valS := reflect.ValueOf(s)
	var collec iter.Seq[T]
	switch valS.Kind() {
	case reflect.String:
		return strings.Contains(valS.String(), reflect.ValueOf(lf).String())
	case reflect.Func:
		collec = valS.Interface().(iter.Seq[T])
	default:
		collec = func(yield func(T) bool) {
			for i := range valS.Len() {
				if !yield(valS.Index(i).Interface().(T)) {
					return
				}
			}
		}
	}
	for v := range collec {
		if v == lf {
			return true
		}
	}
	return false
}

func Contains[T comparable, S StringOrSet[T]](t testing.TB, s S, lf T, out ...any) {
	t.Helper()

	if !contains(s, lf) {
		common := fmt.Sprintf("%v isn't in the collection", lf)
		output(t, common, out)
	}
}

func NotContains[T comparable, S StringOrSet[T]](t testing.TB, s S, lf T, out ...any) {
	t.Helper()

	if contains(s, lf) {
		common := fmt.Sprintf("%v is in the collection", lf)
		output(t, common, out)
	}
}

type Set[T any] interface {
	[]T | iter.Seq[T]
}

func containsFunc[T any, S Set[T]](s S, cmp func(v T) bool) bool {
	var collec iter.Seq[T]
	valS := reflect.ValueOf(s)
	if valS.Kind() == reflect.Func {
		collec = valS.Interface().(iter.Seq[T])
	} else {
		collec = func(yield func(T) bool) {
			for i := range valS.Len() {
				if !yield(valS.Index(i).Interface().(T)) {
					return
				}
			}
		}
	}
	for v := range collec {
		if cmp(v) {
			return true
		}
	}
	return false
}

func ContainsFunc[T any, S Set[T]](t testing.TB, s S, cmp func(T) bool, out ...any) {
	t.Helper()

	if !containsFunc(s, cmp) {
		output(t, "there's no correspondence accordingly to comparator", out)
	}
}

func NotContainsFunc[T any, S Set[T]](t testing.TB, s S, cmp func(T) bool, out ...any) {
	t.Helper()

	if containsFunc(s, cmp) {
		output(t, "there's correspondence accordingly to comparator", out)
	}
}

func HasPrefix(t testing.TB, s string, pfx string, out ...any) {
	t.Helper()

	if !strings.HasPrefix(s, pfx) {
		common := fmt.Sprintf("the %q cannot be a prefix of the %q", pfx, s)
		output(t, common, out)
	}
}

func HasNoPrefix(t testing.TB, s string, pfx string, out ...any) {
	t.Helper()

	if strings.HasPrefix(s, pfx) {
		common := fmt.Sprintf("the %q can be a prefix of the %q", pfx, s)
		output(t, common, out)
	}
}

func HasSuffix(t testing.TB, s string, sfx string, out ...any) {
	t.Helper()

	if !strings.HasSuffix(s, sfx) {
		common := fmt.Sprintf("the %q cannot be a suffix of the %q", sfx, s)
		output(t, common, out)
	}
}

func HasNoSuffix(t testing.TB, s string, sfx string, out ...any) {
	t.Helper()

	if strings.HasSuffix(s, sfx) {
		common := fmt.Sprintf("the %q can be a suffix of the %q", sfx, s)
		output(t, common, out)
	}
}

func Panics(t testing.TB, fn func(), out ...any) {
	t.Helper()

	defer func() {
		t.Helper()
		if recover() == nil {
			output(t, "didn't panic", out)
		}
	}()

	fn()

}

func NotPanics(t testing.TB, fn func(), out ...any) {
	t.Helper()

	defer func() {
		t.Helper()
		if r := recover(); r != nil {
			common := fmt.Sprintf("did panic, %v", r)
			output(t, common, out)
		}
	}()

	fn()
}

func PanicIs(t testing.TB, fn func(), exp any, out ...any) {
	t.Helper()

	defer func() {
		t.Helper()
		if r := recover(); r != exp {
			common := fmt.Sprintf("can't get the expected panic, got %v", r)
			output(t, common, out)
		}
	}()

	fn()
}

func Greater[T cmp.Ordered](t testing.TB, a, b T) {
	t.Helper()

	if a <= b {
		t.Fatalf("%v is actually smaller than %v", a, b)
	}
}

func Smaller[T cmp.Ordered](t testing.TB, a, b T) {
	t.Helper()

	if a >= b {
		t.Fatalf("%v is actually greater than %v", a, b)
	}
}
