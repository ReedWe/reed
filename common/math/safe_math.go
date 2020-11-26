package math

import "github.com/reed/errors"

var (
	mathOverflowErr  = errors.New("math overflow")
	mathUnderflowErr = errors.New("math underflow")
)

func AddUint64(a, b uint64) (sum uint64,err error) {
	c := a + b
	if c < a {
		return 0, errors.Wrapf(mathOverflowErr, "AddUint64(a=%d,b=%d)", a, b)
	}
	return c, nil
}

func SubUint64(a, b uint64) (uint64, error) {
	if b > a {
		return 0, errors.Wrapf(mathUnderflowErr, "SubUint64(a=%d,b=%d)", a, b)
	}
	return a - b, nil
}
