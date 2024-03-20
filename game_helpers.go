package chess

import (
	"fmt"
	"math"
)

func GetCheckingSquares(board *Board, isWhiteKing bool) []*Square {
	checkingSquares := make([]*Square, 0)
	kingSquare := board.GetKingSquare(isWhiteKing)
	knightCheckSquares := []*Square{
		{kingSquare.Rank + 2, kingSquare.File + 1},
		{kingSquare.Rank + 2, kingSquare.File - 1},
		{kingSquare.Rank + 1, kingSquare.File + 2},
		{kingSquare.Rank + 1, kingSquare.File - 2},
		{kingSquare.Rank - 1, kingSquare.File + 2},
		{kingSquare.Rank - 1, kingSquare.File - 2},
		{kingSquare.Rank - 2, kingSquare.File + 1},
		{kingSquare.Rank - 2, kingSquare.File - 1},
	}
	for _, knightCheckSquare := range knightCheckSquares {
		if !knightCheckSquare.IsValidBoardSquare() {
			continue
		}
		piece := board.GetPieceOnSquare(knightCheckSquare)
		if isWhiteKing && piece == BLACK_KNIGHT {
			checkingSquares = append(checkingSquares, knightCheckSquare)
		} else if !isWhiteKing && piece == WHITE_KNIGHT {
			checkingSquares = append(checkingSquares, knightCheckSquare)
		}
	}
	var pawnCheckSquares []*Square
	if isWhiteKing {
		pawnCheckSquares = []*Square{
			{kingSquare.Rank + 1, kingSquare.File - 1},
			{kingSquare.Rank + 1, kingSquare.File + 1},
		}
	} else {
		pawnCheckSquares = []*Square{
			{kingSquare.Rank - 1, kingSquare.File - 1},
			{kingSquare.Rank - 1, kingSquare.File + 1},
		}
	}
	for _, pawnCheckSquare := range pawnCheckSquares {
		if !pawnCheckSquare.IsValidBoardSquare() {
			continue
		}
		piece := board.GetPieceOnSquare(pawnCheckSquare)
		if isWhiteKing && piece == BLACK_PAWN {
			checkingSquares = append(checkingSquares, pawnCheckSquare)
		} else if !isWhiteKing && piece == WHITE_PAWN {
			checkingSquares = append(checkingSquares, pawnCheckSquare)
		}
	}
	for _, diagDir := range [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
		for dis := 1; dis < 8; dis++ {
			diagSquare := Square{
				kingSquare.Rank + uint8(dis*diagDir[0]),
				kingSquare.File + uint8(dis*diagDir[1])}
			if !diagSquare.IsValidBoardSquare() {
				break
			}
			piece := board.GetPieceOnSquare(&diagSquare)
			if isWhiteKing {
				if piece == BLACK_BISHOP || piece == BLACK_QUEEN {
					checkingSquares = append(checkingSquares, &diagSquare)
				}
			} else {
				if piece == WHITE_BISHOP || piece == WHITE_QUEEN {
					checkingSquares = append(checkingSquares, &diagSquare)
				}
			}
			if piece != EMPTY {
				break
			}
		}
	}
	for _, straightDir := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		for dis := 1; dis < 8; dis++ {
			straightSquare := Square{
				kingSquare.Rank + uint8(dis*straightDir[0]),
				kingSquare.File + uint8(dis*straightDir[1]),
			}
			if !straightSquare.IsValidBoardSquare() {
				break
			}
			piece := board.GetPieceOnSquare(&straightSquare)
			if isWhiteKing {
				if piece == BLACK_ROOK || piece == BLACK_QUEEN {
					checkingSquares = append(checkingSquares, &straightSquare)
				}
			} else {
				if piece == WHITE_ROOK || piece == WHITE_QUEEN {
					checkingSquares = append(checkingSquares, &straightSquare)
				}
			}
			if piece != EMPTY {
				break
			}
		}
	}
	return checkingSquares
}

