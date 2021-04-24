package board

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
