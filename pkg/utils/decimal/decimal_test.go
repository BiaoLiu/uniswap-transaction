package decimal

import (
	"fmt"
	"testing"
)

func TestCastDecimalFromInt(t *testing.T) {
	d := CastDecimalFromInt(10000000)
	fmt.Println(d.String())
}

func TestFloat64Mul(t *testing.T) {
	d := Float64Mul(1.00, 100)
	fmt.Println(d)
}

func TestCastCentToYuanString(t *testing.T) {
	v := CastCentToYuanString(10)
	v = fmt.Sprintf("%s元", v)
	fmt.Println(v)
}
