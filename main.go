package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	data "github.com/ped-alm/paa-e1and2/student"
	"github.com/ped-alm/paa-e1and2/utils"
)

func main() {
	// doesn't process files again if they already exists
	if !utils.FileExists("files/data") {
		readAndMarshalData()
	}
	if !utils.FileExists("files/index") {
		indexDataFile()
	}

	// opening data file
	dataFile, err := os.Open("files/data")
	utils.CheckErr(err)
	defer dataFile.Close()

	// opening index file
	indexFile, err := os.Open("files/index")
	utils.CheckErr(err)
	defer indexFile.Close()

	// user menu
	// not treating invalid insertions
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Choose an option to search by")
		fmt.Println("1 - Position")
		fmt.Println("2 - Key")
		fmt.Println("3 - Exit")
		fmt.Print(": ")

		choice, _, err := reader.ReadRune()
		utils.CheckErr(err)

		if choice == '3' {
			break
		}

		fmt.Print("\nInsert Value: ")

		// to fix precocious read
		reader.ReadRune()
		valueText, err := reader.ReadString('\n')
		utils.CheckErr(err)
		// to Work on Windows
		valueText = strings.Replace(valueText, "\r\n", "", -1)
		// to Work on Linux
		valueText = strings.Replace(valueText, "\n", "", -1)

		value, err := strconv.Atoi(valueText)
		utils.CheckErr(err)

		// aux vars
		var student data.Student
		var eof bool

		switch choice {
		case '1':
			student, eof = data.SeekStudent(dataFile, value)
		case '2':
			student, eof = data.FindStudent(dataFile, indexFile, value)
		}

		if eof {
			fmt.Println("position not found")
		} else {
			fmt.Printf("\nKey: %d \n", student.Key)
			fmt.Printf("Name: %s \n", student.Name)
			fmt.Printf("Grade: %f \n\n", student.Grade)
		}

	}
}

// opens the read file and writes its students into a binary file
func readAndMarshalData() {
	// opening read file
	readFile, err := os.Open("files/read")
	utils.CheckErr(err)
	defer readFile.Close()

	// opening data file
	dataFile, err := os.Create("files/data")
	utils.CheckErr(err)
	defer dataFile.Close()

	// marshal to data file process
	scanner := bufio.NewScanner(readFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ";")

		key, err := strconv.Atoi(fields[0])
		utils.CheckErr(err)
		grade, err := strconv.ParseFloat(strings.ReplaceAll(fields[2], ",", "."), 32)
		utils.CheckErr(err)
		student, err := data.NewStudent(key, grade, fields[1])
		utils.CheckErr(err)

		data.WriteStudent(dataFile, student)
	}
}

// Opens the data file and creates an index for it.
// The index file is compound by the student key and a counter
func indexDataFile() {
	// opening index file
	indexFile, err := os.Create("files/index")
	utils.CheckErr(err)
	defer indexFile.Close()

	// opening data file
	dataFile, err := os.Open("files/data")
	utils.CheckErr(err)
	defer dataFile.Close()

	// aux variables
	var student data.Student
	eof := false
	var counter int32 = 0

	for {
		// will stop on EOF
		if student, eof = data.ReadStudent(dataFile); eof {
			break
		}

		utils.Write(indexFile, student.Key)
		utils.Write(indexFile, counter)
		counter++
	}
}
