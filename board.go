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

// 移动，j-左移，l-右移，i-上移，k-下移
func (b *Board) Move(dir string) {
	switch dir {
	case "j":
		b.MoveLeft()
	case "l":
		b.MoveRight()
	case "i":
		b.MoveUp()
	case "k":
		b.MoveDown()
	}
	// 重新添加一个随机值
	b.fillRandomData()
	// 打印棋盘
	b.PrettyPrintBoard()
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

func (b *Board) exchange(list [][]int) [][]int {
	for i := 0; i < b.Size; i++ {
		for j := i; j < b.Size; j++ {
			list[i][j], list[j][i] = list[j][i], list[i][j]
		}
	}
	return list
}

func (b *Board) exchangeRow(list [][]int) [][]int {
	for row := 0; row < len(list)/2; row++ {
		list[row], list[b.Size-row-1] = list[b.Size-row-1], list[row]
	}
	return list
}

func (b *Board) exchangeColumn(list [][]int) [][]int {
	for row := 0; row < len(list); row++ {
		for i := 0; i < len(list[row])-1; i += 2 {
			list[row][i], list[row][b.Size-i-1] = list[row][b.Size-i-1], list[row][i]
		}
	}
	return list
}

func (b *Board) MoveLeft() {
	fmt.Println("左移")
	var list = make([][]int, b.Size, b.Size)
	for row, rowGrids := range b.Grids {
		for _, grid := range rowGrids {
			list[row] = append(list[row], grid)
		}
	}
	list = b.moveClose(list)
	list = b.Combine(list)
	list = b.moveClose(list)

	// 还原
	for i, l := range list {
		for j, v := range l {
			b.Grids[i][j] = v
		}
	}
}

func (b *Board) MoveRight() {
	fmt.Println("右移")
	var list = make([][]int, b.Size, b.Size)
	for row, rowGrids := range b.Grids {
		for j := b.Size - 1; j >= 0; j-- {
			list[row] = append(list[row], rowGrids[j])
		}
	}
	list = b.moveClose(list)
	list = b.Combine(list)
	list = b.moveClose(list)

	// 还原
	for i, l := range list {
		for j := b.Size - 1; j >= 0; j-- {
			b.Grids[i][b.Size-j-1] = l[j]
		}
	}
}

func (b *Board) MoveUp() {
	fmt.Println("上移")

	var list = make([][]int, b.Size, b.Size)
	for row, rowGrids := range b.Grids {
		for _, grid := range rowGrids {
			list[row] = append(list[row], grid)
		}
	}
	// 交换
	list = b.exchange(list)
	// 交换行
	list = b.exchangeRow(list)

	list = b.moveClose(list)
	list = b.Combine(list)
	list = b.moveClose(list)

	// 交换行
	list = b.exchangeRow(list)
	// 交换
	list = b.exchange(list)

	// 还原
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			b.Grids[i][j] = list[i][j]
		}
	}
}

func (b *Board) MoveDown() {
	fmt.Println("下移")

	var list = make([][]int, b.Size, b.Size)
	for row, rowGrids := range b.Grids {
		for _, grid := range rowGrids {
			list[row] = append(list[row], grid)
		}
	}

	// 交换行
	list = b.exchange(list)
	list = b.exchangeColumn(list)

	list = b.moveClose(list)
	list = b.Combine(list)
	list = b.moveClose(list)

	// 交换行
	list = b.exchangeColumn(list)
	list = b.exchange(list)
	// 还原
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			b.Grids[i][j] = list[i][j]
		}
	}
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
