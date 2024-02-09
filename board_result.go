package chess

type BoardResult string

const (
	BOARD_RESULT_IN_PROGRESS                   BoardResult = "in_progress"
	BOARD_RESULT_WHITE_WINS_BY_CHECKMATE       BoardResult = "white_wins_by_checkmate"
	BOARD_RESULT_BLACK_WINS_BY_CHECKMATE       BoardResult = "black_wins_by_checkmate"
	BOARD_RESULT_DRAW_BY_STALEMATE             BoardResult = "draw_by_stalemate"
	BOARD_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL BoardResult = "draw_by_insufficient_material"
	BOARD_RESULT_DRAW_BY_THREEFOLD_REPETITION  BoardResult = "draw_by_threefold_repetition"
	BOARD_RESULT_DRAW_BY_FIFTY_MOVE_RULE       BoardResult = "draw_by_fifty_move_rule"
)
