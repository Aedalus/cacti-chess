package engine

import (
	"fmt"
	"strconv"
	"strings"
)

//"fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func parseFen(fen string) (*State, error) {
	fenPieces := strings.Split(fen, " ")
	if len(fenPieces) != 6 {
		return nil, fmt.Errorf("fen should have 6 parts")
	}

	state := &State{}
	state.Reset()

	// pieces
	pieces, pieceErr := parsePiecesStr(fenPieces[0])
	if pieceErr != nil {
		return nil, fmt.Errorf("error parsing pieceStr: %v", pieceErr)
	}
	for i := 0; i < 64; i++ {
		sq120 := SQ120(i)
		state.pieces[sq120] = pieces[i]
	}

	// side
	side := fenPieces[1]
	if side == "w" {
		state.side = WHITE
	} else if side == "b" {
		state.side = BLACK
	} else {
		return nil, fmt.Errorf("error parsing sideStr: found %q", side)
	}

	// castlePerms
	permStr := fenPieces[2]
	state.castlePerm = parseCastlePerms(permStr)

	// enPas
	enPasStr := fenPieces[3]
	enPasSq, err := parseEnPas(enPasStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing enPasStr: %v", err)
	}
	state.enPas = enPasSq

	fiftyMove, err := strconv.Atoi(fenPieces[4])
	if err != nil {
		return nil, fmt.Errorf("error parsing fiftyMoveStr: %v", err)
	}
	state.fiftyMove = fiftyMove

	fullmovesCount, err := strconv.Atoi(fenPieces[5])
	if err != nil {
		return nil, fmt.Errorf("error parsing fullmoveStr: %v", err)
	}

	state.halfMoveCount = fullmovesCount / 2
	if state.side == BLACK {
		state.halfMoveCount++
	}

	// posKey
	state.posKey = state.GenPosKey()
	return state, nil
}

func parseEnPas(enPasStr string) (int, error) {
	if enPasStr == "-" {
		return NO_SQ, nil
	}

	fileStr := enPasStr[0:1]
	rankStr := enPasStr[1:2]
	var file int
	switch fileStr {
	case "a":
		file = 0
	case "b":
		file = 1
	case "c":
		file = 2
	case "d":
		file = 3
	case "e":
		file = 4
	case "f":
		file = 5
	case "g":
		file = 6
	case "h":
		file = 7
	default:
		return -1, fmt.Errorf("could not parse file %q", fileStr)
	}
	rank, rankErr := strconv.Atoi(rankStr)
	if rankErr != nil {
		return -1, rankErr
	}
	rankNum := rank - 1 // We start at 0, not 1
	return fileRankToSq(file, rankNum), nil
}

func parseCastlePerms(permStr string) *castlePerm {
	perms := &castlePerm{}
	if strings.ContainsRune(permStr, 'K') {
		perms.Set(CASTLE_PERMS_WK)
	}
	if strings.ContainsRune(permStr, 'Q') {
		perms.Set(CASTLE_PERMS_WQ)
	}
	if strings.ContainsRune(permStr, 'k') {
		perms.Set(CASTLE_PERMS_BK)
	}
	if strings.ContainsRune(permStr, 'q') {
		perms.Set(CASTLE_PERMS_BQ)
	}
	return perms
}

// parsePieceStr returns a len 64 array, with the start being
// A1, B1, etc, opposed to the fen string which starts with A8.
func parsePiecesStr(pieces string) ([64]piece, error) {
	rankStrings := strings.Split(pieces, "/")
	if len(rankStrings) != 8 {
		return [64]piece{}, fmt.Errorf("pieceStr should have 8 parts total")
	}

	// reverse the rankStrings so we start with 1st rank, not 8th
	for i := 0; i < len(rankStrings)/2; i++ {
		j := len(rankStrings) - i - 1
		rankStrings[i], rankStrings[j] = rankStrings[j], rankStrings[i]
	}

	board := [64]piece{}

	for rankPos := 7; rankPos >= 0; rankPos-- {
		filePos := 0
		rankStr := rankStrings[rankPos]
		for _, c := range rankStr {
			var pieceType piece

			// switch off of character
			switch c {
			case 'p':
				pieceType = bP
			case 'r':
				pieceType = bR
			case 'n':
				pieceType = bN
			case 'b':
				pieceType = bB
			case 'k':
				pieceType = bK
			case 'q':
				pieceType = bQ
			case 'P':
				pieceType = wP
			case 'R':
				pieceType = wR
			case 'N':
				pieceType = wN
			case 'B':
				pieceType = wB
			case 'K':
				pieceType = wK
			case 'Q':
				pieceType = wQ
			case '1', '2', '3', '4', '5', '6', '7', '8':
				pieceType = EMPTY
				emptySpaces, err := strconv.Atoi(string(c))
				if err != nil {
					return [64]piece{}, err
				}
				filePos += emptySpaces
			}
			if pieceType != EMPTY {
				index := (int(rankPos) * 8) + filePos
				board[index] = pieceType
				filePos++
			}
		}
	}
	return board, nil
}
