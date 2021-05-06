package position

import (
	"fmt"
	"strings"
)

/*
The Movekey condenses all the information
about a player's move into 64 bits

                 11 1100 0000 0000 0000
               0100 0001 0110 1011 0110
4 ------- 3 ------- 2 ------- 1 -------
0000 0000 0000 0000 0000 0000 0xxx xxxx -> From
0000 0000 0000 0000 00xx xxxx x000 0000 -> To
0000 0000 0000 00xx xx00 0000 0000 0000 -> Captured (Piece)
0000 0000 0000 0x00 0000 0000 0000 0000 -> EnPas
0000 0000 0000 x000 0000 0000 0000 0000 -> Pawn Start
0000 0000 xxxx 0000 0000 0000 0000 0000 -> Promoted Piece (Piece)
0000 000x 0000 0000 0000 0000 0000 0000 -> Castle

Hexidecimal is easier to Count

0 -- 1 -- 8 -- 5 -- C -- 5 -- 8 -- f -- -> 0x185C58f
0000 0001 1000 0101 1100 0101 1000 1111
*/
type Movekey uint64

func (p *Position) ParseMove(str string) (Movekey, error) {
	mvb := []rune(str)
	if mvb[0] > 'h' || mvb[0] < 'a' {
		return Movekey(0), fmt.Errorf("str[0] must be a <= x <= h")
	}
	if mvb[1] > '8' || mvb[1] < '1' {
		return Movekey(0), fmt.Errorf("str[1] must be 1 <= x <= 8")
	}
	if mvb[2] > 'h' || mvb[2] < 'a' {
		return Movekey(0), fmt.Errorf("str[2] must be a <= x <= h")
	}
	if mvb[3] > '8' || mvb[3] < '1' {
		return Movekey(0), fmt.Errorf("str[3] must be 1 <= x <= 8")
	}

	from := fileRankToSq(int(mvb[0]-'a'), int(mvb[1]-'1'))
	to := fileRankToSq(int(mvb[2]-'a'), int(mvb[3]-'1'))
	prChar := ' '
	if len(mvb) > 4 {
		prChar = mvb[4]
	}

	possibleMoves := p.GenerateAllMoves()

	for _, mv := range *possibleMoves {
		// find a matching to/from move
		if mv.Key.getFrom() != from || mv.Key.getTo() != to {
			continue
		}
		// check if it was promoted
		pr := mv.Key.getPromoted()
		if pr != EMPTY {
			if (pr == PwR || pr == PbR) && prChar == 'r' {
				return mv.Key, nil
			} else if (pr == PwQ || pr == PbQ) && prChar == 'q' {
				return mv.Key, nil
			} else if (pr == PwN || pr == PbN) && prChar == 'n' {
				return mv.Key, nil
			} else if (pr == PwB || pr == PbB) && prChar == 'b' {
				return mv.Key, nil
			}
			continue
		} else {
			return mv.Key, nil
		}
	}

	return Movekey(0), nil
}

func (p *Position) MoveExists(m Movekey) bool {
	movelist := p.GenerateAllMoves()

	for _, mv := range *movelist {
		if mv.Key == m {
			return true
		}
	}

	return false
}

func (m *Movekey) ShortString() string {
	str := fmt.Sprintf("%s%s", printSq(m.getFrom()), printSq(m.getTo()))
	if m.getPromoted() != EMPTY {
		switch m.getPromoted() {
		case PwR, PbR:
			str += "r"
		case PwB, PbB:
			str += "b"
		case PwN, PbN:
			str += "n"
		case PwQ, PbQ:
			str += "q"
		}
	}
	return str
}

func (m *Movekey) String() string {
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
	moveKeyCapturedBitmask      uint64 = 0x3c000
	moveKeyEnPasBitmask         uint64 = 0x40000
	moveKeyPawnStartBitmask     uint64 = 0x80000
	moveKeyPromotedPieceBitmask uint64 = 0xf00000
	moveKeyCastleBitmask        uint64 = 0x1000000
)

// from
func (m Movekey) getFrom() int {
	return int(uint64(m) & moveKeyFromBitmask)
}
func (m Movekey) setFrom(sq int) Movekey {
	return Movekey((uint64(m) & ^moveKeyFromBitmask) | uint64(sq))
}

// to
func (m Movekey) getTo() int {
	return int((uint64(m) & moveKeyToBitmask) >> 7)
}
func (m Movekey) setTo(sq int) Movekey {
	return Movekey((uint64(m) & ^ moveKeyToBitmask) | (uint64(sq << 7)))
}

// enPas
func (m Movekey) isEnPas() bool {
	return uint64(m)&moveKeyEnPasBitmask != 0
}
func (m Movekey) setEnPas() Movekey {
	return Movekey(uint64(m) | moveKeyEnPasBitmask)
}
func (m Movekey) clearEnPas() Movekey {
	return Movekey(uint64(m) & ^moveKeyEnPasBitmask)
}

// pawnStart
func (m Movekey) isPawnStart() bool {
	return uint64(m)&moveKeyPawnStartBitmask != 0
}
func (m Movekey) setPawnStart() Movekey {
	return Movekey(uint64(m) | moveKeyPawnStartBitmask)
}
func (m Movekey) clearPawnStart() Movekey {
	return Movekey(uint64(m) & ^ moveKeyPawnStartBitmask)
}

// castle
func (m Movekey) isCastle() bool {
	return uint64(m)&moveKeyCastleBitmask != 0
}
func (m Movekey) setCastle() Movekey {
	return Movekey(uint64(m) | moveKeyCastleBitmask)
}
func (m Movekey) clearCastle() Movekey {
	return Movekey(uint64(m) & ^ moveKeyCastleBitmask)
}

// captured
func (m Movekey) getCaptured() Piece {
	return Piece((uint64(m) & moveKeyCapturedBitmask) >> 14)
}

func (m Movekey) setCaptured(p Piece) Movekey {
	return Movekey(uint64(m) & ^moveKeyCapturedBitmask | (uint64(p) << 14))
}

// promoted
func (m Movekey) getPromoted() Piece {
	return Piece((uint64(m) & moveKeyPromotedPieceBitmask) >> 20)
}

func (m Movekey) setPromoted(p Piece) Movekey {
	return Movekey(uint64(m) & ^moveKeyPromotedPieceBitmask | (uint64(p) << 20))
}

// no move
func (m Movekey) IsNoMove() bool {
	return m == 0
}