func filterMovesByKingSafety(board *Board, moves []*Move) []*Move {
	filteredMoves := make([]*Move, 0, len(moves))
	for _, move := range moves {
		boardBuilder := NewBoardBuilder().FromBoard(board)
		UpdatePiecesFromMove(board, boardBuilder, move)
		newBoard := boardBuilder.Build()
		checkingSquares := GetCheckingSquares(newBoard, board.IsWhiteTurn)
		if len(checkingSquares) == 0 {
			filteredMoves = append(filteredMoves, move)
		}
	}
	return filteredMoves
}

func addKingChecksToMoves(board *Board, moves *[]*Move) {
	for _, move := range *moves {
		boardBuilder := NewBoardBuilder().FromBoard(board)
		UpdatePiecesFromMove(board, boardBuilder, move)
		newBoard := boardBuilder.Build()
		move.KingCheckingSquares = GetCheckingSquares(newBoard, !board.IsWhiteTurn)
	}
}
func GetLegalMovesForPawn(board *Board, square *Square) ([]*Move, error) {
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		emptyMoves := make([]*Move, 0)
		return emptyMoves, nil
	}
	pawnMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	var upgradePieces [4]Piece
	var squareInFront Square
	var leftCaptureSquare Square
	var rightCaptureSquare Square
	if board.IsWhiteTurn {
		if piece != WHITE_PAWN {
			return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_PAWN", piece)
		}
		upgradePieces = [4]Piece{WHITE_KNIGHT, WHITE_BISHOP, WHITE_ROOK, WHITE_QUEEN}
		squareInFront = Square{square.Rank + 1, square.File}
		leftCaptureSquare = Square{square.Rank + 1, square.File - 1}
		rightCaptureSquare = Square{square.Rank + 1, square.File + 1}
	} else {
		if piece != BLACK_PAWN {
			return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_PAWN", piece)
		}
		upgradePieces = [4]Piece{BLACK_KNIGHT, BLACK_BISHOP, BLACK_ROOK, BLACK_QUEEN}
		squareInFront = Square{square.Rank - 1, square.File}
		leftCaptureSquare = Square{square.Rank - 1, square.File - 1}
		rightCaptureSquare = Square{square.Rank - 1, square.File + 1}
	}
	pieceInFront := board.GetPieceOnSquare(&squareInFront)
	if pieceInFront == EMPTY {
		if squareInFront.Rank == 8 || squareInFront.Rank == 1 {
			move0 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[0]}
			move1 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[1]}
			move2 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[2]}
			move3 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[3]}
			pawnMoves = append(pawnMoves, &move0, &move1, &move2, &move3)
		} else {
			move := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), EMPTY}
			pawnMoves = append(pawnMoves, &move)
		}
		if (board.IsWhiteTurn && square.Rank == 2) || (!board.IsWhiteTurn && square.Rank == 7) {
			var squareTwoInFront Square
			if board.IsWhiteTurn {
				squareTwoInFront = Square{square.Rank + 2, square.File}
			} else {
				squareTwoInFront = Square{square.Rank - 2, square.File}
			}
			pieceTwoInFront := board.GetPieceOnSquare(&squareTwoInFront)
			if pieceTwoInFront == EMPTY {
				move := Move{piece, square, &squareTwoInFront, EMPTY, make([]*Square, 0), EMPTY}
				pawnMoves = append(pawnMoves, &move)
			}
		}
	}
	if leftCaptureSquare.IsValidBoardSquare() {
		var leftCapturePiece Piece
		if board.OptEnPassantSquare != nil && leftCaptureSquare.EqualTo(board.OptEnPassantSquare) {
			leftCapturePiece = board.GetPieceOnSquare(&Square{square.Rank, square.File - 1})
		} else {
			leftCapturePiece = board.GetPieceOnSquare(&leftCaptureSquare)
		}
		if leftCapturePiece != EMPTY && leftCapturePiece.IsWhite() != piece.IsWhite() {
			if (piece.IsWhite() && square.Rank == 7) || (!piece.IsWhite() && square.Rank == 2) {
				move0 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[0]}
				move1 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[1]}
				move2 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[2]}
				move3 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[3]}
				pawnMoves = append(pawnMoves, &move0, &move1, &move2, &move3)
			} else {
				move := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), EMPTY}
				pawnMoves = append(pawnMoves, &move)
			}
		}
	}
	if rightCaptureSquare.IsValidBoardSquare() {
		var rightCapturePiece Piece
		if board.OptEnPassantSquare != nil && rightCaptureSquare.EqualTo(board.OptEnPassantSquare) {
			rightCapturePiece = board.GetPieceOnSquare(&Square{square.Rank, square.File + 1})
		} else {
			rightCapturePiece = board.GetPieceOnSquare(&rightCaptureSquare)
		}
		if rightCapturePiece != EMPTY && rightCapturePiece.IsWhite() != piece.IsWhite() {
			if (piece.IsWhite() && square.Rank == 7) || (!piece.IsWhite() && square.Rank == 2) {
				move0 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[0]}
				move1 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[1]}
				move2 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[2]}
				move3 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[3]}
				pawnMoves = append(pawnMoves, &move0, &move1, &move2, &move3)
			} else {
				move := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), EMPTY}
				pawnMoves = append(pawnMoves, &move)
			}
		}
	}
	pawnMoves = filterMovesByKingSafety(board, pawnMoves)
	addKingChecksToMoves(board, &pawnMoves)
	return pawnMoves, nil
}

