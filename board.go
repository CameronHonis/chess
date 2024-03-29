package chess

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Board struct {
	Pieces                  [8][8]Piece      `json:"pieces"`
	OptEnPassantSquare      *Square          `json:"enPassantSquare"`
	IsWhiteTurn             bool             `json:"isWhiteTurn"`
	CanWhiteCastleQueenside bool             `json:"canWhiteCastleQueenside"`
	CanWhiteCastleKingside  bool             `json:"canWhiteCastleKingside"`
	CanBlackCastleQueenside bool             `json:"canBlackCastleQueenside"`
	CanBlackCastleKingside  bool             `json:"canBlackCastleKingside"`
	HalfMoveClockCount      uint8            `json:"halfMoveClockCount"`
	FullMoveCount           uint16           `json:"fullMoveCount"`
	RepetitionsByMiniFEN    map[string]uint8 `json:"repetitionsByMiniFEN"`
	Result                  BoardResult      `json:"result"`
	// memoizers
	optMaterialCount   *MaterialCount
	optWhiteKingSquare *Square
	optBlackKingSquare *Square
}

func NewBoard(pieces *[8][8]Piece,
	enPassantSquare *Square,
	isWhiteTurn bool,
	canWhiteCastleKingside bool,
	canWhiteCastleQueenside bool,
	canBlackCastleKingside bool,
	canBlackCastleQueenside bool,
	halfMoveClockCount uint8,
	fullMoveCount uint16,
	repetitionsByMiniFEN map[string]uint8,
	result BoardResult,
) *Board {
	return &Board{
		*pieces, enPassantSquare, isWhiteTurn,
		canWhiteCastleQueenside, canWhiteCastleKingside,
		canBlackCastleQueenside, canBlackCastleKingside,
		halfMoveClockCount, fullMoveCount, repetitionsByMiniFEN,
		result, nil, nil, nil,
	}
}

