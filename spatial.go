package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*===============================================================
* Functions to manipulate a "field" of cells --- the main data
* that must be managed by this program.
*==============================================================*/

// The data stored in a single cell of a field
type Cell struct {
	kind  string
	score float64
	newKind string
}


// createField creates a new field of the ysize rows and xsize columns,
// so that field[r][c] gives the Cell at position (r,c).
func createField(rsize, csize int) [][]Cell {
	f := make([][]Cell, rsize)
	for i := range f {
		f[i] = make([]Cell, csize)
	}
	return f
}


// inField returns true iff (row,col) is a valid cell in the field
func inField(field [][]Cell, row, col int) bool {
	return row >= 0 && row < len(field) && col >= 0 && col < len(field[0])
}


// readFieldFromFile opens the given file and read the initial
// values for the field. The first line of the file will contain
// two space-separated integers saying how many rows and columns
// the field should have:
//    10 15
// each subsequent line will consist of a string of Cs and Ds, which
// are the initial strategies for the cells:
//    CCCCCCDDDCCCCCC
//
// If there is ever an error reading, this function should cause the
// program to quit immediately.
func readFieldFromFile(filename string) [][]Cell {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: something went wrong opening the file.")
		fmt.Println("Probably you gave the wrong filename.")
	}
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())

	}
	var dimensions string = lines[0]
	var rowcol []string  = strings.Split(dimensions," ")
	var row = rowcol[0]
	var col = rowcol[1]
	var r,_ = strconv.Atoi(row)
	var c,_ = strconv.Atoi(col)
	cell := make([][]Cell, r)
	for i := range cell {
		cell[i] = make([]Cell, c)
	}

	for i:= 0;i <= r-1; i++{
		for j:= 0; j <= c-1;j++{
			cell[i][j].kind=string(lines[i+1][j])
		}
	}
	if scanner.Err() != nil {
		fmt.Println("Error: there was a problem reading the file")
		os.Exit(3)
	}

	file.Close()

	return cell 
}


// drawField draws a representation of the field on a canvas and save the
// canvas to a PNG file with a name given by the parameter filename.  Each cell
// in the field is a 5-by-5 square, and cells of the "D" kind is
// drawn red and cells of the "C" kind is drawn blue.
func drawField(field [][]Cell, filename string) {
	pic := CreateNewCanvas(500,500)
	for i:= 0;i <= len(field)-1; i++{
		for j:= 0; j <= len(field[0])-1;j++{
			if field[i][j].kind  == "C" {
				x1, y1 := float64(i) * 5, float64(j) * 5
				x2, y2 := float64(i+1) * 5, float64(j+1) * 5
				pic.SetStrokeColor(MakeColor(0, 0, 255))
				pic.MoveTo(x1, y1)
				pic.LineTo(x1, y2)
				pic.LineTo(x2, y2)
				pic.LineTo(x2, y1)
				pic.LineTo(x1, y1)
				pic.FillStroke()
			}
			if field[i][j].kind == "D" {
				x1, y1 := float64(i) * 5, float64(j) * 5
				x2, y2 := float64(i+1) * 5, float64(j+1) * 5
				pic.SetStrokeColor(MakeColor(255, 0, 0))
				pic.MoveTo(x1, y1)
				pic.LineTo(x1, y2)
				pic.LineTo(x2, y2)
				pic.LineTo(x2, y1)
				pic.LineTo(x1, y1)
				pic.FillStroke()

			}
		}
	}

	pic.SaveToPNG("Prisoners.png")
}


/*===============================================================
* Functions to simulate the spatial games
*==============================================================*/

// play a game between a cell of type "me" and a cell of type "them" (both me
// and them should be either "C" or "D"). This returns the reward that "me"
// gets when playing against them.
func gameBetween(me, them string, b float64) float64 {
	if me == "C" && them == "C" {
		return 1
	} else if me == "C" && them == "D" {
		return 0
	} else if me == "D" && them == "C" {
		return b
	} else if me == "D" && them == "D" {
		return 0
	} else {
		fmt.Println("type ==", me, them)
		panic("This shouldn't happen")
	}
}


// updateScores goes through every cell, and plays the Prisoner's dilema game
// with each of it's in-field nieghbors (including itself). It updates the
// score of each cell to be the sum of that cell's winnings from the game.
func updateScores(field [][]Cell, b float64) {
	for i:= 0;i <= len(field)-1; i++{
		for j:= 0; j <= len(field[0])-1;j++{
			var me=field[i][j]
			me.score=0
			for k:= i-1;k <=i+1 ; k++{
				for l:=j-1; l <= j+1 ;l++{
					if k==-1||k==len(field)||l==-1||l==len(field[0]){
						me.score +=0 }else{
							var them=field[k][l].kind
							me.score += gameBetween(me.kind,them,b)
						}
					}
				}
				field[i][j]=me
			}
		}
	}


// updateStrategies create a new field by going through every cell (r,c), and
// looking at each of the cells in its neighborhood (including itself) and the
// setting the kind of cell (r,c) in the new field to be the kind of the
// neighbor with the largest score
func updateStrategies(field [][]Cell) [][]Cell {
	for i:= 0;i <= len(field)-1; i++{
		for j:= 0; j <= len(field[0])-1;j++{
			max_i:=i
			max_j:=j
			score:=field[i][j].score
			for k:=-1; k<=1; k++{
				for l:=-1; l<=1; l++{
					in_i:= i+k
					in_j:= j+l
					if in_i>=0 && in_i<len(field) && in_j>=0 && in_j<len(field[0]){
						if score<field[in_i][in_j].score{
							max_i=in_i
							max_j=in_j
							score=field[in_i][in_j].score
							}
						}
					}
				}
				field[i][j].newKind=field[max_i][max_j].kind
			}
		}
		for i:= 0;i <= len(field)-1; i++{
			for j:= 0; j <= len(field[0])-1;j++{
				field[i][j].kind=field[i][j].newKind
			}
		}
		return field
	}


// evolve takes an intial field and evolves it for nsteps according to the game
// rule. At each step, it should call "updateScores()" and the updateStrategies
func evolve(field [][]Cell, nsteps int, b float64) [][]Cell {
	for i := 1; i < nsteps; i++ {
		updateScores(field, b)
		field = updateStrategies(field)
		}
		return field
	}



// Implements a Spatial Games version of prisoner's dilemma. The command-line
// usage is:
//     ./spatial field_file b nsteps
// where 'field_file' is the file continaing the initial arrangment of cells, b
// is the reward for defecting against a cooperator, and nsteps is the number
// of rounds to update stategies.
func main() {
	// parse the command line
	if len(os.Args) != 4 {
		fmt.Println("Error: should spatial field_file b nsteps")
		return
	}
	fieldFile := os.Args[1]
	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil || b <= 0 {
		fmt.Println("Error: bad b parameter.")
		return
		}
		nsteps, err := strconv.Atoi(os.Args[3])
		if err != nil || nsteps < 0 {
			fmt.Println("Error: bad number of steps.")
			return
		}
		field := readFieldFromFile(fieldFile)
		fmt.Println("Field dimensions are:", len(field), "by", len(field[0]))
		// evolve the field for nsteps and write it as a PNG
		evolve(field, nsteps, b)
		drawField(field, "Prisoners.png")
	}