func GetLegalMovesForKnight(board *Board, square *Square) ([]*Move, error) {
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		emptyMoves := make([]*Move, 0)
		return emptyMoves, nil
	}
	knightMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_KNIGHT {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_KNIGHT", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_KNIGHT {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_KNIGHT", piece)
	}
	landSquares := []*Square{
		{square.Rank + 2, square.File - 1},
		{square.Rank + 2, square.File + 1},
		{square.Rank + 1, square.File - 2},
		{square.Rank + 1, square.File + 2},
		{square.Rank - 1, square.File - 2},
		{square.Rank - 1, square.File + 2},
		{square.Rank - 2, square.File - 1},
		{square.Rank - 2, square.File + 1},
	}
	for _, landSquare := range landSquares {
		if !landSquare.IsValidBoardSquare() {
			continue
		}
		landPiece := board.GetPieceOnSquare(landSquare)
		if landPiece == EMPTY {
			move := Move{piece, square, landSquare, EMPTY, make([]*Square, 0), EMPTY}
			knightMoves = append(knightMoves, &move)
		} else if landPiece.IsWhite() != board.IsWhiteTurn {
			move := Move{piece, square, landSquare, landPiece, make([]*Square, 0), EMPTY}
			knightMoves = append(knightMoves, &move)
		}
	}
	knightMoves = filterMovesByKingSafety(board, knightMoves)
	addKingChecksToMoves(board, &knightMoves)
	return knightMoves, nil
}

func GetLegalMovesForBishop(board *Board, square *Square) ([]*Move, error) {
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		emptyMoves := make([]*Move, 0)
		return emptyMoves, nil
	}
	bishopMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_BISHOP {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_BISHOP", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_BISHOP {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_BISHOP", piece)
	}
	for _, dir := range [4][2]int{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}} {
		dis := 0
		for {
			dis++
			landSquare := Square{square.Rank + uint8(dis*dir[0]), square.File + uint8(dis*dir[1])}
			if !landSquare.IsValidBoardSquare() {
				break
			}
			landPiece := board.GetPieceOnSquare(&landSquare)
			if landPiece == EMPTY {
				move := Move{piece, square, &landSquare, EMPTY, make([]*Square, 0), EMPTY}
				bishopMoves = append(bishopMoves, &move)
			} else {
				if piece.IsWhite() != landPiece.IsWhite() {
					move := Move{piece, square, &landSquare, landPiece, make([]*Square, 0), EMPTY}
					bishopMoves = append(bishopMoves, &move)
				}
				break
			}
		}
	}
	bishopMoves = filterMovesByKingSafety(board, bishopMoves)
	addKingChecksToMoves(board, &bishopMoves)
	return bishopMoves, nil
}

