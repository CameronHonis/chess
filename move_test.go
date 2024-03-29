package chess_test

import (
	. "github.com/CameronHonis/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Move", func() {
	var move *Move
	Describe("::DoesAllowEnPassant", func() {
		When("the move does allow en passant", func() {
			BeforeEach(func() {
				move = &Move{WHITE_PAWN, &Square{2, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
			})
			It("returns true", func() {
				Expect(move.DoesAllowEnPassant()).To(BeTrue())
			})
		})
		When("the move does not allow en passant", func() {
			BeforeEach(func() {
				move = &Move{WHITE_PAWN, &Square{2, 5}, &Square{3, 5}, EMPTY, make([]*Square, 0), EMPTY}
			})
			It("returns false", func() {
				Expect(move.DoesAllowEnPassant()).To(BeFalse())
			})
		})
	})

	Describe("::Equal", func() {
		When("the moves are equal", func() {
			It("returns true", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.Equal(b)).To(BeTrue())
			})
			When("the moves list king checking squares in a different order", func() {
				It("returns true", func() {
					a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{1, 1}, {2, 2}, {3, 3}}, EMPTY}
					b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{3, 3}, {2, 2}, {1, 1}}, EMPTY}
					Expect(a.Equal(b)).To(BeTrue())
				})
			})
		})
		When("the moves only differ by piece", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_KNIGHT, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.Equal(b)).To(BeFalse())
			})
		})
		When("the moves only differ by start square", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 6}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.Equal(b)).To(BeFalse())
			})
		})
		When("the moves only differ by end square", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 6}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.Equal(b)).To(BeFalse())
			})
		})
		When("the moves only differ by captured piece", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, WHITE_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.Equal(b)).To(BeFalse())
			})
		})
		When("the moves only differ by checking squares", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{1, 1}, {2, 2}, {3, 3}}, EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{2, 2}, {3, 3}, {4, 4}}, EMPTY}
				Expect(a.Equal(b)).To(BeFalse())
			})
		})
	})

	Describe("::ToAlgebraic", func() {
		When("the move is a pawn move", func() {
			When("the move is a white pawn single jump", func() {
				It("includes only the target square", func() {
					move := &Move{WHITE_PAWN, &Square{2, 5}, &Square{3, 5}, EMPTY, make([]*Square, 0), EMPTY}
					board := GetInitBoard()
					Expect(move.ToAlgebraic(board)).To(Equal("e3"))
				})
			})
			When("the move is a black pawn single jump", func() {
				It("treats it the same as a white pawn jump", func() {
					move := &Move{BLACK_PAWN, &Square{7, 5}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR w kq - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("e6"))
				})
			})
			When("the move is a pawn double jump", func() {
				It("only includes the target square", func() {
					move := &Move{WHITE_PAWN, &Square{2, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
					board := GetInitBoard()
					Expect(move.ToAlgebraic(board)).To(Equal("e4"))
				})
			})
			When("the move is a pawn capture", func() {
				It("includes the file of the pawn capturing", func() {
					move := &Move{WHITE_PAWN, &Square{3, 4}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("rnbqkb1r/pppppppp/8/8/4n3/3P4/PPP1PPPP/RNBQKBNR w kq - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("dxe4"))
				})
				When("two pawns could have taken on the same square", func() {
					It("includes the file of the pawn", func() {
						move := &Move{WHITE_PAWN, &Square{3, 4}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
						board, _ := BoardFromFEN("rnbqkb1r/pppppppp/8/8/4n3/3P1P2/PPP1P1PP/RNBQKBNR w kq - 0 1")
						Expect(move.ToAlgebraic(board)).To(Equal("dxe4"))
					})
				})
			})
		})
		When("the move is a knight move", func() {
			It("returns N followed by the target square", func() {
				move := &Move{WHITE_KNIGHT, &Square{1, 2}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY}
				board := GetInitBoard()
				Expect(move.ToAlgebraic(board)).To(Equal("Nc3"))
			})
			When("the move is a capture", func() {
				It("returns Nx followed by the target square", func() {
					move := &Move{WHITE_KNIGHT, &Square{1, 2}, &Square{3, 3}, BLACK_PAWN, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("rnbqkbnr/pp1ppppp/8/8/8/2p5/PPPPPPPP/RNBQKBNR w kq - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("Nxc3"))
				})
			})
			When("two knights on different files could have moved to the same square", func() {
				It("distinguishes the knights by file", func() {
					move := &Move{BLACK_KNIGHT, &Square{8, 7}, &Square{7, 5}, EMPTY, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/2N5/PPPP1PPP/R1BQK1NR b kq - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("Nge7"))
				})
			})
			When("two knights on the same file could have moved to the same square", func() {
				It("distinguishes the knights by rank", func() {
					move := &Move{WHITE_KNIGHT, &Square{1, 7}, &Square{2, 5}, EMPTY, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("r1bqkb1r/pppp1ppp/2n5/4p3/2B1P3/6N1/PPPP1PPP/R1BQK1NR w kq - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("N1e2"))
				})
			})
			When("three knights cannot be distinguished by either rank or file", func() {
				It("distinguishes the knights by both rank and file", func() {
					move := &Move{WHITE_KNIGHT, &Square{7, 3}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("8/2N3N1/8/2N5/8/8/8/5K1k w - - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("Nc7e6"))
				})
			})
		})
		When("the move is a bishop move", func() {
			It("returns B followed by the target square", func() {
				move := &Move{WHITE_BISHOP, &Square{1, 4}, &Square{3, 2}, EMPTY, make([]*Square, 0), EMPTY}
				board, _ := BoardFromFEN("8/1p6/8/8/8/8/8/3B1K1k w - - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("Bb3"))
			})
			When("the move is a capture", func() {
				It("returns Bx followed by the target square", func() {
					move := &Move{WHITE_BISHOP, &Square{6, 1}, &Square{7, 2}, BLACK_PAWN, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("8/1p6/B7/8/8/8/7k/5K2 w - - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("Bxb7"))
				})
			})
			When("two bishops on different files could have moved to the same square", func() {
				It("distinguishes the bishops by file", func() {
					move := &Move{WHITE_BISHOP, &Square{2, 3}, &Square{3, 4}, EMPTY, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("8/1p6/8/8/8/8/2B1B2k/5K2 w - - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("Bcd3"))
				})
			})
			When("two bishops on the same file could have moved to the same square", func() {
				It("distinguishes the bishops by rank", func() {
					move := &Move{WHITE_BISHOP, &Square{4, 3}, &Square{3, 4}, EMPTY, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("8/1p6/8/8/2B5/8/2B4k/5K2 w - - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("B4d3"))
				})
			})
		})
		When("the move is a rook move", func() {
			It("returns R followed by the target square", func() {
				move := &Move{WHITE_ROOK, &Square{1, 5}, &Square{8, 5}, EMPTY, make([]*Square, 0), EMPTY}
				board, _ := BoardFromFEN("8/1p6/8/8/8/8/7k/4RK2 w - - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("Re8"))
			})
			When("the move is a capture", func() {
				It("returns Rx followed by the target square", func() {
					move := &Move{WHITE_ROOK, &Square{1, 2}, &Square{7, 2}, BLACK_PAWN, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("8/1p6/8/8/8/8/7k/1R3K2 w - - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("Rxb7"))
				})
			})
		})
		When("the move is a queen move", func() {
			It("returns Q followed by the target square", func() {
				move := &Move{WHITE_QUEEN, &Square{4, 1}, &Square{7, 1}, EMPTY, make([]*Square, 0), EMPTY}
				board, _ := BoardFromFEN("6rk/5bpp/5pn1/8/Q7/8/8/5K2 w - - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("Qa7"))
			})
			When("the move is a capture", func() {
				It("returns Qx followed by the target square", func() {
					move := &Move{WHITE_QUEEN, &Square{7, 2}, &Square{7, 6}, BLACK_BISHOP, make([]*Square, 0), EMPTY}
					board, _ := BoardFromFEN("6rk/1Q3bpp/5pn1/8/8/8/8/5K2 w - - 0 1")
					Expect(move.ToAlgebraic(board)).To(Equal("Qxf7"))
				})
				When("two queens on the same rank could have moved to the same square", func() {
					It("distinguishes the queens by file", func() {
						move := &Move{WHITE_QUEEN, &Square{4, 6}, &Square{6, 4}, EMPTY, make([]*Square, 0), EMPTY}
						board, _ := BoardFromFEN("6rk/5bpp/5pn1/8/1Q3Q2/8/8/5K2 w - - 0 1")
						Expect(move.ToAlgebraic(board)).To(Equal("Qfd6"))
					})
				})
				When("two queens on the same file could have moved to the same square", func() {
					It("distinguishes the queens by rank", func() {
						move := &Move{WHITE_QUEEN, &Square{8, 2}, &Square{6, 4}, EMPTY, make([]*Square, 0), EMPTY}
						board, _ := BoardFromFEN("1Q4rk/5bpp/5pn1/8/1Q6/8/8/5K2 w - - 0 1")
						Expect(move.ToAlgebraic(board)).To(Equal("Q8d6"))
					})
				})
				When("three queens cannot be distinguished by either rank or file", func() {
					It("distinguishes the queens by both rank and file, with Qx prefix", func() {
						move := &Move{WHITE_QUEEN, &Square{4, 2}, &Square{2, 4}, EMPTY, make([]*Square, 0), EMPTY}
						board, _ := BoardFromFEN("6rk/5bpp/5pn1/8/1Q1Q4/8/1Q6/5K2 w - - 0 1")
						Expect(move.ToAlgebraic(board)).To(Equal("Qb4d2"))
					})
				})
			})
		})
		When("the move is a king move", func() {
			It("returns K followed by the target square", func() {
				move := &Move{WHITE_KING, &Square{1, 6}, &Square{2, 7}, EMPTY, make([]*Square, 0), EMPTY}
				board, _ := BoardFromFEN("6rk/5bpp/5pn1/8/8/8/8/5K2 w - - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("Kg2"))
			})
		})
		When("the move is kingside castling", func() {
			It("returns O-O", func() {
				move := &Move{WHITE_KING, &Square{1, 5}, &Square{1, 7}, EMPTY, make([]*Square, 0), EMPTY}
				board, _ := BoardFromFEN("6rk/5bpp/5pn1/8/8/8/8/4K2R w K - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("O-O"))
			})
		})
		When("the move is queenside castling", func() {
			It("returns O-O-O", func() {
				move := &Move{BLACK_KING, &Square{8, 5}, &Square{8, 3}, EMPTY, make([]*Square, 0), EMPTY}
				board, _ := BoardFromFEN("r3k3/5bpp/5pn1/8/8/8/8/4K2R b q - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("O-O-O"))
			})
		})
		When("the move is a pawn promotion", func() {
			It("returns the format {target square}8={promotion piece}", func() {
				move := &Move{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_ROOK}
				board, _ := BoardFromFEN("6rk/1P3bpp/5pn1/8/8/8/8/4K2R w - - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("b8=R"))
			})
		})
		When("the move results in a check", func() {
			It("it appends '+' to the end", func() {
				move := &Move{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, []*Square{{8, 2}}, WHITE_QUEEN}
				board, _ := BoardFromFEN("7k/1P3bpp/5pn1/8/8/8/8/4K2R w - - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("b8=Q+"))
			})
		})
		When("the move results in a checkmate", func() {
			It("appends '#'", func() {
				move := &Move{WHITE_ROOK, &Square{1, 4}, &Square{8, 4}, EMPTY, []*Square{{8, 4}}, EMPTY}
				board, _ := BoardFromFEN("7k/5ppp/8/8/8/8/8/3RK3 w - - 0 1")
				Expect(move.ToAlgebraic(board)).To(Equal("Rd8#"))
			})
		})
	})

	Describe("::ToLongAlgebraic", func() {
		When("the move is a pawn move", func() {
			BeforeEach(func() {
				move = &Move{WHITE_PAWN, &Square{1, 1}, &Square{4, 8}, EMPTY, make([]*Square, 0), EMPTY}
			})
			It("includes only the start & end square", func() {
				Expect(move.ToLongAlgebraic()).To(Equal("a1h4"))
			})
		})
		When("the move is a capture", func() {
			When("not en passant", func() {
				BeforeEach(func() {
					move = &Move{WHITE_PAWN, &Square{1, 1}, &Square{4, 8}, BLACK_ROOK, make([]*Square, 0), EMPTY}
				})
				It("includes only the start & end square", func() {
					Expect(move.ToLongAlgebraic()).To(Equal("a1h4"))
				})
			})
			When("is en passant", func() {
				BeforeEach(func() {
					move = &Move{WHITE_PAWN, &Square{6, 5}, &Square{7, 4}, BLACK_PAWN, make([]*Square, 0), EMPTY}
				})
				It("includes only the start & end square", func() {
					Expect(move.ToLongAlgebraic()).To(Equal("e6d7"))
				})
			})
		})
		When("the move is a pawn upgrade", func() {
			When("upgrade to queen", func() {
				BeforeEach(func() {
					move = &Move{WHITE_PAWN, &Square{7, 7}, &Square{8, 7}, EMPTY, make([]*Square, 0), WHITE_QUEEN}
				})
				It("includes the start & end square", func() {
					Expect(strings.HasPrefix(move.ToLongAlgebraic(), "g7g8")).To(BeTrue())
				})
				It("includes 'q' at the end", func() {
					Expect(strings.HasSuffix(move.ToLongAlgebraic(), "q")).To(BeTrue())
				})
			})
			When("upgrade to rook", func() {
				BeforeEach(func() {
					move = &Move{WHITE_PAWN, &Square{7, 7}, &Square{8, 7}, EMPTY, make([]*Square, 0), WHITE_ROOK}
				})
				It("includes the start & end square", func() {
					Expect(strings.HasPrefix(move.ToLongAlgebraic(), "g7g8")).To(BeTrue())
				})
				It("includes 'q' at the end", func() {
					Expect(strings.HasSuffix(move.ToLongAlgebraic(), "r")).To(BeTrue())
				})
			})
			When("upgrade to bishop", func() {
				BeforeEach(func() {
					move = &Move{WHITE_PAWN, &Square{7, 7}, &Square{8, 7}, EMPTY, make([]*Square, 0), WHITE_BISHOP}
				})
				It("includes the start & end square", func() {
					Expect(strings.HasPrefix(move.ToLongAlgebraic(), "g7g8")).To(BeTrue())
				})
				It("includes 'q' at the end", func() {
					Expect(strings.HasSuffix(move.ToLongAlgebraic(), "b")).To(BeTrue())
				})
			})
			When("upgrade to knight", func() {
				BeforeEach(func() {
					move = &Move{WHITE_PAWN, &Square{7, 7}, &Square{8, 7}, EMPTY, make([]*Square, 0), WHITE_KNIGHT}
				})
				It("includes the start & end square", func() {
					Expect(strings.HasPrefix(move.ToLongAlgebraic(), "g7g8")).To(BeTrue())
				})
				It("includes 'q' at the end", func() {
					Expect(strings.HasSuffix(move.ToLongAlgebraic(), "n")).To(BeTrue())
				})
			})
		})
	})

	Describe("#MoveFromAlgebraic", func() {
		var board *Board
		var algMove string
		When("the move is a pawn move", func() {
			When("the move is a white pawn single jump", func() {
				BeforeEach(func() {
					algMove = "e3"
					board = GetInitBoard()
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_PAWN, &Square{2, 5}, &Square{3, 5}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("the move is a black pawn single jump", func() {
				BeforeEach(func() {
					algMove = "e6"
					board, _ = BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b kq - 0 1")
				})
				It("generates the correct move", func() {
					expMove := &Move{BLACK_PAWN, &Square{7, 5}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("the move is a pawn double jump", func() {
				BeforeEach(func() {
					algMove = "e4"
					board = GetInitBoard()
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_PAWN, &Square{2, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("the move is a pawn capture", func() {
				When("only one pawn can take", func() {
					BeforeEach(func() {
						algMove = "dxe4"
						board, _ = BoardFromFEN("rnbqkb1r/pppppppp/8/8/4n3/3P4/PPP1PPPP/RNBQKBNR w kq - 0 1")
					})
					It("generates the correct move", func() {
						expMove := &Move{WHITE_PAWN, &Square{3, 4}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
						Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
					})
				})
				When("two pawns could have taken on the same square", func() {
					BeforeEach(func() {
						algMove = "dxe4"
						board, _ = BoardFromFEN("rnbqkb1r/pppppppp/8/8/4n3/3P1P2/PPP1P1PP/RNBQKBNR w kq - 0 1")
					})
					It("generates the correct move", func() {
						expMove := &Move{WHITE_PAWN, &Square{3, 4}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
						Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
					})
				})
			})
		})
		When("the move is a knight move", func() {
			When("the move is not a capture", func() {
				BeforeEach(func() {
					algMove = "Nc3"
					board = GetInitBoard()
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_KNIGHT, &Square{1, 2}, &Square{3, 3}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("the move is a capture", func() {
				BeforeEach(func() {
					algMove = "Nxc3"
					board, _ = BoardFromFEN("rnbqkbnr/pp1ppppp/8/8/8/2p5/PPPPPPPP/RNBQKBNR w kq - 0 1")
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_KNIGHT, &Square{1, 2}, &Square{3, 3}, BLACK_PAWN, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("two knights on different files could have moved to the same square", func() {
				BeforeEach(func() {
					algMove = "Nge7"
					board, _ = BoardFromFEN("r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/2N5/PPPP1PPP/R1BQK1NR b kq - 0 1")
				})
				It("generates the correct move", func() {
					expMove := &Move{BLACK_KNIGHT, &Square{8, 7}, &Square{7, 5}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("two knights on the same file could have moved to the same square", func() {
				BeforeEach(func() {
					algMove = "N1e2"
					board, _ = BoardFromFEN("r1bqkb1r/pppp1ppp/2n5/4p3/2B1P3/6N1/PPPP1PPP/R1BQK1NR w kq - 0 1")
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_KNIGHT, &Square{1, 7}, &Square{2, 5}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("three knights cannot be distinguished by either rank or file", func() {
				BeforeEach(func() {
					algMove = "Nc7e6"
					board, _ = BoardFromFEN("8/2N3N1/8/2N5/8/8/8/5K1k w - - 0 1")
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_KNIGHT, &Square{7, 3}, &Square{6, 5}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
		})
		When("the move is a bishop move", func() {
			BeforeEach(func() {
				algMove = "Bb3"
				board, _ = BoardFromFEN("8/1p6/8/8/8/8/8/3B1K1k w - - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{WHITE_BISHOP, &Square{1, 4}, &Square{3, 2}, EMPTY, make([]*Square, 0), EMPTY}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
		When("the move is a rook move", func() {
			BeforeEach(func() {
				algMove = "Re8"
				board, _ = BoardFromFEN("8/1p6/8/8/8/8/7k/4RK2 w - - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{WHITE_ROOK, &Square{1, 5}, &Square{8, 5}, EMPTY, make([]*Square, 0), EMPTY}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
		When("the move is a queen move", func() {
			When("a single queen can move to the destination square", func() {
				BeforeEach(func() {
					algMove = "Qa7"
					board, _ = BoardFromFEN("6rk/5bpp/5pn1/8/Q7/8/8/5K2 w - - 0 1")
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_QUEEN, &Square{4, 1}, &Square{7, 1}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
			When("queens on a different rank and file can move to the destination square", func() {
				BeforeEach(func() {
					algMove = "Qd3g6"
					board, _ = BoardFromFEN("krr5/ppp5/nprQ4/8/8/3Q2Q1/8/4K3 w - - 0 1")
				})
				It("generates the correct move", func() {
					expMove := &Move{WHITE_QUEEN, &Square{3, 4}, &Square{6, 7}, EMPTY, make([]*Square, 0), EMPTY}
					Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
				})
			})
		})
		When("the move is a king move", func() {
			BeforeEach(func() {
				algMove = "Kg2"
				board, _ = BoardFromFEN("6rk/5bpp/5pn1/8/8/8/8/5K2 w - - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{WHITE_KING, &Square{1, 6}, &Square{2, 7}, EMPTY, make([]*Square, 0), EMPTY}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
		When("the move is kingside castling", func() {
			BeforeEach(func() {
				algMove = "O-O"
				board, _ = BoardFromFEN("6rk/5bpp/5pn1/8/8/8/8/4K2R w K - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{WHITE_KING, &Square{1, 5}, &Square{1, 7}, EMPTY, make([]*Square, 0), EMPTY}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
		When("the move is queenside castling", func() {
			BeforeEach(func() {
				algMove = "O-O-O"
				board, _ = BoardFromFEN("r3k3/5bpp/5pn1/8/8/8/8/4K2R b q - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{BLACK_KING, &Square{8, 5}, &Square{8, 3}, EMPTY, make([]*Square, 0), EMPTY}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
		When("the move is a pawn promotion", func() {
			BeforeEach(func() {
				algMove = "b8=R"
				board, _ = BoardFromFEN("6rk/1P3bpp/5pn1/8/8/8/8/4K2R w - - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_ROOK}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
		When("the move results in a check from a promotion", func() {
			BeforeEach(func() {
				algMove = "b8=Q+"
				board, _ = BoardFromFEN("7k/1P3bpp/5pn1/8/8/8/8/4K2R w - - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, []*Square{{8, 2}}, WHITE_QUEEN}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
		When("the move results in a checkmate", func() {
			BeforeEach(func() {
				algMove = "Rd8#"
				board, _ = BoardFromFEN("7k/5ppp/8/8/8/8/8/3RK3 w - - 0 1")
			})
			It("generates the correct move", func() {
				expMove := &Move{WHITE_ROOK, &Square{1, 4}, &Square{8, 4}, EMPTY, []*Square{{8, 4}}, EMPTY}
				Expect(MoveFromAlgebraic(algMove, board)).To(BeComparableTo(expMove))
			})
		})
	})

	Describe("#MoveFromLongAlgebraic", func() {
		It("handles a pawn upgrade checking move", func() {
			longAlgMove := "b7b8q"
			board, _ := BoardFromFEN("7k/1P3bpp/5pn1/8/8/8/8/4K2R w - - 0 1")
			expMove := &Move{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, []*Square{{8, 2}}, WHITE_QUEEN}
			Expect(MoveFromLongAlgebraic(longAlgMove, board)).To(BeComparableTo(expMove))
		})
	})
})

//var _ = FIt("go-cmp test", func() {
//	m1 := &Move{WHITE_PAWN, &Square{1, 2}, &Square{3, 4}, EMPTY, []*Square{{4, 4}, {2, 2}}, WHITE_QUEEN}
//	m2 := &Move{WHITE_PAWN, &Square{1, 2}, &Square{3, 4}, EMPTY, []*Square{{4, 4}, {2, 2}}, WHITE_QUEEN}
//	m3 := &Move{WHITE_PAWN, &Square{1, 2}, &Square{3, 4}, EMPTY, []*Square{{4, 4}, {3, 3}}, WHITE_QUEEN}
//	Expect(m1).To(BeComparableTo(m2))
//	Expect(m1).ToNot(BeComparableTo(m3))
//})
