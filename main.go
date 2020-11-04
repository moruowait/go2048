package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// 初始化棋盘 4x4
	var board = NewBoard(4)
	// 随机初始化 N 个数字
	board.InitStartData(4)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if board.IsOver() {
			fmt.Println("game over！")
			return
		}
		switch scanner.Text() {
		case "j", "l", "k", "i":
			board.Move(scanner.Text())
		default:
			fmt.Println("error dir")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