func GetLegalMovesForRook(board *Board, square *Square) ([]*Move, error) {
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		emptyMoves := make([]*Move, 0)
		return emptyMoves, nil
	}
	rookMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_ROOK {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_ROOK", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_ROOK {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_ROOK", piece)
	}
	for _, dir := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		dis := 0
		for {
			dis++
			landSquare := Square{square.Rank + uint8(dis*dir[0]), square.File + uint8(dis*dir[1])}
			if !landSquare.IsValidBoardSquare() {
				break
			}
			landPiece := board.GetPieceOnSquare(&landSquare)
			if landPiece == EMPTY {
				move := Move{piece, square, &landSquare, EMPTY, make([]*Square, 0), EMPTY}
				rookMoves = append(rookMoves, &move)
			} else {
				if piece.IsWhite() != landPiece.IsWhite() {
					move := Move{piece, square, &landSquare, landPiece, make([]*Square, 0), EMPTY}
					rookMoves = append(rookMoves, &move)
				}
				break
			}
		}
	}
	rookMoves = filterMovesByKingSafety(board, rookMoves)
	addKingChecksToMoves(board, &rookMoves)
	return rookMoves, nil
}

func GetLegalMovesForQueen(board *Board, square *Square) ([]*Move, error) {
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		emptyMoves := make([]*Move, 0)
		return emptyMoves, nil
	}
	queenMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_QUEEN {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_ROOK", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_QUEEN {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_ROOK", piece)
	}
	for _, dir := range [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
		dis := 0
		for {
			dis++
			landSquare := Square{square.Rank + uint8(dis*dir[0]), square.File + uint8(dis*dir[1])}
			if !landSquare.IsValidBoardSquare() {
				break
			}
			landPiece := board.GetPieceOnSquare(&landSquare)
			if landPiece == EMPTY {
				move := Move{piece, square, &landSquare, EMPTY, make([]*Square, 0), EMPTY}
				queenMoves = append(queenMoves, &move)
			} else {
				if piece.IsWhite() != landPiece.IsWhite() {
					move := Move{piece, square, &landSquare, landPiece, make([]*Square, 0), EMPTY}
					queenMoves = append(queenMoves, &move)
				}
				break
			}
		}
	}
	queenMoves = filterMovesByKingSafety(board, queenMoves)
	addKingChecksToMoves(board, &queenMoves)
	return queenMoves, nil
}

func GetLegalMovesForKing(board *Board) []*Move {
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		emptyMoves := make([]*Move, 0)
		return emptyMoves
	}
	kingMoves := make([]*Move, 0)
	whiteKingSquare, blackKingSquare := board.ComputeKingPositions()
	var square *Square
	var piece Piece
	if board.IsWhiteTurn {
		square = whiteKingSquare
		piece = WHITE_KING
	} else {
		square = blackKingSquare
		piece = BLACK_KING
	}
	landSquares := []*Square{
		{square.Rank + 1, square.File - 1},
		{square.Rank + 1, square.File},
		{square.Rank + 1, square.File + 1},
		{square.Rank, square.File - 1},
		{square.Rank, square.File + 1},
		{square.Rank - 1, square.File - 1},
		{square.Rank - 1, square.File},
		{square.Rank - 1, square.File + 1},
	}
	enemyKingSquare := board.GetKingSquare(!board.IsWhiteTurn)
	for _, landSquare := range landSquares {
		if !landSquare.IsValidBoardSquare() {
			continue
		}
		enemyKingRankDis := math.Abs(float64(int(enemyKingSquare.Rank) - int(landSquare.Rank)))
		enemyKingFileDis := math.Abs(float64(int(enemyKingSquare.File) - int(landSquare.File)))
		if enemyKingRankDis < 2 && enemyKingFileDis < 2 {
			continue
		}
		landPiece := board.GetPieceOnSquare(landSquare)
		if landPiece == EMPTY || landPiece.IsWhite() != piece.IsWhite() {
			move := Move{piece, square, landSquare, landPiece, make([]*Square, 0), EMPTY}
			kingMoves = append(kingMoves, &move)
		}
	}
	if canCastleKingside(board, square) {
		kingDestSquare := Square{square.Rank, square.File + 2}
		kingMoves = append(kingMoves, &Move{piece, square, &kingDestSquare, EMPTY, make([]*Square, 0), EMPTY})
	}
	if canCastleQueenside(board, square) {
		kingDestSquare := Square{square.Rank, square.File - 2}
		kingMoves = append(kingMoves, &Move{piece, square, &kingDestSquare, EMPTY, make([]*Square, 0), EMPTY})
	}
	kingMoves = filterMovesByKingSafety(board, kingMoves)
	return kingMoves
}

