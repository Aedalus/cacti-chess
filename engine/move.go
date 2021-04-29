package engine

import (
	"fmt"
	"strings"
)

/*
The movekey condenses all the information
about a player's move into 64 bits

4 ------- 3 ------- 2 ------- 1 -------
0000 0000 0000 0000 0000 0000 0xxx xxxx -> From
0000 0000 0000 0000 00xx xxxx x000 0000 -> To
0000 0000 0000 00xx xx00 0000 0000 0000 -> Captured (piece)
0000 0000 0000 0x00 0000 0000 0000 0000 -> EnPas
0000 0000 0000 x000 0000 0000 0000 0000 -> Pawn Start
0000 0000 xxxx 0000 0000 0000 0000 0000 -> Promoted Piece (piece)
0000 000x 0000 0000 0000 0000 0000 0000 -> Castle

Hexidecimal is easier to count

0 -- 1 -- 8 -- 5 -- C -- 5 -- 8 -- f -- -> 0x185C58f
0000 0001 1000 0101 1100 0101 1000 1111
*/
type movekey struct {
	val uint64
}

func newMovekey(from, to int, captured, promoted piece, enPas, pawnStart bool) *movekey {
	k := &movekey{}

	k.setFrom(from)
	k.setTo(to)
	k.setCaptured(captured)
	k.setPromoted(promoted)

	if enPas {
		k.setEnPas()
	}

	if pawnStart {
		k.setPawnStart()
	}

	return k
}

func (m *movekey) ShortString() string {
	str := fmt.Sprintf("%s%s", printSq(m.getFrom()), printSq(m.getTo()))
	if m.getPromoted() != EMPTY {
		switch m.getPromoted() {
		case wR, bR:
			str += "r"
		case wB, bB:
			str += "b"
		case wN, bN:
			str += "n"
		case wQ, bQ:
			str += "q"
		}
	}
	return str
}

func (m *movekey) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("from: %v\n", printSq(m.getFrom())))
	b.WriteString(fmt.Sprintf("to: %v\n", printSq(m.getTo())))
	b.WriteString(fmt.Sprintf("enPas: %v\n", m.isEnPas()))
	b.WriteString(fmt.Sprintf("castle: %v\n", m.isCastle()))
	b.WriteString(fmt.Sprintf("pawnStart: %v\n", m.isPawnStart()))

	if m.getCaptured() == EMPTY {
		b.WriteString(fmt.Sprintf("captured: EMPTY\n"))
	} else {
		b.WriteString(fmt.Sprintf("captured: %v\n", m.getCaptured()))
	}

	if m.getPromoted() == EMPTY {
		b.WriteString(fmt.Sprintf("promoted: EMPTY\n"))
	} else {
		b.WriteString(fmt.Sprintf("promoted: %v\n", m.getPromoted()))
	}

	return b.String()
}

func printSq(sq int) string {
	f := fileLookups[sq]
	r := rankLookups[sq]

	return fmt.Sprintf("%c%c", 'a'+f, '1'+r)
}

const (
	moveKeyFromBitmask          uint64 = 0x7f
	moveKeyToBitmask            uint64 = 0x7F << 7
	moveKeyCapturedBitmask      uint64 = 0x7c000
	moveKeyEnPasBitmask         uint64 = 0x40000
	moveKeyPawnStartBitmask     uint64 = 0x80000
	moveKeyPromotedPieceBitmask uint64 = 0xf00000
	moveKeyCastleBitmask        uint64 = 0x1000000
)

// from
func (m *movekey) getFrom() int {
	return int(m.val & moveKeyFromBitmask)
}
func (m *movekey) setFrom(sq int) {
	m.val = (m.val & ^moveKeyFromBitmask) | uint64(sq)
}

// to
func (m *movekey) getTo() int {
	return int((m.val & moveKeyToBitmask) >> 7)
}
func (m *movekey) setTo(sq int) {
	m.val = (m.val & ^ moveKeyToBitmask) | (uint64(sq << 7))
}

// enPas
func (m *movekey) isEnPas() bool {
	return m.val&moveKeyEnPasBitmask != 0
}
func (m *movekey) setEnPas() {
	m.val = m.val | moveKeyEnPasBitmask
}
func (m *movekey) clearEnPas() {
	m.val = m.val & ^moveKeyEnPasBitmask
}

// pawnStart
func (m *movekey) isPawnStart() bool {
	return m.val&moveKeyPawnStartBitmask != 0
}
func (m *movekey) setPawnStart() {
	m.val = m.val | moveKeyPawnStartBitmask
}
func (m *movekey) clearPawnStart() {
	m.val = m.val & ^ moveKeyPawnStartBitmask
}

// castle
func (m *movekey) isCastle() bool {
	return m.val&moveKeyCastleBitmask != 0
}
func (m *movekey) setCastle() {
	m.val = m.val | moveKeyCastleBitmask
}
func (m *movekey) clearCastle() {
	m.val = m.val & ^ moveKeyCastleBitmask
}

// captured
func (m *movekey) getCaptured() piece {
	return piece((m.val & moveKeyCapturedBitmask) >> 14)
}

func (m *movekey) setCaptured(p piece) {
	m.val = m.val & ^moveKeyCapturedBitmask | (uint64(p) << 14)
}

// promoted
func (m *movekey) getPromoted() piece {
	return piece((m.val & moveKeyPromotedPieceBitmask) >> 20)
}

func (m *movekey) setPromoted(p piece) {
	m.val = m.val & ^moveKeyPromotedPieceBitmask | (uint64(p) << 20)
}
