package chess

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
