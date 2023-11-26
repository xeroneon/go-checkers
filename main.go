package main

import (
	"bufio"
	"checkers/utils"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EmptySquare    = ' '
	BlackSquare    = '█'
	PlayerOnePiece = 'O'
	PlayerTwoPiece = 'X'
	PlayerOneKing  = 'o'
	PlayerTwoKing  = 'x'
)

type Player int

const (
	None Player = iota
	PlayerOne
	PlayerTwo
)

const (
	StartGame = "Start Game"
	ExitGame  = "Exit Game"
)

type Cell struct {
	Piece  rune
	Player Player
	IsKing bool
}

type Board [8][8]Cell

var (
	PlayerOnePiecesLeft = 12
	PlayerTwoPiecesLeft = 12
	PlayerOneWins       = 0
	PlayerTwoWins       = 0
)

func main() {
	for {
		selection := askForSelection()
		fmt.Println(selection)
		switch selection {
		case StartGame:
			utils.ClearScreen()
			startGame()
		case ExitGame:
			os.Exit(0)
		default:
			os.Exit(0)
		}
	}
}

func askForSelection() string {
	reader := bufio.NewReader(os.Stdin)
	options := []string{"Start Game", "Exit Game"}

	fmt.Println("Choose an option")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	fmt.Println("Enter number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)

	if err != nil || choice < 1 || choice > len(options) {
		utils.ClearScreen()
		utils.PrintErr("Invalid selection")
		return askForSelection()
	} else {
		return options[choice-1]
	}
}

func startGame() {
	PlayerOnePiecesLeft = 12
	PlayerTwoPiecesLeft = 12
	var board Board
	initializeBoard(&board)

	scanner := bufio.NewScanner(os.Stdin)
	currentPlayer := PlayerOne

	for {
		if currentPlayer == PlayerTwo {
			var playerTwoPieces []string
			var openSpaces []string
			for i, row := range board {
				for j, col := range row {
					if col.Player == PlayerTwo {
						X := rune(j + 'A')
						Y := i + 1
						playerTwoPieces = append(playerTwoPieces, fmt.Sprintf("%c%d", X, Y))
					}
					if col.Player == None && col.Piece == BlackSquare {
						X := rune(j + 'A')
						Y := i + 1
						openSpaces = append(openSpaces, fmt.Sprintf("%c%d", X, Y))
					}
				}
			}
			rand.Seed(time.Now().UnixNano())
			startCell := playerTwoPieces[rand.Intn(len(playerTwoPieces))]
			endCell := openSpaces[rand.Intn(len(openSpaces))]
			move := fmt.Sprintf("%v %v", startCell, endCell)
			if !makeMove(&board, move, currentPlayer) {
				continue
			}
		}
		if currentPlayer == PlayerOne {
			printBoard(board)
			fmt.Println("Player 1 Piece: O (o for kings)")
			fmt.Println("Player 2 Piece: X (x for kings)")
			fmt.Printf("Player %d's turn. Enter move (e.g., 'E4 F4'): ", currentPlayer)
			scanner.Scan()
			move := scanner.Text()

			if !makeMove(&board, move, currentPlayer) {
				utils.ClearScreen()
				red := "\033[31m"
				reset := "\033[0m"
				fmt.Println()
				fmt.Println(red + "Invalid move. Try again." + reset)
				continue
			}

		}
		utils.ClearScreen()
		if PlayerOnePiecesLeft == 0 {
			PlayerTwoWins = PlayerTwoWins + 1
			break
		}
		if PlayerTwoPiecesLeft == 0 {
			PlayerOneWins = PlayerOneWins + 1
			break
		}

		if currentPlayer == PlayerOne {
			currentPlayer = PlayerTwo
		} else {
			currentPlayer = PlayerOne
		}

	}
}

func initializeBoard(board *Board) {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 == 0 {
				if row < 3 {
					board[row][col] = Cell{Piece: PlayerOnePiece, Player: PlayerOne, IsKing: false}
				} else if row > 4 {
					board[row][col] = Cell{Piece: PlayerTwoPiece, Player: PlayerTwo, IsKing: false}
				} else {
					board[row][col] = Cell{Piece: BlackSquare, Player: None, IsKing: false}
				}
			} else {
				board[row][col] = Cell{Piece: EmptySquare, Player: None, IsKing: false}
			}
		}
	}
}

