package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCastlePerm_Set(t *testing.T) {
	p := &castlePerm{}
	assert.Equal(t, &castlePerm{0}, p)

	p.Set(CASTLE_PERMS_WK)
	assert.Equal(t, 1, p.val)

	p.Set(CASTLE_PERMS_BK)
	assert.Equal(t, 5, p.val)

	p.Set(CASTLE_PERMS_ALL)
	assert.Equal(t, 15, p.val)

	p.Set(CASTLE_PERMS_NONE)
	assert.Equal(t, 15, p.val)
}

func TestCastlePerm_Clear(t *testing.T) {
	p := &castlePerm{}
	assert.Equal(t, &castlePerm{0}, p)

	p.Set(CASTLE_PERMS_WK)
	assert.Equal(t, 1, p.val)

	p.Set(CASTLE_PERMS_BK)
	assert.Equal(t, 5, p.val)

	p.Clear(CASTLE_PERMS_WK)
	assert.Equal(t, 4, p.val)

	p.Set(CASTLE_PERMS_ALL)
	assert.Equal(t, 15, p.val)

	p.Clear(CASTLE_PERMS_ALL)
	assert.Equal(t, 0, p.val)
}

func TestCastlePerm_Has(t *testing.T) {
	p := &castlePerm{}
	assert.Equal(t, &castlePerm{0}, p)

	p.Set(CASTLE_PERMS_WK)
	assert.True(t, p.Has(CASTLE_PERMS_WK))
	assert.False(t, p.Has(CASTLE_PERMS_WQ))
	assert.False(t, p.Has(CASTLE_PERMS_BK))
	assert.False(t, p.Has(CASTLE_PERMS_BQ))
	assert.False(t, p.Has(CASTLE_PERMS_ALL))

	p.Set(CASTLE_PERMS_ALL)
	assert.True(t, p.Has(CASTLE_PERMS_WK))
	assert.True(t, p.Has(CASTLE_PERMS_WQ))
	assert.True(t, p.Has(CASTLE_PERMS_BK))
	assert.True(t, p.Has(CASTLE_PERMS_BQ))
	assert.True(t, p.Has(CASTLE_PERMS_ALL))
}