func BoardFromFEN(fen string) (*Board, error) {
	fen = strings.TrimSpace(fen)
	pieceByFENrune := map[rune]Piece{
		'p': BLACK_PAWN,
		'n': BLACK_KNIGHT,
		'b': BLACK_BISHOP,
		'r': BLACK_ROOK,
		'q': BLACK_QUEEN,
		'k': BLACK_KING,
		'P': WHITE_PAWN,
		'N': WHITE_KNIGHT,
		'B': WHITE_BISHOP,
		'R': WHITE_ROOK,
		'Q': WHITE_QUEEN,
		'K': WHITE_KING,
	}
	pieces := [8][8]Piece{
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	}
	boardBuilder := NewBoardBuilder()
	fenSegs := strings.Split(fen, " ")
	if len(fenSegs) != 6 {
		return nil, fmt.Errorf("invalid FEN: wrong number of FEN segments. Expected 6 vs. actual %d", len(fenSegs))
	}
	miniFEN := strings.Join(fenSegs[:4], " ")
	boardBuilder.WithRepetitionsByMiniFEN(map[string]uint8{miniFEN: uint8(1)})
	for fenSegIdx, fenSeg := range fenSegs {
		if fenSegIdx == 0 {
			materialCountBuilder := NewMaterialCountBuilder()
			rank := uint8(8)
			file := uint8(1)
			for _, FENRune := range []rune(fenSeg) {
				if rank < 1 {
					return nil, fmt.Errorf("invalid FEN: too many rows")
				}
				if FENRune == '/' {
					if file < 9 {
						return nil, fmt.Errorf("invalid FEN: not enough files at rank: %d", rank)
					}
					rank--
					file = 1
					continue
				}
				if file > 8 {
					return nil, fmt.Errorf("invalid FEN: too many files at rank: %d", rank)
				}
				if unicode.IsDigit(FENRune) {
					file += uint8(FENRune) - 48
				} else {
					piece, exists := pieceByFENrune[FENRune]
					if !exists {
						coords := (&Square{rank, file}).ToAlgebraicCoords()
						return nil, fmt.Errorf("invalid FEN: unknown piece char '%c' at %s", FENRune, coords)
					}
					pieces[rank-1][file-1] = piece
					materialCountBuilder.WithPiece(piece, &Square{rank, file})
					file++
				}
			}
			if rank > 1 {
				return nil, fmt.Errorf("invalid FEN: not enough rows")
			}
			boardBuilder.WithPieces(pieces)
			boardBuilder.board.optMaterialCount = materialCountBuilder.Build()
		} else if fenSegIdx == 1 {
			if fenSeg == "w" {
				boardBuilder.WithIsWhiteTurn(true)
			} else if fenSeg == "b" {
				boardBuilder.WithIsWhiteTurn(false)
			} else {
				return nil, fmt.Errorf("invalid FEN: unknown turn specifier %s", fenSeg)
			}
		} else if fenSegIdx == 2 {
			if fenSeg == "-" || fenSeg == "_" {
				continue
			}
			for _, castleRightsRune := range []rune(fenSeg) {
				if castleRightsRune == 'K' {
					boardBuilder.WithCanWhiteCastleKingside(true)
				} else if castleRightsRune == 'Q' {
					boardBuilder.WithCanWhiteCastleQueenside(true)
				} else if castleRightsRune == 'k' {
					boardBuilder.WithCanBlackCastleKingside(true)
				} else if castleRightsRune == 'q' {
					boardBuilder.WithCanBlackCastleQueenside(true)
				} else {
					return nil, fmt.Errorf("invalid FEN: unknown castle rights specifier, got %c", castleRightsRune)
				}
			}
		} else if fenSegIdx == 3 {
			if fenSeg == "-" || fenSeg == "_" {
				continue
			}
			enPassantSquare, err := SquareFromAlgebraicCoords(fenSeg)
			if err != nil {
				return nil, err
			}
			boardBuilder.WithEnPassantSquare(enPassantSquare)
		} else if fenSegIdx == 4 {
			halfMoveClockCount, intErr := strconv.Atoi(fenSeg)
			if intErr != nil {
				err := fmt.Errorf("invalid FEN: could not parse half move clock count, got error: %w", intErr)
				return nil, err
			}
			if halfMoveClockCount < 0 || halfMoveClockCount > 255 {
				err := fmt.Errorf("invalid FEN: half move clock count outside expected range [0, 255], got (%d)", halfMoveClockCount)
				return nil, err
			}
			boardBuilder.WithHalfMoveClockCount(uint8(halfMoveClockCount))
		} else if fenSegIdx == 5 {
			fullMoveCount, intErr := strconv.Atoi(fenSeg)
			if intErr != nil {
				err := fmt.Errorf("invalid FEN: could not parse full move count, got error: %w", intErr)
				return nil, err
			}
			if fullMoveCount < 0 || fullMoveCount > 65535 {
				err := fmt.Errorf("invalid FEN: full move count outside expected range [0, 65535], got (%d)", fullMoveCount)
				return nil, err
			}
			boardBuilder.WithFullMoveCount(uint16(fullMoveCount))
		}
	}

	boardBuilder.WithResult(BOARD_RESULT_IN_PROGRESS)
	prevBoard := NewBoardBuilder().FromBoard(boardBuilder.Build()).WithIsWhiteTurn(!boardBuilder.board.IsWhiteTurn).Build()
	UpdateBoardResult(prevBoard, boardBuilder, 0)
	return boardBuilder.Build(), nil
}

func GetInitBoard() *Board {
	miniFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -"
	repetitionsByMiniFEN := map[string]uint8{
		miniFEN: uint8(1),
	}
	return NewBoard(&[8][8]Piece{
		{WHITE_ROOK, WHITE_KNIGHT, WHITE_BISHOP, WHITE_QUEEN, WHITE_KING, WHITE_BISHOP, WHITE_KNIGHT, WHITE_ROOK},
		{WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN},
		{BLACK_ROOK, BLACK_KNIGHT, BLACK_BISHOP, BLACK_QUEEN, BLACK_KING, BLACK_BISHOP, BLACK_KNIGHT, BLACK_ROOK},
	}, nil, true, true, true, true,
		true, 0, 1, repetitionsByMiniFEN, BOARD_RESULT_IN_PROGRESS)
}

func (board *Board) GetPieceOnSquare(square *Square) Piece {
	return board.Pieces[square.Rank-1][square.File-1]
}