func printBoard(board Board) {
	fmt.Println()
	for row := 0; row < 8; row++ {
		fmt.Printf("%d ", row+1)

		for col := 0; col < 8; col++ {
			fmt.Printf("%c ", board[row][col].Piece)
		}
		fmt.Println()
	}
	fmt.Print("  ")
	for col := 0; col < 8; col++ {
		fmt.Printf("%c ", 'A'+col)
	}
	fmt.Println()
	fmt.Println()
}

func makeMove(board *Board, move string, player Player) bool {
	// INFO: turn move from string to array coordinates
	parts := strings.Split(move, " ")
	if len(parts) != 2 {
		return false
	}

	start, end := parts[0], parts[1]
	startCol, startRow := utils.ParsePosition(start)
	endCol, endRow := utils.ParsePosition(end)

	// INFO: dont allow the move if the coordinates dont exist
	if startCol < 0 || startCol >= 8 || startRow < 0 || startRow >= 8 ||
		endCol < 0 || endCol >= 8 || endRow < 0 || endRow >= 8 {
		return false
	}

	startCell := &board[startRow][startCol]
	endCell := &board[endRow][endCol]

	// INFO: dont allow players to move anything that isn't their piece
	if startCell.Player != player {
		return false
	}

	if endCell.Piece != BlackSquare {
		return false
	}

	if !isValidMove(startRow, startCol, endRow, endCol, player, *board) {
		return false
	}

	var pieceType rune

	isKing := startCell.IsKing
	// INFO: this is to figure out if the piece needs to be a king or not
	if player == PlayerOne {
		if isKing {
			pieceType = PlayerOneKing
		} else {
			pieceType = PlayerOnePiece
		}
	} else {
		if isKing {
			pieceType = PlayerTwoKing
		} else {
			pieceType = PlayerTwoPiece
		}
	}

	if isKing || utils.Abs(endRow-startRow) == 2 {
		jumped, row, col := isJumping(startRow, startCol, endRow, endCol, player, *board)
		fmt.Println(jumped, row, col)
		if jumped {
			board[row][col].Piece = BlackSquare
			board[row][col].Player = None
			board[row][col].IsKing = false
			if player == PlayerOne {
				PlayerTwoPiecesLeft = PlayerTwoPiecesLeft - 1
			} else {
				PlayerOnePiecesLeft = PlayerOnePiecesLeft - 1
			}
		}
	}

	endCell.Piece = pieceType
	endCell.Player = player

	// INFO: this is to see if the piece should be upgraded to a king
	if !isKing {
		if (player == PlayerOne && endRow == 7) || (player == PlayerTwo && endRow == 0) {
			endCell.IsKing = true
			if player == PlayerOne {
				endCell.Piece = PlayerOneKing
			} else {
				endCell.Piece = PlayerTwoKing
			}
		}
	}

	startCell.Piece = BlackSquare
	startCell.Player = None
	startCell.IsKing = false
	return true
}

func isJumping(startX, startY, endX, endY int, player Player, board Board) (bool, int, int) {
	if utils.Abs(endY-startY) == 2 {
		jumpedCol := startX + (endX-startX)/2
		jumpedRow := startY + (endY-startY)/2
		return true, jumpedCol, jumpedRow
	} else {
		return false, 0, 0
	}
}

func isValidMove(startRow, startCol, endRow, endCol int, player Player, board Board) bool {
	isKing := board[startRow][startCol].IsKing
	deltaX := utils.Abs(endRow - startRow)
	deltaY := utils.Abs(endCol - startCol)
	_, jumpedRow, jumpedCol := isJumping(startRow, startCol, endRow, endCol, player, board)

	jumpedCell := board[jumpedRow][jumpedCol].Player

	if isKing {
		return (deltaX == 1 && deltaY == 1) ||
			(jumpedCell != player && jumpedCell != None && deltaX == 2 && deltaY == 2)
	} else {
		var forwardMove bool
		if player == PlayerOne {
			forwardMove = endRow > startRow
		} else {
			forwardMove = endRow < startRow
		}

		return forwardMove && ((deltaX == 1 && deltaY == 1) || (jumpedCell != player && jumpedCell != None && deltaX == 2 && deltaY == 2))
	}
}
