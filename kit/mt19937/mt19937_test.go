package mt19937

import (
	"testing"
)

func TestNextInt32(t *testing.T) {
	expect := []int32{-1937831252}
	ctx := NewContext(0)
	got := ctx.NextInt32()
	if got != expect[0] {
		t.Error("expect:", expect[0], "got:", got)
	}
}
