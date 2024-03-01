package chess

import (
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

func (move *Move) EqualTo(otherMove *Move) bool {
	if len(move.KingCheckingSquares) != len(otherMove.KingCheckingSquares) {
		return false
	}
	for _, checkingSquare := range move.KingCheckingSquares {
		foundMatch := false
		for _, otherCheckingSquare := range otherMove.KingCheckingSquares {
			if checkingSquare.EqualTo(otherCheckingSquare) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			return false
		}
	}
	return move.Piece == otherMove.Piece &&
		move.StartSquare.EqualTo(otherMove.StartSquare) &&
		move.EndSquare.EqualTo(otherMove.EndSquare) &&
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
					if sharedKnightSquare.EqualTo(move.StartSquare) {
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
						if targetSquare.EqualTo(move.StartSquare) {
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
