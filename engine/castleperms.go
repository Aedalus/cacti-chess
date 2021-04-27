package engine

import "strings"

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
func (perm *castlePerm) String() string {
	if perm.val == 0 {
		return "-"
	}

	builder := strings.Builder{}

	if perm.Has(CASTLE_PERMS_WK) {
		builder.WriteString("K")
	}
	if perm.Has(CASTLE_PERMS_WQ) {
		builder.WriteString("Q")
	}
	if perm.Has(CASTLE_PERMS_BK) {
		builder.WriteString("k")
	}
	if perm.Has(CASTLE_PERMS_BQ) {
		builder.WriteString("q")
	}

	return builder.String()
}