func IsLegalMove(board *Board, move *Move) bool {
	if move.Piece == EMPTY {
		return false
	}
	piece := board.GetPieceOnSquare(move.StartSquare)
	if move.Piece != piece {
		return false
	}
	if piece.IsWhite() != board.IsWhiteTurn {
		return false
	}

	var legalMoves []*Move
	var err error
	if piece.IsPawn() {
		legalMoves, err = GetLegalMovesForPawn(board, move.StartSquare)
	} else if piece.IsKnight() {
		legalMoves, err = GetLegalMovesForKnight(board, move.StartSquare)
	} else if piece.IsBishop() {
		legalMoves, err = GetLegalMovesForBishop(board, move.StartSquare)
	} else if piece.IsRook() {
		legalMoves, err = GetLegalMovesForRook(board, move.StartSquare)
	} else if piece.IsQueen() {
		legalMoves, err = GetLegalMovesForQueen(board, move.StartSquare)
	} else {
		legalMoves = GetLegalMovesForKing(board)
	}
	if err != nil {
		return false
	}
	for _, legalMove := range legalMoves {
		if legalMove.EqualTo(move) {
			return true
		}
	}
	return false
}

func GetLegalMovesFromOrigin(board *Board, square *Square) ([]*Move, error) {
	piece := board.GetPieceOnSquare(square)
	if piece == EMPTY {
		return make([]*Move, 0), nil
	}
	if board.IsWhiteTurn != piece.IsWhite() {
		return nil, fmt.Errorf("cannot generate moves, not player's turn")
	}
	if piece.IsPawn() {
		return GetLegalMovesForPawn(board, square)
	} else if piece.IsKnight() {
		return GetLegalMovesForKnight(board, square)
	} else if piece.IsBishop() {
		return GetLegalMovesForBishop(board, square)
	} else if piece.IsRook() {
		return GetLegalMovesForRook(board, square)
	} else if piece.IsQueen() {
		return GetLegalMovesForQueen(board, square)
	} else if piece.IsKing() {
		return GetLegalMovesForKing(board), nil
	}
	return nil, fmt.Errorf("unexpected piece %s while generating move", piece)
}

func HasLegalMove(board *Board) bool {
	for rank := uint8(1); rank < 9; rank++ {
		for file := uint8(1); file < 9; file++ {
			square := Square{rank, file}
			squareMoves, _ := GetLegalMovesFromOrigin(board, &square)
			if squareMoves != nil && len(squareMoves) > 0 {
				return true
			}
		}
	}
	return false
}

