package engine

import (
	"fmt"
	"strings"
)

type undo struct {
	move       int
	castlePerm castlePerm
	enPas      uint
	fiftyMove  uint
	posKey     uint64
}

// NewState returns a freshly initialized state
func NewState() *State {
	s := &State{}
	s.Reset()
	return s
}

type State struct {
	pieces *board120 // Source of truth for pieces
	side   int       // white/black

	pawns [3]*bitboard64 // white/black/both quick pawn lookups

	// piece list for fast lookup
	// pieceList[wN][0] = E1 etc
	pieceList [13][10]int

	kingSq [2]int // king quick lookups

	castlePerm *castlePerm // permissions to castle

	enPas int // if en passant is available

	fiftyMove int // 50 move counter (100 since we're using half moves)

	searchPly uint // how far we are in the search

	posKey uint64 // unique index for position

	pieceCount    [13]int // count of all pieces on the board
	bigPieceCount [2]int  // white/black non-pawn pieces
	majPieceCount [2]int  // queen/rook
	minPieceCount [2]int  // bishop/knight
	materialCount [2]int  // count of material

	// history
	halfMoveCount int // how many half moves have been made in the whole game
	history       *[2048]undo
}

// GenPosKey generates a statistically unique uint64
// key for the current state of the engine
func (s State) GenPosKey() uint64 {
	var finalKey uint64 = 0
	var p piece = piece(0)

	// pieces
	for sq := 0; sq < BOARD_SQ_NUMBER; sq++ {
		p = s.pieces[sq]
		if p != NO_SQ && p != EMPTY {
			finalKey ^= hashPieceKeys[p][sq]
		}
	}

	// side
	if s.side == WHITE {
		finalKey ^= hashSideKey
	}

	// enPas
	if s.enPas != NO_SQ {
		finalKey ^= hashPieceKeys[EMPTY][s.enPas]
	}

	// castle keys
	finalKey ^= hashCastleKeys[s.castlePerm.val]

	return finalKey
}

func (s *State) updateListCaches() {
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		piece := s.pieces[i]
		if piece != NO_SQ && piece != EMPTY {
			color := pieceColorMap[piece]
			if pieceBigMap[piece] {
				s.bigPieceCount[color]++
			}
			if pieceMajorMap[piece] {
				s.majPieceCount[color]++
			}
			if pieceMinorMap[piece] {
				s.minPieceCount[color]++
			}

			s.materialCount[color] += pieceValueMap[piece]

			// update the pieceList, then increment the counter
			// conceptually like
			// [wP][0] = A2
			// [wP][1] = B2 etc
			curPieceCount := s.pieceCount[piece]
			s.pieceList[piece][curPieceCount] = i
			s.pieceCount[piece]++

			// update king positions
			if piece == wK {
				s.kingSq[WHITE] = i
			}
			if piece == bK {
				s.kingSq[BLACK] = i
			}
		}
	}
}

// Reset will re-initialize the engine to an empty state
func (s *State) Reset() {
	// initialize piece array
	s.pieces = &board120{}
	for i := 0; i < BOARD_SQ_NUMBER; i++ {
		s.pieces[i] = NO_SQ
	}
	for i := 0; i < 64; i++ {
		s.pieces[SQ120(i)] = EMPTY
	}

	// castle perms
	s.castlePerm = &castlePerm{CASTLE_PERMS_NONE}

	// piece counts
	for i := 0; i < 13; i++ {
		s.pieceCount[i] = 0
	}
	for i := 0; i < 2; i++ {
		s.bigPieceCount[i] = 0
		s.majPieceCount[i] = 0
		s.minPieceCount[i] = 0
		s.materialCount[i] = 0
		s.pawns[i] = &bitboard64{0}
	}

	s.pawns[WHITE] = &bitboard64{}
	s.pawns[BLACK] = &bitboard64{}
	s.pawns[BOTH] = &bitboard64{}

	s.kingSq[WHITE] = NO_SQ
	s.kingSq[BLACK] = NO_SQ

	s.side = BOTH
	s.enPas = NO_SQ
	s.fiftyMove = 0
	s.searchPly = 0
	s.halfMoveCount = 0
	s.history = &[2048]undo{}
	s.posKey = 0
}

func (s State) String() string {
	output := strings.Builder{}
	output.WriteString(s.PrintBoard())
	output.WriteString("--------\n")
	output.WriteString(fmt.Sprintf("side: %v\n", s.side))
	output.WriteString(fmt.Sprintf("enPas: %v\n", s.enPas))
	output.WriteString(fmt.Sprintf("castle: %v\n", s.castlePerm))
	output.WriteString(fmt.Sprintf("posKey: %x\n", s.posKey))

	return output.String()
}
func (s State) PrintBoard() string {
	output := strings.Builder{}
	for r := RANK_8; r >= RANK_1; r-- {
		for f := FILE_A; f <= FILE_H; f++ {
			sq := fileRankToSq(int(f), int(r))
			piece := s.pieces[sq]
			output.WriteString(piece.String())
		}
		output.WriteString("\n")
	}
	return output.String()
}
