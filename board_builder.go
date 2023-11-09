package chess

type BoardBuilder struct {
	board *Board
}

func NewBoardBuilder() *BoardBuilder {
	board := Board{}
	board.RepetitionsByMiniFEN = make(map[string]uint8)
	return &BoardBuilder{
		board: &board,
	}
}

func (bb *BoardBuilder) WithPieces(pieces *[8][8]Piece) *BoardBuilder {
	bb.board.Pieces = *pieces
	bb.board.optMaterialCount = nil
	return bb
}

func (bb *BoardBuilder) WithPiece(piece Piece, square *Square) *BoardBuilder {
	prevPiece := bb.board.Pieces[square.Rank-1][square.File-1]
	bb.board.Pieces[square.Rank-1][square.File-1] = piece
	if bb.board.optMaterialCount != nil {
		materialCountBuilder := NewMaterialCountBuilder().WithMaterialCount(bb.board.optMaterialCount)
		if prevPiece != EMPTY {
			materialCountBuilder.WithoutPiece(prevPiece, square)
		}
		materialCountBuilder.WithPiece(piece, square)
		bb.board.optMaterialCount = materialCountBuilder.Build()
	}
	if prevPiece == WHITE_KING {
		bb.board.optWhiteKingSquare = nil
	}
	if prevPiece == BLACK_KING {
		bb.board.optBlackKingSquare = nil
	}
	if piece == WHITE_KING {
		bb.board.optWhiteKingSquare = square
	}
	if piece == BLACK_KING {
		bb.board.optBlackKingSquare = square
	}
	return bb
}

func (bb *BoardBuilder) WithEnPassantSquare(enPassantSquare *Square) *BoardBuilder {
	bb.board.OptEnPassantSquare = enPassantSquare
	return bb
}

func (bb *BoardBuilder) WithIsWhiteTurn(isWhiteTurn bool) *BoardBuilder {
	bb.board.IsWhiteTurn = isWhiteTurn
	return bb
}

func (bb *BoardBuilder) WithCanWhiteCastleQueenside(canWhiteCastleQueenside bool) *BoardBuilder {
	bb.board.CanWhiteCastleQueenside = canWhiteCastleQueenside
	return bb
}

func (bb *BoardBuilder) WithCanWhiteCastleKingside(canWhiteCastleKingside bool) *BoardBuilder {
	bb.board.CanWhiteCastleKingside = canWhiteCastleKingside
	return bb
}
func (bb *BoardBuilder) WithCanBlackCastleQueenside(canBlackCastleQueenside bool) *BoardBuilder {
	bb.board.CanBlackCastleQueenside = canBlackCastleQueenside
	return bb
}

func (bb *BoardBuilder) WithCanBlackCastleKingside(canBlackCastleKingside bool) *BoardBuilder {
	bb.board.CanBlackCastleKingside = canBlackCastleKingside
	return bb
}

func (bb *BoardBuilder) WithHalfMoveClockCount(halfMoveClockCount uint8) *BoardBuilder {
	bb.board.HalfMoveClockCount = halfMoveClockCount
	return bb
}

func (bb *BoardBuilder) WithFullMoveCount(fullMoveCount uint16) *BoardBuilder {
	bb.board.FullMoveCount = fullMoveCount
	return bb
}

func (bb *BoardBuilder) WithRepetitionsByMiniFEN(repetitionsByMiniFEN map[string]uint8) *BoardBuilder {
	bb.board.RepetitionsByMiniFEN = repetitionsByMiniFEN
	return bb
}

func (bb *BoardBuilder) WithMiniFENCount(miniFEN string, count uint8) *BoardBuilder {
	bb.board.RepetitionsByMiniFEN[miniFEN] = count
	return bb
}

func (bb *BoardBuilder) WithIsTerminal(isTerminal bool) *BoardBuilder {
	bb.board.IsTerminal = isTerminal
	return bb
}

func (bb *BoardBuilder) WithIsWhiteWinner(isWhiteWinner bool) *BoardBuilder {
	bb.board.IsWhiteWinner = isWhiteWinner
	return bb
}
func (bb *BoardBuilder) WithIsBlackWinner(isBlackWinner bool) *BoardBuilder {
	bb.board.IsBlackWinner = isBlackWinner
	return bb
}

func (bb *BoardBuilder) FromBoard(board *Board) *BoardBuilder {
	boardCopy := *board
	bb.board = &boardCopy
	return bb
}

func (bb *BoardBuilder) Build() *Board {
	return bb.board
}
