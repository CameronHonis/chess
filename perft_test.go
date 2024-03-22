package chess_test

import (
	"bufio"
	"fmt"
	"github.com/CameronHonis/chess"
	. "github.com/onsi/ginkgo/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const PRINT_LEAFS = true
const FOCUS_TEST_IDX = 1
const MAX_DEPTH = 3

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
		if leaf && PRINT_LEAFS {
			fmt.Printf("leaf %s\n", nextBoard.ToFEN())
		}
		if nextBoard.Result != chess.BOARD_RESULT_IN_PROGRESS {
			if leaf {
				nodeCnt++
			}
			continue
		}
		if leaf {
			nodeCnt++
		} else {
			nodeCnt += perft(nextBoard, depth-1)
		}
	}
	return nodeCnt
}

func parsePerftLine(line string) (fen string, depthNodeCntPairs [][2]int) {
	splitTxt := strings.Split(line, ";")
	fen = splitTxt[0]
	depthNodeCntPairs = make([][2]int, 0)
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
		depthNodeCntPairs = append(depthNodeCntPairs, [2]int{depth, expNodeCnt})
	}
	return fen, depthNodeCntPairs
}

func perft_from_file() {
	file, err := os.Open("./perft")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	testIdx := 0
	for scanner.Scan() {
		currTestIdx := testIdx
		testIdx++

		shouldSkipTest := FOCUS_TEST_IDX >= 0 && FOCUS_TEST_IDX != currTestIdx
		if shouldSkipTest {
			continue
		}
		line := scanner.Text()
		fmt.Printf("%s\n", line)
		fen, depthNodeCntPairs := parsePerftLine(line)

		board, boardErr := chess.BoardFromFEN(fen)
		if boardErr != nil {
			log.Fatalf("could not construct board from FEN %s:\n\t%s", fen, boardErr)
		}

		for _, depthNodeCntPair := range depthNodeCntPairs {
			depth := depthNodeCntPair[0]
			expNodeCnt := depthNodeCntPair[1]

			start := time.Now()
			actNodeCnt := perft(board, depth)
			if actNodeCnt != expNodeCnt {
				log.Fatalf("node count mismatch at depth %d,  actual (%d) vs expected (%d)", depth, actNodeCnt, expNodeCnt)
			} else {
				elapsed := time.Since(start)
				fmt.Printf("depth %d passed in %s\n", depth, elapsed)
			}
		}
	}

	_ = file.Close()
}

var _ = It("perft", func() {
	perft_from_file()
})

var _ = FIt("temp", func() {
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	board, _ := chess.BoardFromFEN(fen)
	perft(board, 3)
})
