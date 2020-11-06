package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Board struct {
	Size  int
	Grids [][]int
}

// 初始化棋盘
func NewBoard(size int) *Board {
	return &Board{
		Size:  size,
		Grids: initGrids(size),
	}
}

func initGrids(size int) [][]int {
	var grids = make([][]int, size)
	for row := 0; row < size; row++ {
		grids[row] = make([]int, size)
	}
	return grids
}

func (b *Board) InitStartData(num int) {
	for i := 0; i < num; i++ {
		b.fillRandomData()
	}
	b.PrettyPrintBoard()
}

// 随机填充数据
func (b *Board) fillRandomData() {
	if !b.hasEmptyCells() {
		return
	}
	cell, err := b.SelectCell()
	if err != nil {
		return
	}
	// 更新格子
	b.updateGrid(cell.Row, cell.Column, b.randValue())
}

func (b *Board) randValue() int {
	rand.Seed(time.Now().UnixNano())
	var n = rand.Intn(10)
	if n < 9 {
		return 2
	}
	return 4
}

func (b *Board) updateGrid(row, column, value int) {
	b.Grids[row][column] = value
}

// 空格子
func (b *Board) emptyCells() []*Cell {
	var cells []*Cell
	for row, rowGrids := range b.Grids {
		for column, grid := range rowGrids {
			if grid != 0 {
				continue
			}
			cells = append(cells, &Cell{
				Row:    row,
				Column: column,
			})
		}
	}
	return cells
}

// 选中一个空格子
func (b *Board) SelectCell() (*Cell, error) {
	var cells = b.emptyCells()
	if len(cells) == 0 {
		return nil, errors.New("empty cells")
	}
	return cells[rand.Intn(len(cells))], nil
}

// 判断是否还有空格子
func (b *Board) hasEmptyCells() bool {
	return len(b.emptyCells()) != 0
}

// 打印棋盘
func (b *Board) PrettyPrintBoard() {
	for _, rowGrids := range b.Grids {
		for _, grid := range rowGrids {
			if grid == 0 {
				fmt.Printf("_ ")
			} else {
				fmt.Printf("%d ", grid)
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

// 判断是否发生改变
func (b *Board) different(grids [][]int) bool {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if grids[i][j] != b.Grids[i][j] {
				return true
			}
		}
	}
	return false
}

// 移动，j-左移，l-右移，i-上移，k-下移
func (b *Board) Move(dir string) {
	// 记录原有棋盘
	var grids = make([][]int, b.Size)
	for i := 0; i < b.Size; i++ {
		grids[i] = append(grids[i], b.Grids[i]...)
	}

	var list = make([][]int, b.Size, b.Size)
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			switch dir {
			case "j": // 左移
				list[i] = append(list[i], b.Grids[i][j])
			case "l": // 右移
				list[i] = append(list[i], b.Grids[i][b.Size-j-1])
			case "i": // 上移
				list[j] = append(list[j], b.Grids[i][j])
			case "k": // 下移
				list[j] = append(list[j], b.Grids[i][b.Size-j-1])
			}
		}
	}

	list = b.moveClose(list)
	list = b.Combine(list)
	list = b.moveClose(list)

	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			switch dir {
			case "j": // 左移
				b.Grids[i][j] = list[i][j]
			case "l": // 右移
				b.Grids[i][b.Size-j-1] = list[i][j]
			case "i": // 上移
				b.Grids[i][j] = list[j][i]
			case "k": // 下移
				b.Grids[i][b.Size-j-1] = list[j][i]
			}
		}
	}

	// 如果移动完没有变化，则不添加随机值
	if b.different(grids) {
		// 重新添加一个随机值
		b.fillRandomData()
	}
	// 打印棋盘
	b.PrettyPrintBoard()

	return
}

func (b *Board) Combine(list [][]int) [][]int {
	for i, l := range list {
		for j := 0; j < len(l)-1; j += 1 {
			if l[j] == l[j+1] && l[j] != 0 {
				list[i][j] += list[i][j+1]
				list[i][j+1] = 0
			}
		}
	}
	return list
}

// 将 [2,0,0,2] 移动为 [2,2,0,0]
func (b *Board) moveClose(list [][]int) [][]int {
	for i, l := range list {
		var cnt int
		for _, v := range l {
			if v != 0 {
				list[i][cnt] = v
				cnt++
			}
		}
		for j := cnt; j < len(l); j++ {
			list[i][j] = 0
		}
	}
	return list
}

// 游戏是否结束：可用格子为空且所有格子上下左右值不等
func (b *Board) IsOver() bool {
	// 如果有空余的格子
	if b.hasEmptyCells() {
		return false
	}
	// 如果没有空余的格子
	// 左右不等
	for i := 0; i < b.Size; i++ {
		for j := 1; j < b.Size; j++ {
			if b.Grids[i][j] == b.Grids[i][j-1] {
				return false
			}
		}
	}
	// 上下不相等
	for j := 0; j < b.Size; j++ {
		for i := 1; i < b.Size; i++ {
			if b.Grids[i][j] == b.Grids[i-1][j] {
				return false
			}
		}
	}
	return true
}
