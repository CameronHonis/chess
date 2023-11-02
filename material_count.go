package chess

type MaterialCount struct {
	WhitePawnCount        uint8
	WhiteKnightCount      uint8
	WhiteLightBishopCount uint8
	WhiteDarkBishopCount  uint8
	WhiteRookCount        uint8
	WhiteQueenCount       uint8
	BlackPawnCount        uint8
	BlackKnightCount      uint8
	BlackLightBishopCount uint8
	BlackDarkBishopCount  uint8
	BlackRookCount        uint8
	BlackQueenCount       uint8
}

type MaterialCountBuilder struct {
	materialCount *MaterialCount
}

func NewMaterialCountBuilder() *MaterialCountBuilder {
	return &MaterialCountBuilder{
		materialCount: &MaterialCount{},
	}
}

func (mcb *MaterialCountBuilder) WithPiece(piece Piece, square *Square) *MaterialCountBuilder {
	return mcb.incrPiece(piece, square, 1)
}

func (mcb *MaterialCountBuilder) WithoutPiece(piece Piece, square *Square) *MaterialCountBuilder {
	return mcb.incrPiece(piece, square, -1)
}

func (mcb *MaterialCountBuilder) incrPiece(piece Piece, square *Square, incr int8) *MaterialCountBuilder {
	switch piece {
	case WHITE_PAWN:
		mcb.materialCount.WhitePawnCount += uint8(incr)
	case WHITE_KNIGHT:
		mcb.materialCount.WhiteKnightCount += uint8(incr)
	case WHITE_BISHOP:
		if square.IsLightSquare() {
			mcb.materialCount.WhiteLightBishopCount += uint8(incr)
		} else {
			mcb.materialCount.WhiteDarkBishopCount += uint8(incr)
		}
	case WHITE_ROOK:
		mcb.materialCount.WhiteRookCount += uint8(incr)
	case WHITE_QUEEN:
		mcb.materialCount.WhiteQueenCount += uint8(incr)
	case BLACK_PAWN:
		mcb.materialCount.BlackPawnCount += uint8(incr)
	case BLACK_KNIGHT:
		mcb.materialCount.BlackKnightCount += uint8(incr)
	case BLACK_BISHOP:
		if square.IsLightSquare() {
			mcb.materialCount.BlackLightBishopCount += uint8(incr)
		} else {
			mcb.materialCount.BlackDarkBishopCount += uint8(incr)
		}
	case BLACK_ROOK:
		mcb.materialCount.BlackRookCount += uint8(incr)
	case BLACK_QUEEN:
		mcb.materialCount.BlackQueenCount += uint8(incr)
	}
	return mcb
}

func (mcb *MaterialCountBuilder) WithMaterialCount(otherMaterialCount *MaterialCount) *MaterialCountBuilder {
	materialCountCopy := *otherMaterialCount
	mcb.materialCount = &materialCountCopy
	return mcb
}

func (mcb *MaterialCountBuilder) Build() *MaterialCount {
	return mcb.materialCount
}
