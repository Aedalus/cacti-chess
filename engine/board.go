package engine

import (
	"strings"
)

/**
A 120 engine has margins around the edges to help
with move computation. The indexes still go from
0-119 total, though the main engine is a
non-contiguous range of 21-98.

-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
-1, 91, 92, 93, 94, 95, 96, 97, 98, -1,
-1, 81, 82, 83, 84, 85, 86, 87, 88, -1,
-1, 71, 72, 73, 74, 75, 76, 77, 78, -1,
-1, 61, 62, 63, 64, 65, 66, 67, 68, -1,
-1, 51, 52, 53, 54, 55, 56, 57, 58, -1,
-1, 41, 42, 43, 44, 45, 46, 47, 48, -1,
-1, 31, 32, 33, 34, 35, 36, 37, 38, -1,
-1, 21, 22, 23, 24, 25, 26, 27, 28, -1,
-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
*/
type board120 [120]piece

// bitboard64 is a representation of just an 8x8
type bitboard64 struct {
	val uint64
}

func (b *bitboard64) set(index int)    { b.val = b.val | 1<<index }
func (b *bitboard64) clear(index int)  { b.val = b.val &^ 1 << index }
func (b *bitboard64) toggle(index int) { b.val = b.val ^ 1<<index }
func (b *bitboard64) has(index int) bool {
	return (b.val & (1 << index)) != 0
}

// mapping indices from 120 <-> 64
var sq120to64 = [120]int{}
var sq64to120 = [64]int{}

func SQ120(sq64 int) int {
	return sq64to120[sq64]
}

func SQ64(sq120 int) int {
	return sq120to64[sq120]
}

func (b bitboard64) Pop() int {
	for i := 0; i < 64; i++ {
		if b.has(i) {
			b.clear(i)
			return i
		}
	}
	return -1
}

func (b bitboard64) Count() int {
	c := 0
	for i := 0; i < 64; i++ {
		if b.has(i) {
			c++
		}
	}
	return c
}

func fileRankToSq(file int, rank int) int {
	return (21 + file) + (rank * 10)
}



func (b bitboard64) String() string {
	str := strings.Builder{}
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			if b.has(rank*8 + file) {
				str.WriteString("X")
			} else {
				str.WriteString(".")
			}
		}
		str.WriteString("\n")
	}
	return str.String()
}

func init() {
	// initialize both boards to not found
	for i := 0; i < 120; i++ {
		sq120to64[i] = -1
	}
	for i := 0; i < 64; i++ {
		sq64to120[i] = -1
	}

	// fill them in
	sq64 := 0 // count the current board64 index
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sq120 := fileRankToSq(file, rank)
			sq120to64[sq120] = sq64
			sq64to120[sq64] = sq120
			sq64++
		}
	}
}
