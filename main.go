package main

import "fmt"
import "strconv"
import "bufio"
import "os"

func main() {
    game := newTicTacGo()

    game.start()
}

func newTicTacGo() ticTacGo {
    game := ticTacGo{
        board: newBoard(),
        player1: "x",
        player2: "o",
    }
    game.currentPlayer = game.player1

    return game
}

type ticTacGo struct {
    board *board
    player1 string
    player2 string
    currentPlayer string
}

func (game *ticTacGo) start() {
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Printf("Turn of %s\n", game.currentPlayer)
        fmt.Print(game.board.drawBoard())

        for {
            fmt.Println("Enter column:")
            scanner.Scan()
            err := scanner.Err()
            if err != nil {
                fmt.Println(err.Error())

                continue
            }
            column, err := strconv.Atoi(scanner.Text())
            if err != nil {
                fmt.Println(err.Error())

                continue
            }

            fmt.Println("Enter row:")
            scanner.Scan()
            err = scanner.Err()
            if err != nil {
                fmt.Println(err.Error())

                continue
            }
            row, err := strconv.Atoi(scanner.Text())
            if err != nil {
                fmt.Println(err.Error())

                continue
            }

            if err := game.board.setFieldValue(column - 1, row - 1, game.currentPlayer); err != nil {
                fmt.Println(err.Error())

                continue
            }
            
            break
        }

        winner, winningLocation, gameOver := game.board.getWinner()

        if winner != "" {
            fmt.Printf("%s wins in %s\n", winner, winningLocation)
            fmt.Print(game.board.drawBoard())

            break
        }

        if gameOver {
            fmt.Println("Game over. Nobody wins")
            fmt.Print(game.board.drawBoard())

            break
        }

        if game.currentPlayer == game.player1 {
            game.currentPlayer = game.player2
        } else {
            game.currentPlayer = game.player1
        }
    }
}

func newBoard() *board {
    return &board{
        fields: [3][3]field{
            [3]field{field{}, field{}, field{}},
            [3]field{field{}, field{}, field{}},
            [3]field{field{}, field{}, field{}},
        },
    }
}

type field struct {
    value string
}

type board struct {
    fields [3][3]field
}

func (board *board) setFieldValue(column int, row int, value string) error {
    if row < 0 || row >= len(board.fields) || column < 0 || column >= len(board.fields[0]) {
        return fmt.Errorf("Invalid field")
    }

    field := &board.fields[row][column]

    if field.value != "" {
        return fmt.Errorf("Field already set")
    }

    field.value = value

    return nil
}

func (board board) getWinner() (string, string, bool) {
    gameOver := true

    // Column
    for column := range board.fields[0] {
        fieldValues := []string{}

        for _, fields := range board.fields {
            field := fields[column]
            fieldValues = append(fieldValues, field.value)
        }

        if winner, ok := getWinnerForFields(fieldValues); winner != "" {
            return winner, fmt.Sprintf("column %d", column + 1), false
        } else if ok {
            gameOver = false
        }
    }

    // Row
    for row, fields := range board.fields {
        fieldValues := []string{}

        for _, field := range fields {
            fieldValues = append(fieldValues, field.value)
        }

        if winner, ok := getWinnerForFields(fieldValues); winner != "" {
            return winner, fmt.Sprintf("row %d", row + 1), false
        } else if ok {
            gameOver = false
        }
    }

    // Diagonal descending
    fieldValues := []string{}

    for row, fields := range board.fields {
        field := fields[row]
        fieldValues = append(fieldValues, field.value)
    }

    if winner, ok := getWinnerForFields(fieldValues); winner != "" {
        return winner, "diagonal descending", false
    } else if ok {
        gameOver = false
    }

    // Diagonal ascending
    fieldValues = []string{}

    for row, fields := range board.fields {
        field := fields[len(fields) - (row + 1)]
        fieldValues = append(fieldValues, field.value)
    }

    if winner, ok := getWinnerForFields(fieldValues); winner != "" {
        return winner, "diagonal ascending", false
    } else if ok {
        gameOver = false
    }

    return "", "", gameOver
}

func (board board) drawBoard() string {
    boardOutput := ""
    indexRow := "  "

    for column := range board.fields[0] {
        indexRow += strconv.Itoa(column + 1) + " "
    }

    boardOutput += indexRow + "\n"

    for row, fields := range board.fields {
        rowOutput := strconv.Itoa(row + 1) + " "

        for _, field := range fields {
            if field.value == "" {
                rowOutput += " "
            } else {
                rowOutput += field.value
            }

            rowOutput += " "
        }

        boardOutput += rowOutput + "\n"
    }

    return boardOutput
}


func countStrings(strings []string) map[string]int {
    count := map[string]int{}

    for _, string := range strings {
        if _, ok := count[string]; !ok {
            count[string] = 0;
        }

        count[string]++
    }

    return count;
}

func getWinnerForFields(fieldValues []string) (string, bool) {
    fieldValuesCount := countStrings(fieldValues)

    if _, ok := fieldValuesCount[""]; !ok {
        if len(fieldValuesCount) == 1 {
            return fieldValues[0], true
        }

        return "", false
    } else if len(fieldValuesCount) > 2 {
        return "", false
    }

    return "", true
}
