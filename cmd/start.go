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

	"github.com/spf13/cobra"
)

func ParseExpr(cell string) string {
	return "EXPR"
}

func ParseText(cell string) string {
	return "TXT"
}

func ParseNum(cell string) string {
	return "NUM"
}

func ParseFile(table [][]string) [][]string{
	newTable := table

	rows := newTable[1:]
	for r, row := range rows {
		for c, cell := range row {
			if cell[0] == '=' {
				modifiedCell := ParseExpr(cell)
				newTable[1+r][c] = modifiedCell
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
				newTable[1+r][c] = modifiedCell
			}
		}
	}
	return newTable
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
		fmt.Println("Couldn't read the file.")
		os.Exit(1)
	}
	table = cleanTable(table)
	table = initRest(table)
	fmt.Print(table)
	newTable := ParseFile(table)
	outFile, err := os.Create(inputFile + "_result.csv")
	
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
