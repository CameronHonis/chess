package chess_test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/CameronHonis/chess"
)

func compareSquares(expSquares []Square, realSquares []*Square) {
	Expect(realSquares).To(HaveLen(len(expSquares)))
	for _, realSquare := range realSquares {
		foundMatch := false
		for _, expSquare := range expSquares {
			if realSquare.Rank == expSquare.Rank && realSquare.File == expSquare.File {
				foundMatch = true
				break
			}
		}
		Expect(foundMatch).To(BeTrue(), "unexpected square %+v", realSquare)
	}
}

func compareMoves(expMoves []Move, realMoves []*Move) {
	Expect(realMoves).To(HaveLen(len(expMoves)))
	for _, realMove := range realMoves {
		foundMatch := false
		for _, expMove := range expMoves {
			if expMove.EqualTo(realMove) {
				foundMatch = true
				break
			}
		}
		Expect(foundMatch).To(BeTrue(), "unexpected move %+v", *realMove)
	}
}

var _ = Describe("GameHelpers", func() {
	Describe("#GetCheckingKingSquares", func() {
		When("the board is the initial board", func() {
			It("returns an empty list", func() {
				board := GetInitBoard()
				checkingSquares := GetCheckingSquares(board, board.IsWhiteTurn)
				expSquares := make([]Square, 0)
				compareSquares(expSquares, checkingSquares)
			})
		})
		When("a single pawn is checking the king", func() {
			It("returns the square of the pawn", func() {
				board, err := BoardFromFEN("8/8/4k3/3P4/8/8/8/7K b - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				checkingSquares := GetCheckingSquares(board, false)
				expSquares := []Square{{5, 4}}
				compareSquares(expSquares, checkingSquares)
			})
		})
		When("there are multiple rooks on the board", func() {
			It("returns the square of each rook that is checking the king", func() {
				board, err := BoardFromFEN("3R2R1/8/2R5/2Rk2R1/4R3/2R5/R2R4/7K w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				rookSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{2, 4},
					{8, 4},
					{5, 3},
					{5, 7},
				}
				compareSquares(expSquares, rookSquares)

			})
		})
		When("there are multiple bishops on the board", func() {
			It("returns the square of each bishop that is checking the king", func() {
				board, err := BoardFromFEN("3BB2B/5B2/B2B1k1B/8/4BB1B/8/8/B5K1 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				bishopSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{1, 1},
					{8, 8},
					{8, 4},
					{4, 8},
				}
				compareSquares(expSquares, bishopSquares)
			})
		})
		When("there are multiple pawns on the board", func() {
			Context("when the pawns are white", func() {
				It("returns the square of each pawn that is checking the king", func() {
					board, err := BoardFromFEN("3PP2P/4PPP1/P2PPkPP/4PPP1/4PP1P/8/1K6/P7 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					pawnSquares := GetCheckingSquares(board, false)
					expSquares := []Square{
						{5, 5},
						{5, 7},
					}
					compareSquares(expSquares, pawnSquares)
				})
			})
			Context("when the pawns are black", func() {
				It("returns the square of each pawn that is checking the king", func() {
					board, err := BoardFromFEN("3pp2p/1k2ppp1/p2ppKpp/4ppp1/4pp1p/8/8/p7 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					pawnSquares := GetCheckingSquares(board, true)
					expSquares := []Square{
						{7, 5},
						{7, 7},
					}
					compareSquares(expSquares, pawnSquares)
				})
			})
		})
		When("there are multiple knights on the board", func() {
			It("returns the square of each knight checking the king", func() {
				board, err := BoardFromFEN("4NNN1/3NNNNN/2N2k1N/3N1N1N/4N1N1/1K6/8/5N2 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				knightSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{4, 5},
					{4, 7},
					{5, 4},
					{5, 8},
					{7, 4},
					{7, 8},
					{8, 5},
					{8, 7},
				}
				compareSquares(expSquares, knightSquares)
			})
		})
		When("there are multiple queens on the board", func() {
			It("returns the square of each knight checking the king", func() {
				board, err := BoardFromFEN("3QQ1Q1/3Q1QQQ/Q4k1Q/3Q1QQQ/4Q1Q1/8/5K2/Q7 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				queenSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{1, 1},
					{6, 1},
					{8, 4},
					{7, 6},
					{7, 7},
					{6, 8},
					{5, 7},
					{5, 6},
				}
				compareSquares(expSquares, queenSquares)
			})
		})
		When("only same color pieces as king exist", func() {
			It("returns no checking squares", func() {
				board, err := BoardFromFEN("3q4/4p3/2q2k2/4p3/3bn1n1/5r2/1K6/8 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				checkingSquares := GetCheckingSquares(board, false)
				Expect(checkingSquares).To(HaveLen(0))
			})
		})
		When("all 'blockable' pieces are blocked", func() {
			It("returns no checking squares", func() {
				board, err := BoardFromFEN("3Q4/4n3/1Qq2kpR/6N1/7B/2N2N2/1B3R2/3K4 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				checkingSquares := GetCheckingSquares(board, false)
				Expect(checkingSquares).To(HaveLen(0))
			})
		})
	})
	Describe("#GetLegalMovesForPawn", func() {
		When("the pawn can capture in either direction", func() {
			Context("and the pawn is not blocked", func() {
				It("returns 2 capture moves and one non-capture move", func() {
					board, err := BoardFromFEN("k1K5/8/8/4p1n1/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 7}, BLACK_KNIGHT, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
				When("a pawn capture results in a check", func() {
					It("returns only the capturing move", func() {
						board, err := BoardFromFEN("3qkbnr/Bp1npb1p/2Pp1p2/p4Pp1/4P3/2N5/PPP3PP/R2QKBNR w KQk g6 1 10")
						Expect(err).ToNot(HaveOccurred())
						realMoves, err := GetLegalMovesForPawn(board, &Square{6, 3})
						Expect(err).ToNot(HaveOccurred())
						expMoves := []Move{
							{WHITE_PAWN, &Square{6, 3}, &Square{7, 4}, BLACK_KNIGHT, []*Square{{7, 4}}, EMPTY},
							{WHITE_PAWN, &Square{6, 3}, &Square{7, 3}, EMPTY, []*Square{}, EMPTY},
							{WHITE_PAWN, &Square{6, 3}, &Square{7, 2}, BLACK_PAWN, []*Square{}, EMPTY},
						}
						compareMoves(expMoves, realMoves)
					})
				})
			})
			Context("and the pawn is blocked", func() {
				It("returns 2 capture moves", func() {
					board, err := BoardFromFEN("k1K5/8/8/4ppn1/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 7}, BLACK_KNIGHT, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
		})
		When("the pawn can capture in one direction", func() {
			Context("and the pawn is not blocked", func() {
				It("returns 1 capture moves and a non-capturing move", func() {
					board, err := BoardFromFEN("k1K5/8/8/4p3/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
			Context("and the pawn is blocked", func() {
				It("returns 1 capture moves", func() {
					board, err := BoardFromFEN("k7/8/8/4rq2/5P2/8/8/7K w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_ROOK, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
		})
		When("the pawn cannot capture either direction", func() {
			Context("and the pawn is not blocked", func() {
				It("returns a non-capturing move", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
			Context("and the pawn is blocked", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k1K5/8/8/5n2/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					Expect(realMoves).To(HaveLen(0))
				})
			})
			Context("and the 'attacked' squares are occupied by friendly pieces", func() {
				It("returns only non-capturing moves", func() {
					board, err := BoardFromFEN("8/8/k7/8/3B1R2/4P3/7K/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
		})
		When("the pawn can capture en passant to the left", func() {
			It("includes the en passant capture move", func() {
				board, err := BoardFromFEN("k1K5/8/8/4pP2/8/8/8/8 w - e6 0 1")
				Expect(err).ToNot(HaveOccurred())
				realMoves, err := GetLegalMovesForPawn(board, &Square{5, 6})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 6}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the pawn can capture en passant to the right", func() {
			It("includes the en passant capture move", func() {
				board, err := BoardFromFEN("k1K5/8/8/5Pp1/8/8/8/8 w - g6 0 1")
				Expect(err).ToNot(HaveOccurred())
				realMoves, err := GetLegalMovesForPawn(board, &Square{5, 6})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 7}, BLACK_PAWN, make([]*Square, 0), EMPTY},
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 6}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the pawn is on the starting row", func() {
			Context("and both squares directly in front are not occupied", func() {
				It("returns all possible moves including a double jump", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/8/8/1P6/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{2, 2})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{2, 2}, &Square{3, 2}, EMPTY, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{2, 2}, &Square{4, 2}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
			Context("and the double jump square is blocked", func() {
				It("returns only the single jump non-capturing move", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/1B6/8/1P6/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{2, 2})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{2, 2}, &Square{3, 2}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
			Context("and the square directly in front is occupied", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/8/1R6/1P6/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{2, 2})
					Expect(err).ToNot(HaveOccurred())
					Expect(realMoves).To(HaveLen(0))
				})
			})
		})
		When("the pawn can be promoted", func() {
			When("and the pawn can capture", func() {
				Context("and the square in front is not occupied", func() {
					It("returns both capturing and non-capturing promotion moves", func() {
						board, err := BoardFromFEN("r7/1P6/8/8/8/8/8/2K3k1 w - - 0 1")
						Expect(err).ToNot(HaveOccurred())
						realMoves, err := GetLegalMovesForPawn(board, &Square{7, 2})
						Expect(err).ToNot(HaveOccurred())
						expMoves := []Move{
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_KNIGHT},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_BISHOP},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_ROOK},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_QUEEN},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_KNIGHT},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_BISHOP},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_ROOK},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_QUEEN},
						}
						compareMoves(expMoves, realMoves)
					})
				})
				Context("and the square in front is occupied", func() {
					It("only returns capture promotion moves", func() {
						board, err := BoardFromFEN("rn6/1P6/8/8/8/8/8/k1K5 w - - 0 1")
						Expect(err).ToNot(HaveOccurred())
						realMoves, err := GetLegalMovesForPawn(board, &Square{7, 2})
						Expect(err).ToNot(HaveOccurred())
						expMoves := []Move{
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_KNIGHT},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_BISHOP},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, []*Square{{8, 1}}, WHITE_ROOK},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, []*Square{{8, 1}}, WHITE_QUEEN},
						}
						compareMoves(expMoves, realMoves)
					})
				})
			})
			Context("and the pawn cannot capture", func() {
				It("returns non-capturing promotion moves", func() {
					board, err := BoardFromFEN("8/1P6/8/8/8/8/8/k1K5 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{7, 2})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_KNIGHT},
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_BISHOP},
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_ROOK},
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_QUEEN},
					}
					compareMoves(expMoves, realMoves)
				})
			})
		})
		When("the pawn is pinned to its king", func() {
			Context("and the pin is coming from a piece on the same file", func() {
				It("returns only moves that block the pin", func() {
					board, err := BoardFromFEN("k7/4r3/8/8/8/4P3/4K3/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
			Context("and the pin is coming from a piece on the same diagonal", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k7/8/1q6/8/8/4P3/5K2/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					Expect(realMoves).To(HaveLen(0))
				})
			})
			Context("and the pin is coming from the same rank", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k7/8/8/8/8/1r2PK2/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					Expect(realMoves).To(HaveLen(0))
				})
			})
		})
		When("the pawn can block a check", func() {
			It("returns the move that would block the check", func() {
				board, err := BoardFromFEN("k7/1b1B4/8/8/8/4PK2/8/8 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
			Context("and the pawn can capture", func() {
				It("only returns the non-capturing move that blocks the check", func() {
					board, err := BoardFromFEN("k7/1b1B4/8/8/3q4/4PK2/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
			Context("and the friendly king is under double check", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k7/1b1B4/8/4n3/8/4PK2/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					Expect(realMoves).To(HaveLen(0))
				})
			})
		})
		When("the board is a terminal board", func() {
			It("returns no moves", func() {
				board, _ := BoardFromFEN("k6p/8/8/8/8/8/8/RQ4K1 b - - 0 1")
				realMoves, err := GetLegalMovesForPawn(board, &Square{8, 8})
				Expect(err).ToNot(HaveOccurred())
				Expect(realMoves).To(HaveLen(0))
			})
		})
	})
	Describe("#GetLegalMovesForKnight", func() {
		When("the knight is in the bottom right corner", func() {
			It("returns only 2 moves on the board", func() {
				board, _ := BoardFromFEN("k7/1p6/8/8/8/8/8/6KN w - - 0 1")
				realMoves, err := GetLegalMovesForKnight(board, &Square{1, 8})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_KNIGHT, &Square{1, 8}, &Square{2, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{1, 8}, &Square{3, 7}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the knight is unimpeded in the middle of the board", func() {
			It("returns all possible knight moves", func() {
				board, _ := BoardFromFEN("k7/8/8/8/p3N3/8/8/6K1 w - - 0 1")
				realMoves, err := GetLegalMovesForKnight(board, &Square{4, 5})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_KNIGHT, &Square{4, 5}, &Square{6, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{4, 5}, &Square{6, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{4, 5}, &Square{5, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{4, 5}, &Square{5, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{4, 5}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{4, 5}, &Square{3, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{4, 5}, &Square{2, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{4, 5}, &Square{2, 6}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("different pieces occupy land squares", func() {
			It("returns only moves that capture enemy pieces", func() {
				board, _ := BoardFromFEN("k4r2/7N/5P2/6q1/p7/8/8/4K3 w - - 0 1")
				realMoves, err := GetLegalMovesForKnight(board, &Square{7, 8})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_KNIGHT, &Square{7, 8}, &Square{8, 6}, BLACK_ROOK, make([]*Square, 0), EMPTY},
					{WHITE_KNIGHT, &Square{7, 8}, &Square{5, 7}, BLACK_QUEEN, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("a move results in a discovered check", func() {
			It("returns the move with the rook's square checking the king", func() {
				board, _ := BoardFromFEN("5B2/7N/5P2/6q1/p7/7k/8/4K3 w - - 0 1")
				realMoves, err := GetLegalMovesForKnight(board, &Square{7, 8})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_KNIGHT, &Square{7, 8}, &Square{5, 7}, BLACK_QUEEN, []*Square{{5, 7}}, EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
			Context("and teh move is a discovered double check", func() {
				It("returns a move with both the knight and the rook squares checking the king", func() {
					board, _ := BoardFromFEN("5B1R/7N/5P2/6q1/p7/7k/8/4K3 w - - 0 1")
					realMoves, err := GetLegalMovesForKnight(board, &Square{7, 8})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_KNIGHT, &Square{7, 8}, &Square{5, 7}, BLACK_QUEEN, []*Square{{8, 8}, {5, 7}}, EMPTY},
					}
					compareMoves(expMoves, realMoves)
				})
			})
		})
		When("the board is a terminal board", func() {
			It("returns no moves", func() {
				board, _ := BoardFromFEN("k6n/8/8/8/8/8/8/RQ4K1 b - - 0 1")
				realMoves, err := GetLegalMovesForKnight(board, &Square{8, 8})
				Expect(err).ToNot(HaveOccurred())
				Expect(realMoves).To(HaveLen(0))
			})
		})
	})
	Describe("#GetLegalMovesForBishop", func() {
		When("the bishop is unimpeded in the middle of the board", func() {
			It("returns all legal bishop moves", func() {
				board, _ := BoardFromFEN("2K5/8/8/4B3/8/8/P7/5k2 w - - 0 1")
				realMoves, err := GetLegalMovesForBishop(board, &Square{5, 5})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_BISHOP, &Square{5, 5}, &Square{1, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{2, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{4, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{6, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{7, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{8, 8}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{6, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{7, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{8, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{4, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{3, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{5, 5}, &Square{2, 8}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("a friendly piece is blocking its path", func() {
			It("does not return moves on or beyond the friendly pieces", func() {
				board, _ := BoardFromFEN("2K5/8/8/8/3P4/8/8/B4k2 w - - 0 1")
				realMoves, err := GetLegalMovesForBishop(board, &Square{1, 1})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_BISHOP, &Square{1, 1}, &Square{2, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{1, 1}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("an enemy piece is blocking its path", func() {
			It("returns only moves up to and capturing the enemy piece", func() {
				board, _ := BoardFromFEN("2K5/8/8/8/3q4/8/8/B4k2 w - - 0 1")
				realMoves, err := GetLegalMovesForBishop(board, &Square{1, 1})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_BISHOP, &Square{1, 1}, &Square{2, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{1, 1}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_BISHOP, &Square{1, 1}, &Square{4, 4}, BLACK_QUEEN, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the move results in a double king check", func() {
			It("contains a move with two checking squares", func() {
				board, _ := BoardFromFEN("2K5/5R2/4P1N1/5B2/4P3/8/8/5k2 w - - 0 1")
				realMoves, err := GetLegalMovesForBishop(board, &Square{5, 6})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_BISHOP, &Square{5, 6}, &Square{4, 7}, EMPTY, []*Square{{7, 6}}, EMPTY},
					{WHITE_BISHOP, &Square{5, 6}, &Square{3, 8}, EMPTY, []*Square{{7, 6}, {3, 8}}, EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the board is a terminal board", func() {
			It("returns no moves", func() {
				board, _ := BoardFromFEN("k6b/8/8/8/8/8/R7/1Q4K1 b - - 0 1")
				realMoves, err := GetLegalMovesForBishop(board, &Square{8, 8})
				Expect(err).ToNot(HaveOccurred())
				Expect(realMoves).To(HaveLen(0))
			})
		})
	})
	Describe("#GetLegalMovesForRook", func() {
		When("the rook is unimpeded in the middle of the board", func() {
			It("returns all legal rook moves", func() {
				board, _ := BoardFromFEN("5K1k/7p/8/8/3R4/8/8/8 w - - 0 1")
				realMoves, err := GetLegalMovesForRook(board, &Square{4, 4})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_ROOK, &Square{4, 4}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{4, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{4, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{4, 8}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{4, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{4, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{4, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{1, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{2, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{3, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{5, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{6, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{7, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{4, 4}, &Square{8, 4}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the rook is impeded by both friendly and enemy pieces", func() {
			It("does not return friendly captures and moves beyond enemy pieces", func() {
				board, _ := BoardFromFEN("R1N5/7p/q7/8/8/4K3/8/6k1 w - - 0 1")
				realMoves, err := GetLegalMovesForRook(board, &Square{8, 1})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_ROOK, &Square{8, 1}, &Square{7, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{8, 1}, &Square{6, 1}, BLACK_QUEEN, make([]*Square, 0), EMPTY},
					{WHITE_ROOK, &Square{8, 1}, &Square{8, 2}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the move results in a check", func() {
			It("returns a move with a checking square", func() {
				board, _ := BoardFromFEN("R1N5/Q6p/8/8/8/4K3/8/1k6 w - - 0 1")
				realMoves, err := GetLegalMovesForRook(board, &Square{8, 1})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_ROOK, &Square{8, 1}, &Square{8, 2}, EMPTY, []*Square{{8, 2}}, EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the board is a terminal board", func() {
			It("returns no moves", func() {
				board, _ := BoardFromFEN("k6r/8/8/8/8/8/8/RQ4K1 b - - 0 1")
				realMoves, err := GetLegalMovesForRook(board, &Square{1, 1})
				Expect(err).ToNot(HaveOccurred())
				Expect(realMoves).To(HaveLen(0))
			})
		})
	})
	Describe("#GetLegalMovesForQueen", func() {
		When("the queen is unimpeded in the middle of the board", func() {
			It("returns all legal queen moves", func() {
				board, _ := BoardFromFEN("kr6/pp6/5K2/8/2Q5/8/8/8 w - - 0 1")
				realMoves, err := GetLegalMovesForQueen(board, &Square{4, 3})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_QUEEN, &Square{4, 3}, &Square{4, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{4, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{4, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{4, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{4, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{4, 8}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{1, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{2, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{5, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{6, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{7, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{8, 3}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{2, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{3, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{5, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{7, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{8, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{6, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{5, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{3, 4}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{2, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{4, 3}, &Square{1, 6}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the queen is impeded by both friendly and enemy pieces", func() {
			It("does not return friendly captures and moves beyond enemy pieces", func() {
				board, _ := BoardFromFEN("2K3pk/R5pp/8/8/8/2p5/8/Q2N4 w - - 0 1")
				realMoves, err := GetLegalMovesForQueen(board, &Square{1, 1})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_QUEEN, &Square{1, 1}, &Square{2, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{3, 3}, BLACK_PAWN, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{2, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{3, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{4, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{5, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{6, 1}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{1, 2}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_QUEEN, &Square{1, 1}, &Square{1, 3}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the move results in a check", func() {
			It("returns a move with a checking square", func() {
				board, _ := BoardFromFEN("2K3pk/R4PpN/5PQR/5PNP/8/2p5/8/3N4 w - - 0 1")
				realMoves, err := GetLegalMovesForQueen(board, &Square{6, 7})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_QUEEN, &Square{6, 7}, &Square{7, 7}, BLACK_PAWN, []*Square{{7, 7}}, EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the board is a terminal board", func() {
			It("returns no moves", func() {
				board, _ := BoardFromFEN("k6q/8/8/8/8/8/R7/1Q4K1 b - - 0 1")
				realMoves, err := GetLegalMovesForQueen(board, &Square{8, 8})
				Expect(err).ToNot(HaveOccurred())
				Expect(realMoves).To(HaveLen(0))
			})
		})
	})
	Describe("#GetLegalMovesForKing", func() {
		When("the king is unimpeded in the middle of the board", func() {
			It("returns all 8 king moves", func() {
				board, _ := BoardFromFEN("k7/p7/5K2/8/8/8/8/8 w - - 0 1")
				realMoves := GetLegalMovesForKing(board)
				expMoves := []Move{
					{WHITE_KING, &Square{6, 6}, &Square{7, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{6, 6}, &Square{7, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{6, 6}, &Square{7, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{6, 6}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{6, 6}, &Square{6, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{6, 6}, &Square{5, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{6, 6}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{6, 6}, &Square{5, 7}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the king is unimpeded on the corner of the board", func() {
			It("does not include moves off the board", func() {
				board, _ := BoardFromFEN("k7/p7/8/8/8/8/8/7K w - - 0 1")
				realMoves := GetLegalMovesForKing(board)
				expMoves := []Move{
					{WHITE_KING, &Square{1, 8}, &Square{1, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{1, 8}, &Square{2, 7}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{1, 8}, &Square{2, 8}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the king is two files away from the enemy king", func() {
			It("does not include moves that 'touch' the enemy king", func() {
				board, _ := BoardFromFEN("8/p7/8/2k1K3/8/8/8/8 w - - 0 1")
				realMoves := GetLegalMovesForKing(board)
				expMoves := []Move{
					{WHITE_KING, &Square{5, 5}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{5, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{5, 5}, &Square{6, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{5, 5}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
					{WHITE_KING, &Square{5, 5}, &Square{4, 6}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(expMoves, realMoves)
			})
		})
		When("the king has kingside castle rights", func() {
			Context("and a piece occupies a square between the king and kingside rook", func() {
				Context("and the king is white", func() {
					It("does not return a king move to castle kingside", func() {
						board, _ := BoardFromFEN("4k3/p7/8/8/8/8/3PPP2/3QK1NR w K - 0 1")
						realMoves := GetLegalMovesForKing(board)
						for _, realMove := range realMoves {
							Expect(realMove.EndSquare.EqualTo(&Square{1, 7})).To(BeFalse())
						}
					})
				})
				Context("and the king is black", func() {
					It("does not return a king move to castle kingside", func() {
						board, _ := BoardFromFEN("3pkb1r/3ppp2/8/8/8/8/8/3K4 b k - 0 1")
						realMoves := GetLegalMovesForKing(board)
						for _, realMove := range realMoves {
							Expect(realMove.EndSquare.EqualTo(&Square{8, 7})).To(BeFalse())
						}
					})
				})
			})
			Context("and all pieces are cleared between the king and kingside rook", func() {
				Context("and the king is white", func() {
					It("returns a king move to castle kingside", func() {
						board, _ := BoardFromFEN("4k3/p7/8/8/8/5N2/3PPP2/3QK2R w K - 0 1")
						realMoves := GetLegalMovesForKing(board)
						foundCastleMove := false
						for _, realMove := range realMoves {
							if realMove.EndSquare.EqualTo(&Square{1, 7}) {
								foundCastleMove = true
								break
							}
						}
						Expect(foundCastleMove).To(BeTrue())
					})
					Context("and an enemy rook stares between the king's start and end castle squares", func() {
						It("does not return a king move to castle kingside", func() {
							board, _ := BoardFromFEN("4k1r1/p7/8/8/8/5N2/3PPP2/3QK2R w K - 0 1")
							realMoves := GetLegalMovesForKing(board)
							for _, realMove := range realMoves {
								Expect(realMove.EndSquare.EqualTo(&Square{1, 7})).To(BeFalse())
							}
						})
					})
				})
				Context("and the king is black", func() {
					It("returns a king move to castle kingside", func() {
						board, _ := BoardFromFEN("3pk2r/3ppp2/3b4/8/8/8/8/3K4 b k - 0 1")
						realMoves := GetLegalMovesForKing(board)
						foundCastleMove := false
						for _, realMove := range realMoves {
							if realMove.EndSquare.EqualTo(&Square{8, 7}) {
								foundCastleMove = true
								break
							}
						}
						Expect(foundCastleMove).To(BeTrue())
					})
				})
			})
			Context("and the king is in check", func() {
				var board *Board
				BeforeEach(func() {
					board, _ = BoardFromFEN("4k2r/8/8/8/8/4Q3/3R4/2K5 b k - 0 1")
				})
				It("does not return a king move to castle kingside", func() {
					realMoves := GetLegalMovesForKing(board)
					foundCastleMove := false
					for _, realMove := range realMoves {
						if realMove.EndSquare.EqualTo(&Square{8, 7}) {
							foundCastleMove = true
							break
						}
					}
					Expect(foundCastleMove).To(BeFalse())
				})
			})
		})
		When("the king has queenside castle rights", func() {
			Context("and a piece occupies a square between the king and queenside rook", func() {
				Context("and the king is white", func() {
					It("does not return a king move to castle queenside", func() {
						board, _ := BoardFromFEN("4k3/8/2p5/8/8/8/3PPP2/RQ2KB2 w Q - 0 1")
						realMoves := GetLegalMovesForKing(board)
						for _, realMove := range realMoves {
							Expect(realMove.EndSquare.EqualTo(&Square{1, 3})).To(BeFalse())
						}
					})
				})
				Context("and the king is black", func() {
					It("does not return a king move to castle queenside", func() {
						board, _ := BoardFromFEN("r2qkb2/3ppp2/8/8/8/8/3PPP2/R2QKB2 b q - 1 1")
						realMoves := GetLegalMovesForKing(board)
						for _, realMove := range realMoves {
							Expect(realMove.EndSquare.EqualTo(&Square{8, 3})).To(BeFalse())
						}
					})
				})
			})
			Context("and all pieces are cleared between the king and queenside rook", func() {
				Context("and the king is white", func() {
					It("returns a king move to castle queenside", func() {
						board, _ := BoardFromFEN("r2qkb2/3ppp2/8/8/8/8/3PPP2/R3KB2 w Qq - 1 1")
						realMoves := GetLegalMovesForKing(board)
						foundCastleMove := false
						for _, realMove := range realMoves {
							if realMove.EndSquare.EqualTo(&Square{1, 3}) {
								foundCastleMove = true
								break
							}
						}
						Expect(foundCastleMove).To(BeTrue())
					})
					Context("and an enemy rook stares down between the king's start and end castle squares", func() {
						It("does not return a king move to castle queenside", func() {
							board, _ := BoardFromFEN("r2qkb2/3ppp2/8/8/8/8/3PPP2/R2QKB2 b q - 1 1")
							realMoves := GetLegalMovesForKing(board)
							for _, realMove := range realMoves {
								Expect(realMove.EndSquare.EqualTo(&Square{8, 3})).To(BeFalse())
							}
						})
					})
				})
				Context("and the king is black", func() {
					It("returns a move to castle queenside", func() {
						board, _ := BoardFromFEN("r3kb2/3ppp2/8/8/8/8/3PPP2/R3KB2 b Qq - 1 1")
						realMoves := GetLegalMovesForKing(board)
						foundCastleMove := false
						for _, realMove := range realMoves {
							if realMove.EndSquare.EqualTo(&Square{8, 3}) {
								foundCastleMove = true
								break
							}
						}
						Expect(foundCastleMove).To(BeTrue())
					})
				})
			})
		})
		When("the board is a terminal board", func() {
			It("returns no moves", func() {
				board, _ := BoardFromFEN("k7/8/8/8/8/8/8/RQ4K1 b - - 0 1")
				realMoves := GetLegalMovesForKing(board)
				Expect(realMoves).To(HaveLen(0))
			})
		})
	})
	Describe("#IsLegalMove", func() {
		When("the move is a legal move", func() {
			It("returns true", func() {
				board := GetInitBoard()
				move := &Move{WHITE_PAWN, &Square{2, 2}, &Square{4, 2}, EMPTY, make([]*Square, 0), EMPTY}
				Expect(IsLegalMove(board, move)).To(BeTrue())
			})
		})
		When("the move is not a legal move", func() {
			It("returns false", func() {
				board := GetInitBoard()
				move := &Move{WHITE_PAWN, &Square{2, 2}, &Square{5, 2}, EMPTY, make([]*Square, 0), EMPTY}
				Expect(IsLegalMove(board, move)).To(BeFalse())
			})
			Context("and the move would be legal except that the board is in a terminal board state", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("8/4k3/8/8/5N2/6K1/8/8 w - - 0 1")
					move := &Move{WHITE_KING, &Square{3, 7}, &Square{4, 7}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(IsLegalMove(board, move)).To(BeFalse())
				})
			})
		})
	})
	Describe("#HasLegalMove", func() {
		var board *Board
		When("the board has legal moves", func() {
			BeforeEach(func() {
				board = GetInitBoard()
			})
			It("returns true", func() {
				Expect(HasLegalMove(board)).To(BeTrue())
			})
		})
		When("the board does not have legal moves", func() {
			BeforeEach(func() {
				var boardErr error
				board, boardErr = BoardFromFEN("8/8/2k2K2/5N2/8/8/8/8 w - - 0 1")
				Expect(boardErr).To(Succeed())
			})
			It("returns false", func() {
				Expect(HasLegalMove(board)).To(BeFalse())
			})
		})
	})
	Describe("#GetLegalMovesByOrigin", func() {
		var board *Board
		BeforeEach(func() {
			board = GetInitBoard()
		})
		DescribeTable("returns the correct length of moves for a square",
			func(r int, c int, expMoveLen int) {
				square := Square{uint8(r), uint8(c)}
				moves, movesErr := GetLegalMovesFromOrigin(board, &square)
				Expect(movesErr).To(Succeed())
				Expect(moves).To(HaveLen(expMoveLen))
			},
			Entry("A-pawn has two moves", 2, 1, 2),
			Entry("B-pawn has two moves", 2, 2, 2),
			Entry("C-pawn has two moves", 2, 3, 2),
			Entry("D-pawn has two moves", 2, 4, 2),
			Entry("E-pawn has two moves", 2, 5, 2),
			Entry("F-pawn has two moves", 2, 6, 2),
			Entry("G-pawn has two moves", 2, 7, 2),
			Entry("H-pawn has two moves", 2, 8, 2),
			Entry("queens rook has no moves", 1, 1, 0),
			Entry("queens knight has two moves", 1, 2, 2),
			Entry("queens bishop has no moves", 1, 3, 0),
			Entry("queen has no moves", 1, 4, 0),
			Entry("king has no moves", 1, 5, 0),
			Entry("king's bishop has no moves", 1, 6, 0),
			Entry("king's knight has two moves", 1, 7, 2),
			Entry("king's rook has no moves", 1, 8, 0))
	})
	Describe("#GetBoardFromMove", func() {
		var board *Board
		var move Move
		It("does not affect the input board", func() {
			board = GetInitBoard()
			move = Move{
				WHITE_KNIGHT,
				&Square{1, 2},
				&Square{3, 3},
				EMPTY,
				make([]*Square, 0),
				EMPTY,
			}
			boardJson, jsonParseErr := json.Marshal(board)
			Expect(jsonParseErr).To(Succeed())

			_ = GetBoardFromMove(board, &move)

			newBoardJson, jsonParseErr := json.Marshal(board)
			Expect(jsonParseErr).To(Succeed())
			Expect(string(newBoardJson)).To(Equal(string(boardJson)))
			fmt.Println(string(newBoardJson))
		})
		When("the move is a capture", func() {
			BeforeEach(func() {
				initBoard, _ := BoardFromFEN("k7/p2n4/8/8/6B1/8/8/7K w - - 0 1")
				initBoard.RepetitionsByMiniFEN["asdf"] = 2
				move = Move{WHITE_BISHOP, &Square{4, 7}, &Square{7, 4}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				board = GetBoardFromMove(initBoard, &move)
			})
			It("resets the half move clock counter", func() {
				Expect(board.HalfMoveClockCount).To(Equal(uint8(0)))
			})
			It("replaces the capturing piece with the capture piece", func() {
				Expect(board.GetPieceOnSquare(&Square{4, 7})).To(Equal(EMPTY))
				Expect(board.GetPieceOnSquare(&Square{7, 4})).To(Equal(WHITE_BISHOP))
			})
			It("adds the miniFEN to the board counter map", func() {
				miniFEN := "k7/p2B4/8/8/8/8/8/7K b - -"
				miniFENRepetitions, ok := board.RepetitionsByMiniFEN[miniFEN]
				Expect(ok).To(BeTrue())
				Expect(miniFENRepetitions).To(Equal(uint8(1)))
			})
			It("clears the repetition counter map", func() {
				_, ok := board.RepetitionsByMiniFEN["asdf"]
				Expect(ok).To(BeFalse())
			})
			When("the capture is an en passant move", func() {
				BeforeEach(func() {
					initBoard, _ := BoardFromFEN("k7/8/4pP2/8/8/8/8/7K w - e7 0 1")
					move = Move{WHITE_PAWN, &Square{6, 6}, &Square{7, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY}
					board = GetBoardFromMove(initBoard, &move)
				})
				It("moves the capturing piece to the en passant square and removes the en passanted pawn", func() {
					Expect(board.GetPieceOnSquare(&Square{6, 6})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{6, 5})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{7, 5})).To(Equal(WHITE_PAWN))
				})
			})
			When("the move is a capture on a 'castle-able' rook", func() {
				BeforeEach(func() {
					initBoard, _ := BoardFromFEN("4k2r/8/6N1/8/8/8/8/R3K3 w k - 0 1")
					move = Move{WHITE_KNIGHT, &Square{6, 7}, &Square{8, 8}, BLACK_ROOK, make([]*Square, 0), EMPTY}
					board = GetBoardFromMove(initBoard, &move)
				})
				It("revokes the enemy king castle right on that side", func() {
					Expect(board.CanBlackCastleKingside).To(BeFalse())
				})
			})
		})
		When("the move is not a capture", func() {
			BeforeEach(func() {
				initBoard, _ := BoardFromFEN("k7/4r3/8/8/8/8/5N2/7K w - - 0 1")
				initBoard.RepetitionsByMiniFEN["asdf"] = 2
				move = Move{WHITE_KNIGHT, &Square{2, 6}, &Square{3, 4}, EMPTY, make([]*Square, 0), EMPTY}
				board = GetBoardFromMove(initBoard, &move)
			})
			It("moves the piece", func() {
				Expect(board.GetPieceOnSquare(&Square{2, 6})).To(Equal(EMPTY))
				Expect(board.GetPieceOnSquare(&Square{3, 4})).To(Equal(WHITE_KNIGHT))
			})
			It("increments the half move counter", func() {
				Expect(board.HalfMoveClockCount).To(Equal(uint8(1)))
			})
			It("adds the miniFEN to the repetition counter map", func() {
				miniFEN := "k7/4r3/8/8/8/3N4/8/7K b - -"
				miniFENRepetitions, ok := board.RepetitionsByMiniFEN[miniFEN]
				Expect(ok).To(BeTrue())
				Expect(miniFENRepetitions).To(Equal(uint8(1)))
			})
			It("does not clear the repetition counter map", func() {
				miniFENRepetitions, ok := board.RepetitionsByMiniFEN["asdf"]
				Expect(ok).To(BeTrue())
				Expect(miniFENRepetitions).To(Equal(uint8(2)))
			})
			Context("and the pawn is being upgraded", func() {
				BeforeEach(func() {
					initBoard, _ := BoardFromFEN("8/5P2/8/8/8/1kb5/8/7K w - - 0 1")
					move := Move{WHITE_PAWN, &Square{7, 6}, &Square{8, 6}, EMPTY, make([]*Square, 0), WHITE_ROOK}
					board = GetBoardFromMove(initBoard, &move)
				})
				It("converts the pawn to its upgraded piece", func() {
					Expect(board.GetPieceOnSquare(&Square{8, 6})).To(Equal(WHITE_ROOK))
					Expect(board.GetPieceOnSquare(&Square{7, 6})).To(Equal(EMPTY))
				})
			})
			Context("and the move allows for en passant", func() {
				BeforeEach(func() {
					initBoard, _ := BoardFromFEN("8/8/8/8/7K/8/4P3/k7 w - - 0 1")
					move := Move{WHITE_PAWN, &Square{2, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
					board = GetBoardFromMove(initBoard, &move)
				})
				It("sets the en passant square", func() {
					Expect(board.OptEnPassantSquare).ToNot(BeNil())
					Expect(board.OptEnPassantSquare.EqualTo(&Square{3, 5}))
				})
			})
		})
		When("its white's move", func() {
			BeforeEach(func() {
				initBoard, _ := BoardFromFEN("8/8/8/8/7K/8/4P3/k7 w - - 0 1")
				move = Move{WHITE_PAWN, &Square{2, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
				board = GetBoardFromMove(initBoard, &move)
			})
			It("sets the next turn to black", func() {
				Expect(board.IsWhiteTurn).To(BeFalse())
			})
		})
		When("its black's move", func() {
			BeforeEach(func() {
				initBoard, _ := BoardFromFEN("8/1r4k1/8/8/7K/8/4P3/8 b - - 0 1")
				move = Move{BLACK_ROOK, &Square{7, 2}, &Square{2, 2}, EMPTY, make([]*Square, 0), EMPTY}
				board = GetBoardFromMove(initBoard, &move)
			})
			It("increments the full move counter", func() {
				Expect(board.FullMoveCount).To(Equal(uint16(2)))
			})
			It("moves the piece", func() {
				Expect(board.GetPieceOnSquare(&Square{7, 2})).To(Equal(EMPTY))
				Expect(board.GetPieceOnSquare(&Square{2, 2})).To(Equal(BLACK_ROOK))
			})
			It("updates the next turn to white", func() {
				Expect(board.IsWhiteTurn).To(BeTrue())
			})
			It("increments the half move clock count", func() {
				Expect(board.HalfMoveClockCount).To(Equal(uint8(1)))
			})
		})
		When("a rook moves, revoking a castling right", func() {
			When("the rook is white", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")
					board.RepetitionsByMiniFEN["asdf"] = 2
					Expect(board.CanWhiteCastleKingside).To(BeTrue())
					Expect(board.CanWhiteCastleQueenside).To(BeTrue())
				})
				Context("and the rook is the kingside rook", func() {
					BeforeEach(func() {
						move = Move{WHITE_ROOK, &Square{1, 8}, &Square{1, 7}, EMPTY, make([]*Square, 0), EMPTY}
						board = GetBoardFromMove(board, &move)
					})
					It("revokes the white kingside castle right", func() {
						Expect(board.CanWhiteCastleKingside).To(BeFalse())
					})
					It("clears the repetition counter map", func() {
						_, ok := board.RepetitionsByMiniFEN["asdf"]
						Expect(ok).To(BeFalse())
					})
				})
				Context("and the rook is the queenside rook", func() {
					BeforeEach(func() {
						move = Move{WHITE_ROOK, &Square{1, 1}, &Square{1, 2}, EMPTY, make([]*Square, 0), EMPTY}
						board = GetBoardFromMove(board, &move)
					})
					It("revokes the white queenside castle right", func() {
						Expect(board.CanWhiteCastleQueenside).To(BeFalse())
					})
					It("clears the repetition counter map", func() {
						_, ok := board.RepetitionsByMiniFEN["asdf"]
						Expect(ok).To(BeFalse())
					})
				})
			})
			When("the rook is black", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")
					board.RepetitionsByMiniFEN["asdf"] = 2
				})
				Context("and the rook is the kingside rook", func() {
					BeforeEach(func() {
						Expect(board.CanBlackCastleKingside).To(BeTrue())
						move = Move{BLACK_ROOK, &Square{8, 8}, &Square{8, 7}, EMPTY, make([]*Square, 0), EMPTY}
						board = GetBoardFromMove(board, &move)
					})
					It("revokes the black kingside castle right", func() {
						Expect(board.CanBlackCastleKingside).To(BeFalse())
					})
					It("clears the repetition counter map", func() {
						_, ok := board.RepetitionsByMiniFEN["asdf"]
						Expect(ok).To(BeFalse())
					})
				})
				Context("and the rook is the queenside rook", func() {
					BeforeEach(func() {
						Expect(board.CanBlackCastleQueenside).To(BeTrue())
						move = Move{BLACK_ROOK, &Square{8, 1}, &Square{8, 2}, EMPTY, make([]*Square, 0), EMPTY}
						board = GetBoardFromMove(board, &move)
					})
					It("revokes the black queenside castle right", func() {
						Expect(board.CanBlackCastleQueenside).To(BeFalse())
					})
					It("clears the repetition counter map", func() {
						_, ok := board.RepetitionsByMiniFEN["asdf"]
						Expect(ok).To(BeFalse())
					})
				})
			})
		})
		When("a king moves with both castle rights", func() {
			Context("the king is white", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")
					board.RepetitionsByMiniFEN["asdf"] = 2
					Expect(board.CanWhiteCastleKingside).To(BeTrue())
					Expect(board.CanWhiteCastleQueenside).To(BeTrue())
					move = Move{WHITE_KING, &Square{1, 5}, &Square{2, 4}, EMPTY, make([]*Square, 0), EMPTY}
					board = GetBoardFromMove(board, &move)
				})
				It("revokes both of white's castle rights", func() {
					Expect(board.CanWhiteCastleKingside).To(BeFalse())
					Expect(board.CanWhiteCastleQueenside).To(BeFalse())
				})
				It("clears the repetition counter map", func() {
					_, ok := board.RepetitionsByMiniFEN["asdf"]
					Expect(ok).To(BeFalse())
				})
			})
			Context("the king is black", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")
					Expect(board.CanBlackCastleKingside).To(BeTrue())
					Expect(board.CanBlackCastleQueenside).To(BeTrue())
					move = Move{BLACK_KING, &Square{8, 5}, &Square{8, 6}, EMPTY, make([]*Square, 0), EMPTY}
					board = GetBoardFromMove(board, &move)
				})
				It("revokes both of black's castle rights", func() {
					Expect(board.CanBlackCastleKingside).To(BeFalse())
					Expect(board.CanBlackCastleQueenside).To(BeFalse())
				})
				It("clears the repetition counter map", func() {
					_, ok := board.RepetitionsByMiniFEN["asdf"]
					Expect(ok).To(BeFalse())
				})
			})
		})
		When("the moves is castles", func() {
			Context("and white castles kingside", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")
					board.RepetitionsByMiniFEN["asdf"] = 2
					move = Move{WHITE_KING, &Square{1, 5}, &Square{1, 7}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(board.GetPieceOnSquare(&Square{1, 5})).To(Equal(WHITE_KING))
					Expect(board.GetPieceOnSquare(&Square{1, 6})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{1, 7})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{1, 8})).To(Equal(WHITE_ROOK))
				})
				It("re-positions the white king and king's rook", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.GetPieceOnSquare(&Square{1, 5})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{1, 6})).To(Equal(WHITE_ROOK))
					Expect(board.GetPieceOnSquare(&Square{1, 7})).To(Equal(WHITE_KING))
					Expect(board.GetPieceOnSquare(&Square{1, 8})).To(Equal(EMPTY))
				})
				It("revokes castling rights for the active player", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.CanWhiteCastleKingside).To(BeFalse())
					Expect(board.CanWhiteCastleQueenside).To(BeFalse())
				})
				It("clears the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					_, ok := board.RepetitionsByMiniFEN["asdf"]
					Expect(ok).To(BeFalse())
				})
				It("adds the current miniFEN to the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					miniFEN := board.ToMiniFEN()
					miniFENRepetitions, ok := board.RepetitionsByMiniFEN[miniFEN]
					Expect(ok).To(BeTrue())
					Expect(miniFENRepetitions).To(Equal(uint8(1)))
				})
			})
			Context("and white castles queenside", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")
					move = Move{WHITE_KING, &Square{1, 5}, &Square{1, 3}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(board.GetPieceOnSquare(&Square{1, 5})).To(Equal(WHITE_KING))
					Expect(board.GetPieceOnSquare(&Square{1, 4})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{1, 3})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{1, 1})).To(Equal(WHITE_ROOK))
				})
				It("re-positions the white king and queen's rook", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.GetPieceOnSquare(&Square{1, 5})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{1, 4})).To(Equal(WHITE_ROOK))
					Expect(board.GetPieceOnSquare(&Square{1, 3})).To(Equal(WHITE_KING))
					Expect(board.GetPieceOnSquare(&Square{1, 1})).To(Equal(EMPTY))
				})
				It("revokes castling rights for the active player", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.CanWhiteCastleKingside).To(BeFalse())
					Expect(board.CanWhiteCastleQueenside).To(BeFalse())
				})
				It("clears the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					_, ok := board.RepetitionsByMiniFEN["asdf"]
					Expect(ok).To(BeFalse())
				})
				It("adds the current miniFEN to the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					miniFEN := board.ToMiniFEN()
					miniFENRepetitions, ok := board.RepetitionsByMiniFEN[miniFEN]
					Expect(ok).To(BeTrue())
					Expect(miniFENRepetitions).To(Equal(uint8(1)))
				})
			})
			Context("and black castles kingside", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")
					move = Move{BLACK_KING, &Square{8, 5}, &Square{8, 7}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(board.GetPieceOnSquare(&Square{8, 5})).To(Equal(BLACK_KING))
					Expect(board.GetPieceOnSquare(&Square{8, 6})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{8, 7})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{8, 8})).To(Equal(BLACK_ROOK))
				})
				It("re-positions the black king and king's rook", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.GetPieceOnSquare(&Square{8, 5})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{8, 6})).To(Equal(BLACK_ROOK))
					Expect(board.GetPieceOnSquare(&Square{8, 7})).To(Equal(BLACK_KING))
					Expect(board.GetPieceOnSquare(&Square{8, 8})).To(Equal(EMPTY))
				})
				It("revokes castling rights for the active player", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.CanBlackCastleKingside).To(BeFalse())
					Expect(board.CanBlackCastleQueenside).To(BeFalse())
				})
				It("clears the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					_, ok := board.RepetitionsByMiniFEN["asdf"]
					Expect(ok).To(BeFalse())
				})
				It("adds the current miniFEN to the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					miniFEN := board.ToMiniFEN()
					miniFENRepetitions, ok := board.RepetitionsByMiniFEN[miniFEN]
					Expect(ok).To(BeTrue())
					Expect(miniFENRepetitions).To(Equal(uint8(1)))
				})
			})
			Context("and black castles queenside", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1")
					move = Move{BLACK_KING, &Square{8, 5}, &Square{8, 3}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(board.GetPieceOnSquare(&Square{8, 5})).To(Equal(BLACK_KING))
					Expect(board.GetPieceOnSquare(&Square{8, 4})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{8, 3})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{8, 1})).To(Equal(BLACK_ROOK))
				})
				It("re-positions the black king and queen's rook", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.GetPieceOnSquare(&Square{8, 5})).To(Equal(EMPTY))
					Expect(board.GetPieceOnSquare(&Square{8, 4})).To(Equal(BLACK_ROOK))
					Expect(board.GetPieceOnSquare(&Square{8, 3})).To(Equal(BLACK_KING))
					Expect(board.GetPieceOnSquare(&Square{8, 1})).To(Equal(EMPTY))
				})
				It("revokes castling rights for the active player", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.CanBlackCastleKingside).To(BeFalse())
					Expect(board.CanBlackCastleQueenside).To(BeFalse())
				})
				It("clears the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					_, ok := board.RepetitionsByMiniFEN["asdf"]
					Expect(ok).To(BeFalse())
				})
				It("adds the current miniFEN to the repetition counter map", func() {
					board = GetBoardFromMove(board, &move)
					miniFEN := board.ToMiniFEN()
					miniFENRepetitions, ok := board.RepetitionsByMiniFEN[miniFEN]
					Expect(ok).To(BeTrue())
					Expect(miniFENRepetitions).To(Equal(uint8(1)))
				})
			})
		})
		When("the resulting board is repeated 3 times", func() {
			BeforeEach(func() {
				board, _ = BoardFromFEN("8/8/4k3/8/3K4/4P3/8/8 w - - 0 1")
				move = Move{WHITE_KING, &Square{4, 4}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
				board.RepetitionsByMiniFEN["8/8/4k3/8/4K3/4P3/8/8 b - -"] = 2
			})
			It("increments the board repetitions counter map", func() {
				board = GetBoardFromMove(board, &move)
				miniFEN := "8/8/4k3/8/4K3/4P3/8/8 b - -"
				miniFENRepetitions, ok := board.RepetitionsByMiniFEN[miniFEN]
				Expect(ok).To(BeTrue())
				Expect(miniFENRepetitions).To(Equal(uint8(3)))
			})
			It("results in a terminal draw state", func() {
				board = GetBoardFromMove(board, &move)
				Expect(board.Result).To(Equal(BOARD_RESULT_DRAW_BY_THREEFOLD_REPETITION))
			})
		})
		When("the resulting board is a stalemate", func() {
			BeforeEach(func() {
				board, _ = BoardFromFEN("5Q2/8/8/8/8/8/4K3/7k w - - 0 1")
				move = Move{WHITE_QUEEN, &Square{8, 6}, &Square{2, 6}, EMPTY, make([]*Square, 0), EMPTY}
			})
			It("results in a terminal draw board", func() {
				board = GetBoardFromMove(board, &move)
				Expect(board.Result).To(Equal(BOARD_RESULT_DRAW_BY_STALEMATE))
			})
		})
		When("the move violates the 50-move rule", func() {
			BeforeEach(func() {
				board = GetInitBoard()
				board.HalfMoveClockCount = 49
				move = Move{WHITE_KNIGHT, &Square{1, 2}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY}
			})
			It("results in a terminal draw board", func() {
				board = GetBoardFromMove(board, &move)
				Expect(board.Result).To(Equal(BOARD_RESULT_DRAW_BY_FIFTY_MOVE_RULE))
			})
		})
		When("the remaining material forces a draw", func() {
			BeforeEach(func() {
				board, _ = BoardFromFEN("8/8/8/5N2/1K6/6q1/8/7k w - - 1 1")
				Expect(board.Result).To(Equal(BOARD_RESULT_IN_PROGRESS))
				move = Move{WHITE_KNIGHT, &Square{5, 6}, &Square{3, 7}, BLACK_QUEEN, []*Square{{3, 7}}, EMPTY}
			})
			It("results in a terminal draw board", func() {
				board = GetBoardFromMove(board, &move)
				Expect(board.Result).To(Equal(BOARD_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL))
			})
		})
		When("the move is a checkmate", func() {
			Context("and the checkmating player is white", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("k1K5/8/8/8/8/1Q6/8/8 w - - 1 1")
					move = Move{WHITE_QUEEN, &Square{3, 2}, &Square{7, 2}, EMPTY, []*Square{{7, 2}}, EMPTY}
				})
				It("results in a terminal board", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.Result).To(Equal(BOARD_RESULT_WHITE_WINS_BY_CHECKMATE))
				})
			})
			Context("and the checkmating player is black", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("8/8/8/B4N2/1K6/R7/4q3/2r4k b - - 1 1")
					move = Move{BLACK_QUEEN, &Square{2, 5}, &Square{4, 3}, EMPTY, []*Square{{4, 3}}, EMPTY}
				})
				It("results in a terminal board", func() {
					board = GetBoardFromMove(board, &move)
					Expect(board.Result).To(Equal(BOARD_RESULT_BLACK_WINS_BY_CHECKMATE))
				})
			})
			When("the previous move was a double pawn move", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("rnbqkbnr/ppppp2p/5p2/6p1/3PP3/8/PPP2PPP/RNBQKBNR w KQkq g6 1 3")
				})
				It("resets the en passant square", func() {
					move = Move{WHITE_QUEEN, &Square{1, 4}, &Square{5, 8}, EMPTY, []*Square{{5, 8}}, EMPTY}
					board = GetBoardFromMove(board, &move)
					Expect(board.Result).To(Equal(BOARD_RESULT_WHITE_WINS_BY_CHECKMATE))
					Expect(board.OptEnPassantSquare).To(BeNil())
				})
			})
		})
		When("the last move was a pawn double jump", func() {
			Context("and this move is a pawn single jump", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1")
					Expect(board.OptEnPassantSquare).ToNot(BeNil())
					move = Move{BLACK_PAWN, &Square{7, 5}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY}
				})
				It("resets the en passant square", func() {
					newBoard := GetBoardFromMove(board, &move)
					Expect(newBoard.OptEnPassantSquare).To(BeNil())
				})
			})
		})
		When("the move updates castling rights", func() {
			BeforeEach(func() {
				initBoard, _ := BoardFromFEN("4k2r/8/8/8/8/8/8/R3K2R w KQ - 1 1")
				move = Move{WHITE_ROOK, &Square{1, 8}, &Square{2, 8}, EMPTY, make([]*Square, 0), EMPTY}
				board = GetBoardFromMove(initBoard, &move)
			})
			It("resets the half move clock", func() {
				Expect(board.HalfMoveClockCount).To(Equal(uint8(0)))
			})
		})
	})
	Describe("#UpdatePiecesFromMove", func() {
		var board *Board
		var boardBuilder *BoardBuilder
		var move *Move
		BeforeEach(func() {
			board, _ = BoardFromFEN("r1bqkb1r/ppp2ppp/2n2n2/1N1pp3/3P1B2/8/PPP1PPPP/R2QKBNR w KQkq - 0 1")
			boardBuilder = NewBoardBuilder().FromBoard(board)
			move = &Move{WHITE_PAWN, &Square{4, 4}, &Square{5, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY}
		})
		It("moves the piece", func() {
			UpdatePiecesFromMove(board, boardBuilder, move)
			board := boardBuilder.Build()
			Expect(board.GetPieceOnSquare(&Square{4, 4})).To(Equal(EMPTY))
			Expect(board.GetPieceOnSquare(&Square{5, 5})).To(Equal(WHITE_PAWN))
		})
		It("updates the cached material count", func() {
			UpdatePiecesFromMove(board, boardBuilder, move)
			board := boardBuilder.Build()
			Expect(board.ComputeMaterialCount().BlackPawnCount).To(Equal(uint8(7)))
		})
	})
})