func (board *Board) IsForcedDrawByMaterial() bool {
	mat := board.ComputeMaterialCount()
	if mat.WhitePawnCount > 0 || mat.BlackPawnCount > 0 {
		return false
	}
	if mat.WhiteQueenCount > 0 || mat.BlackQueenCount > 0 {
		return false
	}
	if mat.WhiteRookCount > 0 || mat.BlackRookCount > 0 {
		return false
	}
	if mat.WhiteKnightCount > 1 || mat.BlackKnightCount > 1 {
		return false
	}

	if mat.WhiteLightBishopCount > 0 && mat.WhiteDarkBishopCount > 0 {
		return false
	}
	whiteBishopCount := mat.WhiteLightBishopCount + mat.WhiteDarkBishopCount
	if whiteBishopCount > 0 && mat.WhiteKnightCount > 0 {
		return false
	}

	if mat.BlackLightBishopCount > 0 && mat.BlackDarkBishopCount > 0 {
		return false
	}
	blackBishopCount := mat.BlackLightBishopCount + mat.BlackDarkBishopCount
	if blackBishopCount > 0 && mat.BlackKnightCount > 0 {
		return false
	}
	return true
}

func (board *Board) HasLegalNextMove() bool {
	moves := GetLegalMovesForKing(board)
	if len(moves) == 0 {
		checkingSquares := GetCheckingSquares(board, board.IsWhiteTurn)
		if len(checkingSquares) > 1 {
			// assume that it must be checkmate if the king is in double check and can't move anywhere
			return false
		} else {
			return HasLegalMove(board)
		}
	} else {
		return true
	}
}

func (board *Board) ComputeKingPositions() (*Square, *Square) {
	if board.optWhiteKingSquare != nil && board.optBlackKingSquare != nil {
		return board.optWhiteKingSquare, board.optBlackKingSquare
	}

	for r := uint8(0); r < 8; r++ {
		for c := uint8(0); c < 8; c++ {
			piece := board.Pieces[r][c]
			if piece == WHITE_KING {
				board.optWhiteKingSquare = &Square{r + 1, c + 1}
				if board.optBlackKingSquare != nil {
					return board.optWhiteKingSquare, board.optBlackKingSquare
				}
			} else if piece == BLACK_KING {
				board.optBlackKingSquare = &Square{r + 1, c + 1}
				if board.optWhiteKingSquare != nil {
					return board.optWhiteKingSquare, board.optBlackKingSquare
				}
			}
		}
	}
	return board.optWhiteKingSquare, board.optBlackKingSquare
}

func (board *Board) ComputeMaterialCount() *MaterialCount {
	if board.optMaterialCount != nil {
		return board.optMaterialCount
	}

	materialCountBuilder := NewMaterialCountBuilder()
	for r := uint8(0); r < 8; r++ {
		for c := uint8(0); c < 8; c++ {
			piece := board.Pieces[r][c]
			materialCountBuilder.WithPiece(piece, &Square{r + 1, c + 1})
		}
	}
	board.optMaterialCount = materialCountBuilder.Build()
	return board.optMaterialCount
}

func (board *Board) GetKingSquare(isWhiteKing bool) *Square {
	whiteSquare, blackSquare := board.ComputeKingPositions()
	if isWhiteKing {
		return whiteSquare
	} else {
		return blackSquare
	}
}

