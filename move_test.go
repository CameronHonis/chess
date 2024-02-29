package chess_test

import (
	. "github.com/CameronHonis/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Move", func() {
	Describe("::DoesAllowEnPassant", func() {
		When("the move does allow en passant", func() {
			move := Move{WHITE_PAWN, &Square{2, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
			It("returns true", func() {
				Expect(move.DoesAllowEnPassant()).To(BeTrue())
			})
		})
		When("the move does not allow en passant", func() {
			move := Move{WHITE_PAWN, &Square{2, 5}, &Square{3, 5}, EMPTY, make([]*Square, 0), EMPTY}
			It("returns false", func() {
				Expect(move.DoesAllowEnPassant()).To(BeFalse())
			})
		})
	})

	Describe("::EqualTo", func() {
		When("the moves are equal", func() {
			It("returns true", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.EqualTo(b)).To(BeTrue())
			})
			When("the moves list king checking squares in a different order", func() {
				It("returns true", func() {
					a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{1, 1}, {2, 2}, {3, 3}}, EMPTY}
					b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{3, 3}, {2, 2}, {1, 1}}, EMPTY}
					Expect(a.EqualTo(b)).To(BeTrue())
				})
			})
		})
		When("the moves only differ by piece", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_KNIGHT, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.EqualTo(b)).To(BeFalse())
			})
		})
		When("the moves only differ by start square", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 6}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.EqualTo(b)).To(BeFalse())
			})
		})
		When("the moves only differ by end square", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 6}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.EqualTo(b)).To(BeFalse())
			})
		})
		When("the moves only differ by captured piece", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, make([]*Square, 0), EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, WHITE_KNIGHT, make([]*Square, 0), EMPTY}
				Expect(a.EqualTo(b)).To(BeFalse())
			})
		})
		When("the moves only differ by checking squares", func() {
			It("returns false", func() {
				a := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{1, 1}, {2, 2}, {3, 3}}, EMPTY}
				b := &Move{WHITE_QUEEN, &Square{2, 5}, &Square{4, 5}, BLACK_KNIGHT, []*Square{{2, 2}, {3, 3}, {4, 4}}, EMPTY}
				Expect(a.EqualTo(b)).To(BeFalse())
			})
		})
	})
})
