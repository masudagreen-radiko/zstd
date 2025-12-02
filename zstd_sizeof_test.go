package zstd

import (
	"testing"
)

func TestGetSizeofCCtx(t *testing.T) {
	ctx := NewCtx().(*ctx)
	size := ctx.GetSizeofCCtx()
	if size == 0 {
		t.Error("Expected non-zero size for CCtx")
	}
	t.Logf("CCtx size: %d bytes", size)
}

func TestGetSizeofDCtx(t *testing.T) {
	ctx := NewCtx().(*ctx)
	size := ctx.GetSizeofDCtx()
	if size == 0 {
		t.Error("Expected non-zero size for DCtx")
	}
	t.Logf("DCtx size: %d bytes", size)
}
