package gosnap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext(t, "test")

	assert.Equal(t, "test", ctx.Dir)
	assert.Equal(t, ".snap", ctx.FileExtension)
	assert.Equal(t, false, ctx.AutoUpdate)
	assert.Equal(t, 0, len(ctx.AutoUpdateSnapshots))
}

func TestContextNewSnapshot(t *testing.T) {
	ctx := NewContext(t, "test")

	s := ctx.NewSnapshot("test")

	assert.Equal(t, "test", s.Name)
	assert.Equal(t, ctx, s.ctx)
}

func TestContextShouldUpdateSnapshot(t *testing.T) {
	ctx := NewContext(t, "test")

	s := ctx.NewSnapshot("test A")
	assert.Equal(t, false, ctx.shouldUpdateSnapshot(s))

	ctx.AutoUpdate = true
	assert.Equal(t, true, ctx.shouldUpdateSnapshot(s))

	ctx.AutoUpdate = false
	ctx.AutoUpdateSnapshots = []string{s.Name}
	assert.Equal(t, true, ctx.shouldUpdateSnapshot(s))
}