func (board *Board) ToFEN() string {
	var fenSegsBuilder strings.Builder

	pieceRuneByPiece := []rune{'x', 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}
	fenPiecesSeg := make([]rune, 0, 64)
	for r := 7; r >= 0; r-- {
		consecEmptyCount := 0
		for c := 0; c < 8; c++ {
			piece := board.Pieces[r][c]
			if piece == EMPTY {
				consecEmptyCount++
				continue
			}
			if consecEmptyCount > 0 {
				fenPiecesSeg = append(fenPiecesSeg, rune(consecEmptyCount+48))
				consecEmptyCount = 0
			}
			fenPiecesSeg = append(fenPiecesSeg, pieceRuneByPiece[piece])
		}
		if consecEmptyCount > 0 {
			fenPiecesSeg = append(fenPiecesSeg, rune(consecEmptyCount+48))
			consecEmptyCount = 0
		}
		if r > 0 {
			fenPiecesSeg = append(fenPiecesSeg, '/')
		}
	}
	fenSegsBuilder.WriteString(string(fenPiecesSeg))
	fenSegsBuilder.WriteRune(' ')

	var turnRune rune
	if board.IsWhiteTurn {
		turnRune = 'w'
	} else {
		turnRune = 'b'
	}
	fenSegsBuilder.WriteRune(turnRune)
	fenSegsBuilder.WriteRune(' ')

	fenCastleSeg := make([]rune, 0, 4)
	if board.CanWhiteCastleKingside {
		fenCastleSeg = append(fenCastleSeg, 'K')
	}
	if board.CanWhiteCastleQueenside {
		fenCastleSeg = append(fenCastleSeg, 'Q')
	}
	if board.CanBlackCastleKingside {
		fenCastleSeg = append(fenCastleSeg, 'k')
	}
	if board.CanBlackCastleQueenside {
		fenCastleSeg = append(fenCastleSeg, 'q')
	}
	if len(fenCastleSeg) == 0 {
		fenCastleSeg = append(fenCastleSeg, '-')
	}
	fenSegsBuilder.WriteString(string(fenCastleSeg))
	fenSegsBuilder.WriteRune(' ')

	if board.OptEnPassantSquare != nil {
		fenSegsBuilder.WriteString(board.OptEnPassantSquare.ToAlgebraicCoords())
	} else {
		fenSegsBuilder.WriteRune('-')
	}
	fenSegsBuilder.WriteRune(' ')

	fenSegsBuilder.WriteString(strconv.Itoa(int(board.HalfMoveClockCount)))
	fenSegsBuilder.WriteRune(' ')

	fenSegsBuilder.WriteString(strconv.Itoa(int(board.FullMoveCount)))

	return fenSegsBuilder.String()
}

func (board *Board) ToMiniFEN() string {
	fen := board.ToFEN()
	fenSegs := strings.Split(fen, " ")
	return strings.Join(fenSegs[:4], " ")
}

func (board *Board) IsInitBoard() bool {
	return board.ToMiniFEN() == "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -" && board.FullMoveCount == uint16(1)
}

func (board *Board) IsCheckmate() bool {
	return board.Result == BOARD_RESULT_WHITE_WINS_BY_CHECKMATE || board.Result == BOARD_RESULT_BLACK_WINS_BY_CHECKMATE
}

func (board *Board) String() string {
	return fmt.Sprintf("Board<%s>", board.ToFEN())
}

// pieceSquareByFile strictly used for finding origin square of a move expressed in algebraic notation
// where the origin square is not explicitly given.
func (board *Board) pieceSquareByRank(piece Piece, rank uint8) (*Square, error) {
	for file := uint8(1); file < 9; file++ {
		sqr := &Square{rank, file}
		_piece := board.GetPieceOnSquare(sqr)
		if piece == _piece {
			return sqr, nil
		}
	}
	return nil, fmt.Errorf("could not find peice %s on rank %s on %s", piece, string(rank), board.ToFEN())
}

// pieceSquareByFile strictly used for finding origin square of a move expressed in algebraic notation
// where the origin square is not explicitly given.
func (board *Board) pieceSquareByFile(piece Piece, file uint8) (*Square, error) {
	for rank := uint8(1); rank < 9; rank++ {
		sqr := &Square{rank, file}
		_piece := board.GetPieceOnSquare(sqr)
		if piece == _piece {
			return sqr, nil
		}
	}
	return nil, fmt.Errorf("could not find piece %s on file %s of %s", piece, string(file), board.ToFEN())
}

func (board *Board) pieceSquaresOnBoard(piece Piece) []*Square {
	var out = make([]*Square, 0)
	for rank := uint8(1); rank < 9; rank++ {
		for file := uint8(1); file < 9; file++ {
			sqr := &Square{rank, file}
			_piece := board.GetPieceOnSquare(sqr)
			if piece == _piece {
				out = append(out, sqr)
			}
		}
	}
	return out
}
