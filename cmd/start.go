/*
Built by Adithya Kaushik
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"unicode"
	"math"
	"strconv"
	"github.com/spf13/cobra"
)

type Pos struct {
	Row int
	Col int
}

func isBinaryOp(c string) bool {
	switch(c) {
		case "+", "-", "/", "*", "^" :
			return true
	}
	return false
}

func GetPrecedence(op string) int {
	switch(op) {
		case "^" : return 3
		case "/", "*" : return 2
		case "+", "-" : return 1
		default : return 0
	}
}

func InfixToPostfix(data string) []string {
	var infix []string
	var token string = ""
	for i := 1; i < len(data); i++ {
		if string(data[i]) == "(" {
			if token != "" { infix = append(infix, token); token = "" }
			infix = append(infix, "(")
		} else if string(data[i]) == ")" {
			if token != "" { infix = append(infix, token); token = "" }
			infix = append(infix, ")")
		} else if !isBinaryOp(string(data[i])) {
			token = token + string(data[i])
		} else {
			infix = append(infix, token)
			infix = append(infix, string(data[i]))
			token = ""
		}
	}
	infix = append(infix, token)	
	s := StringStack{}
	var postfix []string = []string{}
	for _, val := range infix {
		val = strings.Trim(val, " \n\t\r")
		if !isBinaryOp(val) && val != "(" && val != ")" {
			postfix = append(postfix, val)
		} else if val == "(" {
			s.Push(val)
		} else if val == ")" {
			for ; !s.IsEmpty() && s.Top() != "("; {
				popped, _ := s.Pop()
				postfix = append(postfix, popped)
			}
			s.Pop()
		} else { //operator
			if s.IsEmpty() {
				s.Push(val)
			} else {
				for ; !s.IsEmpty() && s.Top() != "(" ; {
					top := s.Top()
					p1 := GetPrecedence(top)
					p2 := GetPrecedence(val)

					if p1 > p2 || (p2 == p1 && val != "^") {
						v, _ := s.Pop()
						postfix = append(postfix, v)
					} else {
						break
					}
				}
				s.Push(val)
			}
		}
	}

	for ; !s.IsEmpty(); {
		postfix = append(postfix, s.Top())
		s.Pop()
	}

	return postfix

}

func Strip(data string) ([]string, []string) {
	var members []string
	var ops []string
	var token string = ""
	for i := 1; i < len(data); i++ {
		if string(data[i]) == "(" {
			if token != "" { members = append(members, token); token = "" }
			members = append(members, "(")
		} else if string(data[i]) == ")" {
			if token != "" { members = append(members, token); token = "" }
			members = append(members, ")")
		} else if !isBinaryOp(string(data[i])) {
			token = token + string(data[i])
		} else {
			ops = append(ops, string(data[i]))
			members = append(members, token)
			token = ""
		}
	}
	members = append(members, token)
	return members, ops

}

func isRef(data string) bool {
    if len(data) < 2 { return false }

    i := 0
    for i < len(data) && unicode.IsLetter(rune(data[i])) {
        i++
    }
    if i == 0 || i == len(data) {
        return false
    }
    for i < len(data) {
        if !unicode.IsDigit(rune(data[i])) {
            return false
        }
        i++
    }

    return true
}

func GetType(data string) string {
    if len(data) == 0 { return "float" }

    if data[0] == '"' {
        return "string"
    } else if data[0] == '=' {
        return "formula"
    } else if isRef(data) {
        return "ref"
    } else if unicode.IsLetter(rune(data[0])) {
        return "string"
    }
    
    return "float"
}

func UnpackRef(data string) (int, int) {
	var col int = 0
	var row int = 0

	for _, c := range data {
		if unicode.IsLetter(c) {
			col = 26 * col + int(unicode.ToUpper(c) - 'A' + 1)
		} else {
			row = 10 * row + int(c - '0')
		}
	}
	return row -1, col-1

}

func CheckValidRef(data string, table [][]string, vis map[Pos]bool) (string, string, string) {
    row, col := UnpackRef(data)

    if row < 0 || col < 0 || row >= len(table) || col >= len(table[0]) {
        return "#BOUNDS", "ERR", "#BOUNDS"
    }

    p := Pos{Row: row, Col: col}
    refVal := table[row][col]
    pType := GetType(refVal)

    if vis[p] {
        return "#CYCLE", "ERR", "#CYCLE"
    }

    if pType == "ref" || pType == "formula" {
        vis[p] = true
        defer delete(vis, p)
    }

    if pType == "ref" {
        return CheckValidRef(refVal, table, vis)
    }

    if pType == "formula" {
        resolvedVal := ParseExpr(refVal, table, row, col, vis)
        return resolvedVal, "float", ""
    }

    return refVal, pType, ""
}

func EvalExpr(postfix []string, table [][]string, r int, c int, vis map[Pos]bool) string {
    s := StringStack{}
    vis[Pos{Row: r, Col: c}] = true

    for _, val := range postfix {
        if !isBinaryOp(val) {
            s.Push(val)
        } else {
            B, _ := s.Pop()
            A, _ := s.Pop()

            valA, valB := A, B
            typeA, typeB := GetType(A), GetType(B)
            var errA, errB string

            if typeA == "ref" {
                valA, typeA, errA = CheckValidRef(A, table, vis)
            }
            if typeB == "ref" {
                valB, typeB, errB = CheckValidRef(B, table, vis)
            }

            if errA != "" { return errA }
            if errB != "" { return errB }
            
            if strings.HasPrefix(valA, "#") { return valA }
            if strings.HasPrefix(valB, "#") { return valB }

            if typeA != "float" || typeB != "float" {
                return "#TYPE!"
            }

            useA, _ := strconv.ParseFloat(valA, 64)
            useB, _ := strconv.ParseFloat(valB, 64)

            var res float64
            switch val {
            case "+":
                res = useA + useB
            case "-":
                res = useA - useB
            case "*":
                res = useA * useB
            case "/":
                if useB == 0 {
                    return "#DIV/0!"
                }
                res = useA / useB
            case "^":
                res = math.Pow(useA, useB)
            }
            s.Push(strconv.FormatFloat(res, 'f', 2, 64))
        }
    }

    if s.IsEmpty() {
        return "0"
    }

    finalVal, _ := s.Pop()
    if GetType(finalVal) == "ref" {
        resolvedVal, _, errStr := CheckValidRef(finalVal, table, vis)
        if errStr != "" {
            return errStr
        }
        return resolvedVal
    }

    return finalVal
}

func ParseExpr(cell string, table [][]string, r int, c int, vis map[Pos]bool) string {
    postfix := InfixToPostfix(cell)
    result := EvalExpr(postfix, table, r, c, vis)
    return result
}

func ParseText(cell string) string {
	return cell
}

func ParseNum(cell string) string {
	return cell
}

func ParseFile(table [][]string) [][]string{
	for r := 1; r < len(table); r++ {
		for c := 0; c < len(table[r]); c++ {
			cell := table[r][c]
			if cell == "" {
                table[r][c] = "0"
                continue
            }
			if cell[0] == '=' {
				var vis map[Pos]bool = make(map[Pos]bool)
				modifiedCell := ParseExpr(cell, table, r, c, vis)
				table[r][c] = modifiedCell
			} else {
				flag := false
				var modifiedCell string
				for i := 0; i < len(cell); i++ {
					if unicode.IsLetter(rune(cell[i])) {
						modifiedCell = ParseText(cell)
						flag = true
						break
					}
				}
				if !flag {
					modifiedCell = ParseNum(cell)
				}
				table[r][c] = modifiedCell
			}
		}
	}
	return table
}

func initRest(table [][]string) [][]string {
	var maxCols int 
	for _, row := range table {
		maxCols = max(maxCols, len(row))
	}
	currHeaderLength := len(table[0])
	if currHeaderLength < maxCols {
		Header := make([]string, maxCols)
		colStart := string(len(table[0]) + 'A')

		copy(Header, table[0])
		for i := currHeaderLength; i < maxCols; i++ {
			colStart = string(i + 'A')
			Header[i] = colStart
		}
		table[0] = Header
	}

	for i := 1; i < len(table); i++ {
		currLen := len(table[i])
		if currLen < maxCols {
			newRow := make([]string, maxCols)
			copy(newRow, table[i])
			for j := currLen; j < maxCols; j++ {
				newRow[j] = "0"
			}
			table[i] = newRow
		}
	}

	return table
}



func cleanTable(table [][]string) [][]string {
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] == "" {
				table[i][j] = "0"
			}
			table[i][j] = strings.TrimRight(table[i][j], " \t\n\r")
		}
	}

	return table
}


func UnpackFile(inputFile string) {
	file, err := os.OpenFile(inputFile, os.O_RDONLY, 0644)
	
	if err != nil {
		fmt.Println("Error in opening the file!")
		os.Exit(1)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	r.Comma = '|'

	table, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	table = cleanTable(table)
	table = initRest(table)
	newTable := ParseFile(table)
	inp, _, _ := strings.Cut(inputFile, ".csv")
	outFile, err := os.Create(inp + "_result.csv")
	
	if err != nil {
		fmt.Println("Error in output file creation.")
		os.Exit(1)
	}
	defer outFile.Close()

	w := csv.NewWriter(outFile)
	w.Comma = '|'
	w.WriteAll(newTable)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts GoExcel",
	Long: `Starts GoExcel, expects an input CSV file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide an input file")
			os.Exit(1)
		} else if len(args) > 1 {
			fmt.Println("Too many files! I require just one")
			os.Exit(1)
		}

		UnpackFile(args[0])
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