func GetLegalMovesBySquare(board *Board, stopAtFirst bool) (*[8][8][]*Move, uint8, error) {
	//for the active player, it returns:
	//	1. A 2d array which maps to board squares, where the element type
	//	   is a slice of legal moves for that piece in that square
	// 2. The total count of all legal moves
	var boardMoves [8][8][]*Move
	movesCount := uint8(0)
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		return nil, 0, fmt.Errorf("cannot generate moves by square on terminal board")
	}
	for rank := uint8(1); rank < 9; rank++ {
		for file := uint8(1); file < 9; file++ {
			square := Square{rank, file}
			movesBySquare, movesErr := GetLegalMovesFromOrigin(board, &square)
			if movesErr != nil {
				return nil, 0, fmt.Errorf("cannot generates moves by square: %s", movesErr)
			}
			boardMoves[rank-1][file-1] = movesBySquare
			movesCount += uint8(len(movesBySquare))

			if stopAtFirst && movesCount > 0 {
				return &boardMoves, movesCount, nil
			}
		}
	}
	return &boardMoves, movesCount, nil
}

func GetLegalMoves(board *Board) ([]*Move, error) {
	var moves = make([]*Move, 0)
	if board.Result != BOARD_RESULT_IN_PROGRESS {
		return moves, fmt.Errorf("cannot generate moves on terminal board")
	}
	for rank := uint8(1); rank < 9; rank++ {
		for file := uint8(1); file < 9; file++ {
			square := &Square{rank, file}
			squareMoves, movesErr := GetLegalMovesFromOrigin(board, square)
			if movesErr != nil {
				return make([]*Move, 0), fmt.Errorf("cannot generate moves on square %s on board %s: %s", square, board, movesErr)
			}
			moves = append(moves, squareMoves...)
		}
	}
	return moves, nil
}

// TODO: do we need square param here?
func canCastleKingside(board *Board, square *Square) bool {
	if board.IsWhiteTurn && !board.CanWhiteCastleKingside {
		return false
	} else if !board.IsWhiteTurn && !board.CanBlackCastleKingside {
		return false
	}
	if len(GetCheckingSquares(board, board.IsWhiteTurn)) > 0 {
		return false
	}
	piece := board.GetPieceOnSquare(square)

	kingRightSquare := Square{square.Rank, square.File + 1}
	kingRightTwoSquare := Square{square.Rank, square.File + 2}
	kingRightPiece := board.GetPieceOnSquare(&kingRightSquare)
	kingTwoRightPiece := board.GetPieceOnSquare(&kingRightTwoSquare)
	if kingRightPiece != EMPTY || kingTwoRightPiece != EMPTY {
		return false
	}

	boardBuilder := NewBoardBuilder().FromBoard(board)
	boardBuilder.WithPiece(EMPTY, square)
	boardBuilder.WithPiece(piece, &kingRightSquare)
	if len(GetCheckingSquares(boardBuilder.board, board.IsWhiteTurn)) > 0 {
		return false
	}
	boardBuilder.WithPiece(EMPTY, &kingRightSquare)
	boardBuilder.WithPiece(piece, &kingRightTwoSquare)
	if len(GetCheckingSquares(boardBuilder.board, board.IsWhiteTurn)) > 0 {
		return false
	}
	return true
}

// TODO: do we need square param here?
func canCastleQueenside(board *Board, square *Square) bool {
	if board.IsWhiteTurn && !board.CanWhiteCastleQueenside {
		return false
	} else if !board.IsWhiteTurn && !board.CanBlackCastleQueenside {
		return false
	}
	if len(GetCheckingSquares(board, board.IsWhiteTurn)) > 0 {
		return false
	}
	piece := board.GetPieceOnSquare(square)

	kingLeftSquare := Square{square.Rank, square.File - 1}
	kingLeftTwoSquare := Square{square.Rank, square.File - 2}
	kingLeftPiece := board.GetPieceOnSquare(&kingLeftSquare)
	kingTwoLeftPiece := board.GetPieceOnSquare(&kingLeftTwoSquare)
	if kingLeftPiece != EMPTY || kingTwoLeftPiece != EMPTY {
		return false
	}

	boardBuilder := NewBoardBuilder().FromBoard(board)
	boardBuilder.WithPiece(EMPTY, square)
	boardBuilder.WithPiece(piece, &kingLeftSquare)
	if len(GetCheckingSquares(boardBuilder.board, board.IsWhiteTurn)) > 0 {
		return false
	}
	boardBuilder.WithPiece(EMPTY, &kingLeftSquare)
	boardBuilder.WithPiece(piece, &kingLeftTwoSquare)
	if len(GetCheckingSquares(boardBuilder.board, board.IsWhiteTurn)) > 0 {
		return false
	}
	return true
}

