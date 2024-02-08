package chess_test

import (
	. "github.com/CameronHonis/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Board", func() {
	Describe("::IsForcedDrawByMaterial", func() {
		When("the remaining material forces a draw", func() {
			When("only the kings are on the board", func() {
				It("returns true", func() {
					board, _ := BoardFromFEN("8/8/8/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
				})
			})
			When("only one bishop exists", func() {
				Context("and its a white bishop", func() {
					It("returns true", func() {
						board, _ := BoardFromFEN("8/4B3/8/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
					})
				})
				Context("and its a black bishop", func() {
					It("returns true", func() {
						board, _ := BoardFromFEN("8/4b3/8/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
					})
				})
			})
			When("both players have one bishop", func() {
				It("returns true", func() {
					board, _ := BoardFromFEN("8/4B3/8/8/7K/8/3b4/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
				})
			})
			When("a player has two same colored bishops", func() {
				Context("and the bishops are white", func() {
					It("returns true", func() {
						board, _ := BoardFromFEN("8/4B3/3B4/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
					})
				})
				Context("and the bishops are black", func() {
					It("returns true", func() {
						board, _ := BoardFromFEN("8/5b2/4b3/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
					})
				})
			})
			When("only one knight exists", func() {
				Context("and the knight is white", func() {
					It("returns true", func() {
						board, _ := BoardFromFEN("8/8/8/8/4N2K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
					})
				})
				Context("and the knight is black", func() {
					It("returns true", func() {
						board, _ := BoardFromFEN("8/8/4n3/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
					})
				})
			})
			When("both players have only one knight", func() {
				It("returns true", func() {
					board, _ := BoardFromFEN("5n2/8/8/8/4N2K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeTrue())
				})
			})
		})
		When("the remaining material does not force a draw", func() {
			When("only one rook exists", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("8/2R5/8/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeFalse())
				})
			})
			When("only one queen exists", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("8/8/4q3/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeFalse())
				})
			})
			When("only one pawn exists", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("8/8/8/6P1/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeFalse())
				})
			})
			When("a player has two different colored bishops", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("8/4B3/4B3/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeFalse())
				})
			})
			When("only one player has a bishop and knight", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("8/4B3/4N3/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeFalse())
				})
			})
			When("only one player has two knights", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("8/8/4n3/4n3/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDrawByMaterial()).To(BeFalse())
				})
			})
		})
	})
	Describe("::IsWhiteCheckmated", func() {
		When("neither player is checkmated", func() {
			When("the board is the initial board", func() {
				It("returns false", func() {
					board := GetInitBoard()
					Expect(board.IsWhiteCheckmated()).To(BeFalse())
				})
			})
			When("white is in check, but has legal moves to escape", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("4K3/4q3/8/8/8/5k2/8/8 w - - 0 1")
					Expect(board.IsWhiteCheckmated()).To(BeFalse())
				})
			})
			When("white is in stalemate", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("K7/8/1q6/8/8/5k2/8/8 w - - 0 1")
					Expect(board.IsWhiteCheckmated()).To(BeFalse())
				})
			})
		})
		When("white is checkmated", func() {
			It("returns true", func() {
				board, _ := BoardFromFEN("K7/1q6/2k5/8/8/8/8/8 w - - 0 1")
				Expect(board.IsWhiteCheckmated()).To(BeTrue())
			})
		})
		When("black is checkmated", func() {
			It("returns false", func() {
				board, _ := BoardFromFEN("k7/1Q6/2K5/8/8/8/8/8 w - - 0 1") // note it is white's turn
				Expect(board.IsWhiteCheckmated()).To(BeFalse())
			})
		})
	})
	Describe("::IsBlackCheckmated", func() {
		When("neither player is checkmated", func() {
			When("the board is the initial board", func() {
				It("returns false", func() {
					board := GetInitBoard()
					Expect(board.IsBlackCheckmated()).To(BeFalse())
				})
			})
			When("black is in check, but has legal moves to escape", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("4k3/4Q3/8/8/8/5K2/8/8 b - - 0 1")
					Expect(board.IsBlackCheckmated()).To(BeFalse())
				})
			})
			When("black is in stalemate", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("k7/8/1Q6/8/8/5K2/8/8 b - - 0 1")
					Expect(board.IsBlackCheckmated()).To(BeFalse())
				})
			})
		})
		When("black is checkmated", func() {
			It("returns true", func() {
				board, _ := BoardFromFEN("k7/1Q6/2K5/8/8/8/8/8 b - - 0 1")
				Expect(board.IsBlackCheckmated()).To(BeTrue())
			})
		})
		When("white is checkmated", func() {
			It("returns false", func() {
				board, _ := BoardFromFEN("K7/1q6/2k5/8/8/8/8/8 b - - 0 1") // note it is black's turn
				Expect(board.IsBlackCheckmated()).To(BeFalse())
			})
		})
	})
	Describe("::IsDrawByStalemate", func() {
		When("there is no stalemate", func() {
			When("the board is the initial board", func() {
				It("returns false", func() {
					board := GetInitBoard()
					Expect(board.IsDrawByStalemate()).To(BeFalse())
				})
			})
			When("white is in check", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("K7/q7/3k4/8/8/8/8/8 w - - 0 1")
					Expect(board.IsDrawByStalemate()).To(BeFalse())
				})
			})
			When("the board is terminal due to the fifty move rule", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("K7/8/3k1r2/8/8/8/8/8 w - - 50 102")
					Expect(board.IsDrawByStalemate()).To(BeFalse())
				})
			})
		})
		When("white is stalemated", func() {
			It("returns true", func() {
				board, _ := BoardFromFEN("K7/8/1qk5/8/8/8/8/8 w - - 0 1")
				Expect(board.IsDrawByStalemate()).To(BeTrue())
			})
		})
		When("black is stalemated", func() {
			It("returns true", func() {
				board, _ := BoardFromFEN("k7/8/1QK5/8/8/8/8/8 b - - 0 1")
				Expect(board.IsDrawByStalemate()).To(BeTrue())
			})
		})
	})
	Describe("::IsDrawByFiftyMoveRule", func() {
		When("the fifty move rule has not been reached", func() {
			It("returns false", func() {
				board := GetInitBoard()
				Expect(board.IsDrawByFiftyMoveRule()).To(BeFalse())
			})
		})
		When("the fifty move rule has been reached", func() {
			It("returns true", func() {
				board, _ := BoardFromFEN("K7/8/3k1r2/8/8/8/8/8 w - - 50 102")
				Expect(board.IsDrawByFiftyMoveRule()).To(BeTrue())
			})
		})
	})
	Describe("::IsDrawByThreefoldRepetition", func() {
		When("the board is terminal", func() {
			When("white is checkmated", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("K7/1q6/2k2r2/8/8/8/8/8 w - - 0 1")
					Expect(board.IsDrawByThreefoldRepetition()).To(BeFalse())
				})
			})
			When("black is checkmated", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("k7/1Q6/2K2R2/8/8/8/8/8 w - - 0 1")
					Expect(board.IsDrawByThreefoldRepetition()).To(BeFalse())
				})
			})
			When("the board is a draw by stalemate", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("K7/8/k7/8/8/8/8/1r6 w - - 0 1")
					Expect(board.IsDrawByThreefoldRepetition()).To(BeFalse())
				})
			})
			When("the board is a draw by the fifty move rule", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("K7/8/8/8/8/5k2/6q1/8 w - - 50 102")
					Expect(board.IsDrawByThreefoldRepetition()).To(BeFalse())
				})
			})
			When("the board is a draw by insufficient material", func() {
				It("returns false", func() {
					board, _ := BoardFromFEN("K7/8/2k5/8/8/8/8/8 w - - 0 1")
					Expect(board.IsDrawByThreefoldRepetition()).To(BeFalse())
				})
			})
			When("the board is a draw by threefold repetition", func() {
				var board *Board
				BeforeEach(func() {
					board = GetInitBoard()
					board.IsTerminal = true // pretend both players shuffle around the same knights 3 times
				})
				It("returns true", func() {
					Expect(board.IsDrawByFiftyMoveRule())
				})
			})
		})
		When("the board is not terminal", func() {
			It("returns false", func() {
				board := GetInitBoard()
				Expect(board.IsDrawByThreefoldRepetition()).To(BeFalse())
			})
		})
	})
	Describe("::HasLegalNextMove", func() {
		var board *Board
		When("the board represents a stalemate", func() {
			BeforeEach(func() {
				board, _ = BoardFromFEN("k7/2Q5/8/8/8/4K3/8/8 b - - 0 1")
			})
			It("returns false", func() {
				Expect(board.HasLegalNextMove()).To(BeFalse())
			})
		})
		When("the board results in checkmate", func() {
			Context("and only one piece is giving a check", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("k7/8/8/8/8/4K3/R7/1R6 b - - 0 1")
				})
				It("returns false", func() {
					Expect(board.HasLegalNextMove()).To(BeFalse())
				})
			})
			Context("and two pieces are giving a check", func() {
				BeforeEach(func() {
					board, _ = BoardFromFEN("k7/2N5/8/8/8/4K3/R7/1R6 b - - 0 1")
				})
				It("returns false", func() {
					Expect(board.HasLegalNextMove()).To(BeFalse())
				})
			})
		})
		When("the board does not represent a stalemate or checkmate", func() {
			BeforeEach(func() {
				board, _ = BoardFromFEN("rnbqkbnr/pppp1ppp/8/4p3/3PP3/8/PPP2PPP/RNBQKBNR b KQkq d3 0 2")
			})
			It("returns true", func() {
				Expect(board.HasLegalNextMove()).To(BeTrue())
			})
		})
	})
	Describe("::getMaterialCount", func() {
		It("counts material of the initiate board", func() {
			board := GetInitBoard()
			materialCount := board.ComputeMaterialCount()
			Expect(materialCount.WhitePawnCount).To(Equal(uint8(8)))
			Expect(materialCount.WhiteKnightCount).To(Equal(uint8(2)))
			Expect(materialCount.WhiteLightBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.WhiteDarkBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.WhiteRookCount).To(Equal(uint8(2)))
			Expect(materialCount.WhiteQueenCount).To(Equal(uint8(1)))
			Expect(materialCount.BlackPawnCount).To(Equal(uint8(8)))
			Expect(materialCount.BlackKnightCount).To(Equal(uint8(2)))
			Expect(materialCount.BlackLightBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.BlackDarkBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.BlackRookCount).To(Equal(uint8(2)))
			Expect(materialCount.BlackQueenCount).To(Equal(uint8(1)))
		})
	})
	DescribeTable("::ToFEN", func(fen string) {
		board, _ := BoardFromFEN(fen)
		generatedFEN := board.ToFEN()
		Expect(generatedFEN).To(Equal(fen))
	},
		Entry("initial board", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
		Entry("no castle rights", "3R2R1/8/2R5/2Rk2R1/4R3/2R5/R2R4/7K w - - 0 1"),
		Entry("en passant square", "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1"),
		Entry("move counters boosted", "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 25 30"),
	)
	Describe("#BoardFromFEN", func() {
		When("the FEN is valid", func() {
			When("the FEN is the initial board", func() {
				It("returns exactly the init board", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
					expBoard := GetInitBoard()
					board, err := BoardFromFEN(fen)
					Expect(err).ToNot(HaveOccurred())
					Expect(board).ToNot(BeNil())
					Expect(board.FullMoveCount).To(Equal(expBoard.FullMoveCount))
					Expect(board.HalfMoveClockCount).To(Equal(expBoard.HalfMoveClockCount))
					Expect(board.CanBlackCastleKingside).To(Equal(expBoard.CanBlackCastleKingside))
					Expect(board.CanBlackCastleQueenside).To(Equal(expBoard.CanBlackCastleQueenside))
					Expect(board.CanWhiteCastleKingside).To(Equal(expBoard.CanWhiteCastleKingside))
					Expect(board.CanWhiteCastleQueenside).To(Equal(expBoard.CanWhiteCastleQueenside))
					Expect(board.IsWhiteTurn).To(Equal(expBoard.IsWhiteTurn))
					Expect(board.OptEnPassantSquare).To(Equal(expBoard.OptEnPassantSquare))
					for i := 0; i < 8; i++ {
						for j := 0; j < 8; j++ {
							piece := board.Pieces[i][j]
							expPiece := board.Pieces[i][j]
							Expect(piece).To(Equal(expPiece))
						}
					}
				})
			})
			When("the FEN specifies that neither player has castle rights", func() {
				It("returns a board with all castle rights revoked", func() {
					fen := "3R2R1/8/2R5/2Rk2R1/4R3/2R5/R2R4/7K w - - 0 1"
					board, err := BoardFromFEN(fen)
					Expect(err).ToNot(HaveOccurred())
					Expect(board.CanWhiteCastleQueenside).To(BeFalse())
					Expect(board.CanWhiteCastleKingside).To(BeFalse())
					Expect(board.CanBlackCastleQueenside).To(BeFalse())
					Expect(board.CanBlackCastleKingside).To(BeFalse())
				})
			})
			When("two white kings exist in the FEN pieces", func() {
				It("parses the board with no errors", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKKBNR w KQkq - 0 1"
					_, err := BoardFromFEN(fen)
					Expect(err).ToNot(HaveOccurred())
				})
			})
			When("the FEN specifies an en passant square", func() {
				It("returns a board with the en passant square set", func() {
					board, err := BoardFromFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
					Expect(err).ToNot(HaveOccurred())
					Expect(board.OptEnPassantSquare).To(Equal(&Square{Rank: 3, File: 5}))
				})
			})
			When("the FEN represents an insufficient material draw board", func() {
				It("returns a terminal draw board", func() {
					board, err := BoardFromFEN("8/8/2k2K2/5N2/8/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					Expect(board.IsTerminal).To(BeTrue())
					Expect(board.IsWhiteWinner).To(BeFalse())
					Expect(board.IsBlackWinner).To(BeFalse())
				})
			})
		})
		When("the FEN is not valid", func() {
			Context("the FEN does not have the correct amount of segments", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0"
					_, err := BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			When("the issue is with the pieces", func() {
				Context("the FEN has too many pieces rows", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN has one too few rows in pieces", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN has too many pieces on the first row", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnrp/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN has too few pieces on the first row", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbn/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN contains invalid piece chars", func() {
					It("returns an error", func() {
						invalidFEN := "xxxxxxxx/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
			})
			Context("the FEN does not have a valid turn specifier character", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR X KQkq - 0 1"
					_, err := BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			Context("the FEN contains an invalid castle rights specifier", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w XQkq - 0 1"
					_, err := BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			Context("the FEN contains an invalid enPassant square", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq jk4 0 1"
					_, err := BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			When("the issue is with the HalfMoveClockCount", func() {
				Context("the FEN contains a halfMoveClockCount greater than the range for a uint8", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 277 1"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN contains a non-integer as the halfMoveClockCount", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - X 1"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
			})
			When("the issue is with the FullMoveCount", func() {
				Context("the FEN contains a fullMoveCount greater than the range for a uint16", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 70700"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN contains a non-integer as the fullMoveCount", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 X"
						_, err := BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
			})
		})
	})
	Describe("::IsInitBoard", func() {
		When("the board is the initial board", func() {
			It("returns true", func() {
				board := GetInitBoard()
				Expect(board.IsInitBoard()).To(BeTrue())
			})
		})
		When("the board is not the initial board", func() {
			It("returns false", func() {
				board := GetInitBoard()
				board = GetBoardFromMove(board, &Move{WHITE_PAWN, &Square{2, 1}, &Square{4, 1}, EMPTY, nil, EMPTY})
				Expect(board.IsInitBoard()).To(BeFalse())
			})
		})
	})
})
