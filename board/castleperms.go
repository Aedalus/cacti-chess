package board

const (
	CASTLE_PERMS_WK = 1 << iota
	CASTLE_PERMS_WQ
	CASTLE_PERMS_BK
	CASTLE_PERMS_BQ
	CASTLE_PERMS_NONE = 0
	CASTLE_PERMS_ALL  = 15
)

type castlePerm struct {
	val int
}

func (perm *castlePerm) Set(b int)      { perm.val = perm.val | b }
func (perm *castlePerm) Clear(b int)    { perm.val = perm.val &^ b }
func (perm *castlePerm) Toggle(b int)   { perm.val = perm.val ^ b }
func (perm *castlePerm) Has(b int) bool { return perm.val&b == b }