func GetBoardFromMove(board *Board, move *Move) *Board {
	boardBuilder := NewBoardBuilder().FromBoard(board)

	UpdatePiecesFromMove(board, boardBuilder, move)
	UpdateBoardCounters(board, boardBuilder, move)
	UpdateBoardEnPassantSquare(board, boardBuilder, move)

	boardBuilder.WithIsWhiteTurn(!board.IsWhiteTurn)

	UpdateCastleRights(board, boardBuilder, move)
	repetitions := UpdateRepetitionsByFENMap(board, boardBuilder, move)
	UpdateBoardResult(board, boardBuilder, repetitions)

	return boardBuilder.Build()
}

func UpdatePiecesFromMove(lastBoard *Board, boardBuilder *BoardBuilder, move *Move) {
	movingPiece := lastBoard.GetPieceOnSquare(move.StartSquare)
	var landingPiece Piece
	if move.PawnUpgradedTo != EMPTY {
		landingPiece = move.PawnUpgradedTo
	} else {
		landingPiece = movingPiece
	}
	boardBuilder.WithPiece(EMPTY, move.StartSquare)
	boardBuilder.WithPiece(landingPiece, move.EndSquare)
	if lastBoard.OptEnPassantSquare != nil && movingPiece.IsPawn() && move.EndSquare.EqualTo(lastBoard.OptEnPassantSquare) {
		// en passant capture
		enPassantedPawnSquare := Square{
			move.StartSquare.Rank,
			lastBoard.OptEnPassantSquare.File,
		}
		boardBuilder.WithPiece(EMPTY, &enPassantedPawnSquare)
	}

	if move.IsCastles() {
		if move.Piece == WHITE_KING {
			if move.EndSquare.EqualTo(&Square{1, 7}) {
				boardBuilder.WithPiece(EMPTY, &Square{1, 8})
				boardBuilder.WithPiece(WHITE_ROOK, &Square{1, 6})
			} else if move.EndSquare.EqualTo(&Square{1, 3}) {
				boardBuilder.WithPiece(EMPTY, &Square{1, 1})
				boardBuilder.WithPiece(WHITE_ROOK, &Square{1, 4})
			}
		} else if move.Piece == BLACK_KING {
			if move.EndSquare.EqualTo(&Square{8, 7}) {
				boardBuilder.WithPiece(EMPTY, &Square{8, 8})
				boardBuilder.WithPiece(BLACK_ROOK, &Square{8, 6})
			} else if move.EndSquare.EqualTo(&Square{8, 3}) {
				boardBuilder.WithPiece(EMPTY, &Square{8, 1})
				boardBuilder.WithPiece(BLACK_ROOK, &Square{8, 4})
			}
		}
	}
}

