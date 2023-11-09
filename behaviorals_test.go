package chess_test

import (
	"github.com/CameronHonis/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// various situations that failed in practice
var _ = Describe("Behaviorals", func() {
	Describe("Scotch Game", func() {
		It("should not generate a terminal board", func() {
			board := chess.GetInitBoard()
			move := &chess.Move{
				Piece:          chess.WHITE_PAWN,
				StartSquare:    &chess.Square{Rank: 2, File: 5},
				EndSquare:      &chess.Square{Rank: 4, File: 5},
				CapturedPiece:  chess.EMPTY,
				PawnUpgradedTo: chess.EMPTY,
			}
			Expect(chess.IsLegalMove(board, move)).To(BeTrue())
			board = chess.GetBoardFromMove(board, move)
			move = &chess.Move{
				Piece:          chess.BLACK_PAWN,
				StartSquare:    &chess.Square{Rank: 7, File: 5},
				EndSquare:      &chess.Square{Rank: 5, File: 5},
				CapturedPiece:  chess.EMPTY,
				PawnUpgradedTo: chess.EMPTY,
			}
			Expect(chess.IsLegalMove(board, move)).To(BeTrue())
			board = chess.GetBoardFromMove(board, move)
			move = &chess.Move{
				Piece:          chess.WHITE_PAWN,
				StartSquare:    &chess.Square{Rank: 2, File: 4},
				EndSquare:      &chess.Square{Rank: 4, File: 4},
				CapturedPiece:  chess.EMPTY,
				PawnUpgradedTo: chess.EMPTY,
			}
			Expect(chess.IsLegalMove(board, move)).To(BeTrue())
			board = chess.GetBoardFromMove(board, move)
			Expect(board.IsTerminal).To(BeFalse())
		})
	})
})
