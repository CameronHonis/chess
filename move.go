package chess

import (
	"fmt"
	"strings"
)

type Move struct {
	Piece               Piece     `json:"piece"`
	StartSquare         *Square   `json:"startSquare"`
	EndSquare           *Square   `json:"endSquare"`
	CapturedPiece       Piece     `json:"capturedPiece"`
	KingCheckingSquares []*Square `json:"kingCheckingSquares"`
	PawnUpgradedTo      Piece     `json:"pawnUpgradedTo"`
}

// Deprecated: Replaced with "Equal" instead - due to support from google/go-cmp
func (move *Move) EqualTo(otherMove *Move) bool {
	return move.Equal(otherMove)
}

func (move *Move) Equal(otherMove *Move) bool {
	if len(move.KingCheckingSquares) != len(otherMove.KingCheckingSquares) {
		return false
	}
	for _, checkingSquare := range move.KingCheckingSquares {
		foundMatch := false
		for _, otherCheckingSquare := range otherMove.KingCheckingSquares {
			if checkingSquare.Equal(otherCheckingSquare) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			return false
		}
	}
	return move.Piece == otherMove.Piece &&
		move.StartSquare.Equal(otherMove.StartSquare) &&
		move.EndSquare.Equal(otherMove.EndSquare) &&
		move.CapturedPiece == otherMove.CapturedPiece &&
		move.PawnUpgradedTo == otherMove.PawnUpgradedTo
}

func (move *Move) DoesAllowEnPassant() bool {
	if !move.Piece.IsPawn() {
		return false
	}
	dis := int(move.StartSquare.Rank) - int(move.EndSquare.Rank)
	return dis == 2 || dis == -2
}

func (move *Move) IsCastles() bool {
	if !move.Piece.IsKing() {
		return false
	}
	return move.StartSquare.File == 5 && (move.EndSquare.File == 3 || move.EndSquare.File == 7)
}