func UpdateBoardCounters(lastBoard *Board, boardBuilder *BoardBuilder, move *Move) {
	if move.CapturedPiece != EMPTY || move.Piece.IsPawn() {
		boardBuilder.WithHalfMoveClockCount(0)
	} else {
		boardBuilder.WithHalfMoveClockCount(lastBoard.HalfMoveClockCount + 1)
	}
	if !lastBoard.IsWhiteTurn {
		boardBuilder.WithFullMoveCount(lastBoard.FullMoveCount + 1)
	}
}
func UpdateBoardEnPassantSquare(lastBoard *Board, boardBuilder *BoardBuilder, move *Move) {
	if move.DoesAllowEnPassant() {
		enPassantSquare := &Square{
			uint8(math.Min(float64(move.StartSquare.Rank), float64(move.EndSquare.Rank))) + 1,
			move.StartSquare.File,
		}
		boardBuilder.WithEnPassantSquare(enPassantSquare)
	} else {
		boardBuilder.WithEnPassantSquare(nil)
	}
}
func UpdateCastleRights(lastBoard *Board, boardBuilder *BoardBuilder, move *Move) {
	if !move.Piece.IsKing() && !move.Piece.IsRook() {
		return
	}
	if move.Piece.IsRook() {
		if lastBoard.CanWhiteCastleQueenside && move.StartSquare.EqualTo(&Square{1, 1}) {
			boardBuilder.WithCanWhiteCastleQueenside(false)
		} else if lastBoard.CanWhiteCastleKingside && move.StartSquare.EqualTo(&Square{1, 8}) {
			boardBuilder.WithCanWhiteCastleKingside(false)
		} else if lastBoard.CanBlackCastleQueenside && move.StartSquare.EqualTo(&Square{8, 1}) {
			boardBuilder.WithCanBlackCastleQueenside(false)
		} else if lastBoard.CanBlackCastleKingside && move.StartSquare.EqualTo(&Square{8, 8}) {
			boardBuilder.WithCanBlackCastleKingside(false)
		}
	} else if move.Piece.IsKing() {
		if move.Piece.IsWhite() {
			boardBuilder.WithCanWhiteCastleKingside(false)
			boardBuilder.WithCanWhiteCastleQueenside(false)
		} else {
			boardBuilder.WithCanBlackCastleKingside(false)
			boardBuilder.WithCanBlackCastleQueenside(false)
		}
	}
}

func UpdateRepetitionsByFENMap(lastBoard *Board, boardBuilder *BoardBuilder, move *Move) uint8 {
	castleRightsChanged := false
	castleRightsChanged = castleRightsChanged || lastBoard.CanWhiteCastleKingside != boardBuilder.board.CanWhiteCastleKingside
	castleRightsChanged = castleRightsChanged || lastBoard.CanWhiteCastleQueenside != boardBuilder.board.CanWhiteCastleQueenside
	castleRightsChanged = castleRightsChanged || lastBoard.CanBlackCastleKingside != boardBuilder.board.CanBlackCastleKingside
	castleRightsChanged = castleRightsChanged || lastBoard.CanBlackCastleQueenside != boardBuilder.board.CanBlackCastleQueenside

	if move.Piece.IsPawn() || move.CapturedPiece != EMPTY || castleRightsChanged {
		boardBuilder.WithRepetitionsByMiniFEN(make(map[string]uint8))
	}
	miniFEN := boardBuilder.board.ToMiniFEN()
	repetitions, _ := lastBoard.RepetitionsByMiniFEN[miniFEN]
	boardBuilder.WithMiniFENCount(miniFEN, repetitions+1)
	return repetitions + 1
}

func UpdateBoardResult(lastBoard *Board, boardBuilder *BoardBuilder, repetitions uint8) {
	if repetitions >= 3 {
		boardBuilder.WithResult(BOARD_RESULT_DRAW_BY_THREEFOLD_REPETITION)
	} else if boardBuilder.board.HalfMoveClockCount >= 50 {
		boardBuilder.WithResult(BOARD_RESULT_DRAW_BY_FIFTY_MOVE_RULE)
	} else if !boardBuilder.board.HasLegalNextMove() {
		checkingSquares := GetCheckingSquares(boardBuilder.board, boardBuilder.board.IsWhiteTurn)
		if len(checkingSquares) > 0 {
			if lastBoard.IsWhiteTurn {
				boardBuilder.WithResult(BOARD_RESULT_WHITE_WINS_BY_CHECKMATE)
			} else {
				boardBuilder.WithResult(BOARD_RESULT_BLACK_WINS_BY_CHECKMATE)
			}
		} else {
			boardBuilder.WithResult(BOARD_RESULT_DRAW_BY_STALEMATE)
		}
	} else if boardBuilder.board.IsForcedDrawByMaterial() {
		boardBuilder.WithResult(BOARD_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL)
	}
}
