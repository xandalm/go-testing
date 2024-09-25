package assert

import (
	"fmt"
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
	panic("assert: output argument must be a fmt string")
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

func Equal[T any](t testing.TB, a, b T, out ...any) {
	t.Helper()

	if !reflect.DeepEqual(a, b) {
		common := fmt.Sprintf("want equal values, got %v and want %v", a, b)
		output(t, common, out)
	}
}

func NotEqual[T any](t testing.TB, a, b T, out ...any) {
	t.Helper()

	if reflect.DeepEqual(a, b) {
		common := fmt.Sprintf("don't want equal values, but %v = %v", a, b)
		output(t, common, out)
	}
}

func EqualFunc[T any](t testing.TB, a, b T, cmp func(T, T) bool, out ...any) {
	t.Helper()

	if !cmp(a, b) {
		common := fmt.Sprintf("want similar values, but %v and %v are not", a, b)
		output(t, common, out)
	}
}

func NotEqualFunc[T any](t testing.TB, a, b T, cmp func(T, T) bool, out ...any) {
	t.Helper()

	if cmp(a, b) {
		common := fmt.Sprintf("don't want similar values, but %v and %v are", a, b)
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

func Contains[T any, S []T](t testing.TB, s S, v T, out ...any) {
	t.Helper()

	for _, e := range s {
		if reflect.DeepEqual(e, v) {
			return
		}
	}
	common := fmt.Sprintf("there's no %v in %v", v, s)
	output(t, common, out)
}

func NotContains[T any, S []T](t testing.TB, s S, v T, out ...any) {
	t.Helper()

	for _, e := range s {
		if reflect.DeepEqual(e, v) {
			common := fmt.Sprintf("there's %v in %v", v, s)
			output(t, common, out)
			return
		}
	}
}

func ContainsFunc[A any, B any, S []A](t testing.TB, s S, v B, fn func(A, B) bool, out ...any) {
	t.Helper()

	for _, e := range s {
		if fn(e, v) {
			return
		}
	}
	common := fmt.Sprintf("there's no %v in %v", v, s)
	output(t, common, out)
}

func NotContainsFunc[A any, B any, S []A](t testing.TB, s S, v B, fn func(A, B) bool, out ...any) {
	t.Helper()

	for _, e := range s {
		if fn(e, v) {
			common := fmt.Sprintf("there's %v in %v", v, s)
			output(t, common, out)
			return
		}
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

func HasSubString(t testing.TB, s string, sub string, out ...any) {
	t.Helper()

	if !strings.Contains(s, sub) {
		common := fmt.Sprintf("the %q cannot be a substring of the %q", sub, s)
		output(t, common, out)
	}
}

func HasNoSubString(t testing.TB, s string, sub string, out ...any) {
	t.Helper()

	if strings.Contains(s, sub) {
		common := fmt.Sprintf("the %q can be a substring of the %q", sub, s)
		output(t, common, out)
	}
}

func Panics(t testing.TB, fn func(), out ...any) {
	t.Helper()

	defer func() {
		if recover() == nil {
			output(t, "didn't panic", out)
		}
	}()

	fn()
}

func NotPanics(t testing.TB, fn func(), out ...any) {
	t.Helper()

	defer func() {
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
		if r := recover(); r != exp {
			common := fmt.Sprintf("can't get the expected panic, got %v", r)
			output(t, common, out)
		}
	}()

	fn()
}