func (move *Move) ToAlgebraic(board *Board) string {
	if move.IsCastles() {
		if move.EndSquare.File == 3 {
			return "O-O-O"
		}
		return "O-O"
	}
	algBuilder := strings.Builder{}
	writeMoveStartSpecifier := func() {
		if !move.Piece.IsPawn() {
			pieceAlg := move.Piece.ToAlgebraic()
			algBuilder.WriteString(pieceAlg)
			sharedTargetSquareSquares := make([]*Square, 0)
			if pieceAlg == "N" {
				sharedKnightSquares := []*Square{
					{Rank: move.EndSquare.Rank + 2, File: move.EndSquare.File + 1},
					{Rank: move.EndSquare.Rank + 1, File: move.EndSquare.File + 2},
					{Rank: move.EndSquare.Rank - 1, File: move.EndSquare.File + 2},
					{Rank: move.EndSquare.Rank - 2, File: move.EndSquare.File + 1},
					{Rank: move.EndSquare.Rank - 2, File: move.EndSquare.File - 1},
					{Rank: move.EndSquare.Rank - 1, File: move.EndSquare.File - 2},
					{Rank: move.EndSquare.Rank + 1, File: move.EndSquare.File - 2},
					{Rank: move.EndSquare.Rank + 2, File: move.EndSquare.File - 1},
				}
				for _, sharedKnightSquare := range sharedKnightSquares {
					if sharedKnightSquare.Equal(move.StartSquare) {
						continue
					}
					if !sharedKnightSquare.IsValidBoardSquare() {
						continue
					}
					if board.GetPieceOnSquare(sharedKnightSquare) == move.Piece {
						sharedTargetSquareSquares = append(sharedTargetSquareSquares, sharedKnightSquare)
					}
				}
			} else if pieceAlg == "B" || pieceAlg == "R" || pieceAlg == "Q" {
				var dirFromEndSquare [][2]int
				if pieceAlg == "B" {
					dirFromEndSquare = [][2]int{{1, 1}, {-1, -1}, {-1, 1}, {1, -1}}
				} else if pieceAlg == "R" {
					dirFromEndSquare = [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
				} else {
					dirFromEndSquare = [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {-1, -1}, {-1, 1}, {1, -1}}
				}
				for _, dir := range dirFromEndSquare {
					for i := 1; i < 8; i++ {
						targetSquare := &Square{
							uint8(int(move.EndSquare.Rank) + dir[0]*i),
							uint8(int(move.EndSquare.File) + dir[1]*i)}
						if !targetSquare.IsValidBoardSquare() {
							break
						}
						if targetSquare.Equal(move.StartSquare) {
							break
						}
						if board.GetPieceOnSquare(targetSquare) != EMPTY {
							if board.GetPieceOnSquare(targetSquare) == move.Piece {
								sharedTargetSquareSquares = append(sharedTargetSquareSquares, targetSquare)
							}
							break
						}
					}
				}
			}
			if len(sharedTargetSquareSquares) > 0 {
				sharedRank := false
				sharedFile := false
				for _, sharedTargetSquare := range sharedTargetSquareSquares {
					if sharedTargetSquare.Rank == move.StartSquare.Rank {
						sharedRank = true
					}
					if sharedTargetSquare.File == move.StartSquare.File {
						sharedFile = true
					}
				}
				if !sharedFile {
					algBuilder.WriteString(move.StartSquare.ToAlgebraicCoords()[:1])
				} else if !sharedRank {
					algBuilder.WriteString(move.StartSquare.ToAlgebraicCoords()[1:])
				} else {
					algBuilder.WriteString(move.StartSquare.ToAlgebraicCoords())
				}
			}
		} else {
			// is pawn
			if move.CapturedPiece != EMPTY {
				algBuilder.WriteString(move.StartSquare.ToAlgebraicCoords()[:1])
			}
		}
	}
	writeCaptureSpecifier := func() {
		if move.CapturedPiece != EMPTY {
			algBuilder.WriteString("x")
		}
	}
	writeLandSquareSpecifier := func() {
		algBuilder.WriteString(move.EndSquare.ToAlgebraicCoords())
	}
	writePawnUpgradeSpecifier := func() {
		if move.PawnUpgradedTo != EMPTY {
			algBuilder.WriteString("=" + move.PawnUpgradedTo.ToAlgebraic())
		}
	}
	writeCheckSpecifier := func() {
		if len(move.KingCheckingSquares) == 0 {
			return
		}
		if GetBoardFromMove(board, move).Result == BOARD_RESULT_IN_PROGRESS {
			algBuilder.WriteString("+")
		} else {
			algBuilder.WriteString("#")
		}
	}

	writeMoveStartSpecifier()
	writeCaptureSpecifier()
	writeLandSquareSpecifier()
	writePawnUpgradeSpecifier()
	writeCheckSpecifier()

	return algBuilder.String()
}

// ToLongAlgebraic adheres to the UCI "long algebraic notation", which differs from the standard
// "long algebraic notation". The UCI version is restricted to the start/end squares and a piece
// delimiter for upgrading pawns.
func (move *Move) ToLongAlgebraic() string {
	var upgradeChar string
	if move.PawnUpgradedTo.IsQueen() {
		upgradeChar = "q"
	} else if move.PawnUpgradedTo.IsRook() {
		upgradeChar = "r"
	} else if move.PawnUpgradedTo.IsBishop() {
		upgradeChar = "b"
	} else if move.PawnUpgradedTo.IsKnight() {
		upgradeChar = "n"
	}
	return fmt.Sprintf("%s%s%s", move.StartSquare.ToAlgebraicCoords(), move.EndSquare.ToAlgebraicCoords(), upgradeChar)
}

func MoveFromAlgebraic(algMove string, priorBoard *Board) (*Move, error) {
	originChars, targetChars, piece, upgradePiece := extractAlgebraicMoveInfo(algMove, priorBoard.IsWhiteTurn)

	landSqr, landSquareErr := SquareFromAlgebraicCoords(targetChars)
	if landSquareErr != nil {
		return nil, fmt.Errorf("could not create move %s: could not make land square: %s", algMove, landSquareErr)
	}

	var validMovesFromOrigin = make([]*Move, 0)
	var movesErr error
	if len(originChars) == 2 {
		originSqr, originSqrErr := SquareFromAlgebraicCoords(originChars)
		if originSqrErr != nil {
			return nil, fmt.Errorf("could not create move %s: could not make origin square: %s", algMove, originSqrErr)
		}
		validMovesFromOrigin, movesErr = GetLegalMovesFromOrigin(priorBoard, originSqr)
	} else if len(originChars) == 1 {
		originChar := originChars[0]
		if originChar >= 'a' && originChar <= 'h' {
			file := originChar - 'a' + 1
			originSqr, originSqrErr := priorBoard.pieceSquareByFile(piece, file)
			if originSqrErr != nil {
				return nil, fmt.Errorf("could not create move %s: %s", algMove, originSqrErr)
			}
			validMovesFromOrigin, movesErr = GetLegalMovesFromOrigin(priorBoard, originSqr)
		} else if originChar >= '1' && originChar <= '8' {
			rank := originChar - '1' + 1
			originSqr, originSqrErr := priorBoard.pieceSquareByRank(piece, rank)
			if originSqrErr != nil {
				return nil, fmt.Errorf("could not create move %s: %s", algMove, originSqrErr)
			}
			validMovesFromOrigin, movesErr = GetLegalMovesFromOrigin(priorBoard, originSqr)
		} else {
			return nil, fmt.Errorf("malformed origin char %s in move %s", string(originChar), algMove)
		}
	} else {
		pieceSqrs := priorBoard.pieceSquaresOnBoard(piece)
		if len(pieceSqrs) == 0 {
			return nil, fmt.Errorf("cannot create move, piece of type %s does not exist on %s", piece, priorBoard)
		}
		for _, pieceSqr := range pieceSqrs {
			var validMovesFromPieceSqr []*Move
			validMovesFromPieceSqr, movesErr = GetLegalMovesFromOrigin(priorBoard, pieceSqr)
			validMovesFromOrigin = append(validMovesFromOrigin, validMovesFromPieceSqr...)
		}
	}

	var validMoves = make([]*Move, 0)
	for _, move := range validMovesFromOrigin {
		if move.EndSquare.Equal(landSqr) {
			validMoves = append(validMoves, move)
		}
	}

	if movesErr != nil {
		return nil, fmt.Errorf("cannot create move, could not generate legal moves that match move %s on %s: %s", algMove, priorBoard, movesErr)
	}
	if len(validMoves) == 4 && upgradePiece != EMPTY {
		for _, validMove := range validMoves {
			if validMove.PawnUpgradedTo == upgradePiece {
				validMoves = []*Move{validMove}
				break
			}
		}
	}
	if len(validMoves) > 1 {
		return nil, fmt.Errorf("cannot create move, move repr %s is ambiguous. Rank and/or file needed on %s", algMove, priorBoard)
	}
	if len(validMoves) == 0 {
		return nil, fmt.Errorf("cannot create move, no legal moves exist with piece %s and landSqr %s on %s", piece, landSqr, priorBoard)
	}
	return validMoves[0], nil
}

func MoveFromLongAlgebraic(algMove string, priorBoard *Board) (*Move, error) {
	originChars, targetChars, _, upgradePiece := extractAlgebraicMoveInfo(algMove, priorBoard.IsWhiteTurn)
	if len(originChars) != 2 {
		return nil, fmt.Errorf("cannot create move, long algebraic origin must be a square, got: %s", originChars)
	}
	originSqr, originSqrErr := SquareFromAlgebraicCoords(originChars)
	if originSqrErr != nil {
		return nil, fmt.Errorf("cannot create move, error reading origin square from move %s: %s", algMove, originSqrErr)
	}
	landSqr, landSqrErr := SquareFromAlgebraicCoords(targetChars)
	if landSqrErr != nil {
		return nil, fmt.Errorf("cannot create move, error reading land square from move %s: %s", algMove, landSqrErr)
	}
	moves, movesErr := GetLegalMovesFromOrigin(priorBoard, originSqr)
	if movesErr != nil {
		return nil, fmt.Errorf("cannot create move, error getting moves from origin square %s on %s: %s", originChars, priorBoard, movesErr)
	}
	for _, move := range moves {
		if move.EndSquare.Equal(landSqr) {
			if move.PawnUpgradedTo != upgradePiece {
				continue
			}
			return move, nil
		}
	}
	return nil, fmt.Errorf("cannot create move, could not find move from %s to %s on %s", originChars, targetChars, priorBoard)
}

// extractAlgebraicMoveInfo supports both standard algebraic notation and long algebraic notation move formats
// It returns:
//   - the string repr of the move origin - either an empty string or rank/file or both
//   - the algebraic notation of the land square
//   - the piece involved in the move
//   - the piece the pawn was upgraded to
func extractAlgebraicMoveInfo(algMove string, isWhiteTurn bool) (string, string, Piece, Piece) {
	if algMove == "O-O" {
		if isWhiteTurn {
			return "e1", "g1", WHITE_KING, EMPTY
		} else {
			return "e8", "g8", BLACK_KING, EMPTY
		}
	} else if algMove == "O-O-O" {
		if isWhiteTurn {
			return "e1", "c1", WHITE_KING, EMPTY
		} else {
			return "e8", "c8", BLACK_KING, EMPTY
		}
	}
	var ptr = 0
	var piece Piece

	firstChar := algMove[0]
	if firstChar == 'N' {
		if isWhiteTurn {
			piece = WHITE_KNIGHT
		} else {
			piece = BLACK_KNIGHT
		}
		ptr++
	} else if firstChar == 'B' {
		if isWhiteTurn {
			piece = WHITE_BISHOP
		} else {
			piece = BLACK_BISHOP
		}
		ptr++
	} else if firstChar == 'R' {
		if isWhiteTurn {
			piece = WHITE_ROOK
		} else {
			piece = BLACK_ROOK
		}
		ptr++
	} else if firstChar == 'Q' {
		if isWhiteTurn {
			piece = WHITE_QUEEN
		} else {
			piece = BLACK_QUEEN
		}
		ptr++
	} else if firstChar == 'K' {
		if isWhiteTurn {
			piece = WHITE_KING
		} else {
			piece = BLACK_KING
		}
		ptr++
	} else {
		if isWhiteTurn {
			piece = WHITE_PAWN
		} else {
			piece = BLACK_PAWN
		}
	}

	var originChars string
	var targetChars string
	var upgradePiece = EMPTY
	for ptr < len(algMove) {
		c := algMove[ptr]
		if c >= 'a' && c <= 'h' {
			if targetChars != "" {
				// slide target to origin
				originChars = targetChars
			}
			targetChars = string(c)
		} else if c >= '1' && c <= '8' {
			targetChars += string(c)
		} else if c == 'Q' || c == 'q' {
			if isWhiteTurn {
				upgradePiece = WHITE_QUEEN
			} else {
				upgradePiece = BLACK_QUEEN
			}
		} else if c == 'R' || c == 'r' {
			if isWhiteTurn {
				upgradePiece = WHITE_ROOK
			} else {
				upgradePiece = BLACK_ROOK
			}
		} else if c == 'B' || c == 'b' {
			if isWhiteTurn {
				upgradePiece = WHITE_BISHOP
			} else {
				upgradePiece = BLACK_BISHOP
			}
		} else if c == 'N' || c == 'n' {
			if isWhiteTurn {
				upgradePiece = WHITE_KNIGHT
			} else {
				upgradePiece = BLACK_KNIGHT
			}
		}
		ptr++
	}

	return originChars, targetChars, piece, upgradePiece
}
