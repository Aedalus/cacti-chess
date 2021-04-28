package engine

type move struct {
	move  int
	score int
}

// need up to 7 bits to represent a move
// 1 2 4 8 15 32 64

/*
4 ------- 3 ------- 2 ------- 1 -------
0000 0000 0000 0000 0000 0000 0xxx xxxx -> From
0000 0000 0000 0000 00xx xxxx x000 0000 -> To >> 7, 0x3F
0000 0000 0000 00xx xx00 0000 0000 0000 -> Captured >> 14 0xF
0000 0000 0000 0x00 0000 0000 0000 0000 -> EP 0x40000
0000 0000 0000 x000 0000 0000 0000 0000 -> Pawn Start 0x80000
0000 0000 xxxx 0000 0000 0000 0000 0000 -> Promoted Piece >> 20, 0xF
0000 000x 0000 0000 0000 0000 0000 0000 -> Castle 0x1000000

Hexidecimal is easier to count

0 -- 1 -- 8 -- 5 -- C -- 5 -- 8 -- f -- -> 0x185C58f
0000 0001 1000 0101 1100 0101 1000 1111
*/
