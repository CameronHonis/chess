package chess_test

import (
	"bufio"
	"fmt"
	"github.com/CameronHonis/chess"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

const PRINT_LEAFS = true

var scanner *bufio.Scanner

func perft(board *chess.Board, depth int) int {
	moves, movesErr := chess.GetLegalMoves(board)
	if movesErr != nil {
		log.Fatalf("error getting moves for board %s: %s", board, movesErr)
	}
	leaf := depth == 1
	nodeCnt := 0
	for _, move := range moves {
		nextBoard := chess.GetBoardFromMove(board, move)
		if leaf {
			if PRINT_LEAFS {
				fmt.Printf("leaf %s\n", nextBoard.ToFEN())
			}
			nodeCnt++
		} else {
			nodeCnt += perft(nextBoard, depth-1)
		}
	}
	return nodeCnt
}

//func compareNodeCounts(board *chess.Board, nodeCountByDepth []int) {
//	depth := len(nodeCountByDepth)
//	var currBoards = []*chess.Board{board}
//	var nextBoards = make([]*chess.Board, 0)
//	for i := 0; i < depth; i++ {
//		expNodes := nodeCountByDepth[i]
//		for _, currBoard := range currBoards {
//			moves, movesErr := chess.GetLegalMoves(board)
//			if movesErr != nil {
//				log.Fatalf("error getting moves for board %s: %s", currBoard, movesErr)
//			}
//			for _, move := range moves {
//				nextBoard := chess.GetBoardFromMove(board, move)
//				nextBoards = append(nextBoards, nextBoard)
//			}
//		}
//
//		log.Println(i, len(nextBoards))
//		if len(nextBoards) != expNodes {
//			log.Fatalf("incorrect number of game nodes from %s at depth %d, actual %d vs expected %d", board, i, len(nextBoards), expNodes)
//		}
//		currBoards = nextBoards
//		nextBoards = make([]*chess.Board, 0)
//	}
//}

func TestPerf(t *testing.T) {
	file, err := os.Open("./perft")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	//var rawTxt []string
	for scanner.Scan() {
		splitTxt := strings.Split(scanner.Text(), ";")
		fen := splitTxt[0]
		board, boardErr := chess.BoardFromFEN(fen)
		if boardErr != nil {
			log.Fatalf("could not construct board from FEN %s:\n\t%s", fen, boardErr)
		}
		for _, perftStr := range splitTxt[1:] {
			perftStrSplit := strings.Split(perftStr, " ")
			depthStr := perftStrSplit[0][1:]
			depth, parseDepthErr := strconv.Atoi(depthStr)
			if parseDepthErr != nil {
				log.Fatalf("could not parse depth from %s:\n\t%s", depthStr, parseDepthErr)
			}

			expNodeCntStr := perftStrSplit[1]
			expNodeCnt, parseNodeCntErr := strconv.Atoi(expNodeCntStr)
			if parseNodeCntErr != nil {
				log.Fatalf("could not parse expNodeCnt from %s:\n\t%s", expNodeCntStr, parseNodeCntErr)
			}
			actNodeCnt := perft(board, depth)
			if actNodeCnt != expNodeCnt {
				log.Fatalf("node count mismatch, actual (%d) vs expected (%d)", actNodeCnt, expNodeCnt)
			}
		}
	}

	_ = file.Close()
}
