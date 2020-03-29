package data

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/ped-alm/paa-e1and2/utils"
)

const nameMaxSize = 25

type Student struct {
	Key   int32
	Grade float32
	Name  [nameMaxSize]byte // chars
}

// Student "constructor"
func NewStudent(key int, grade float64, name string) (Student, error) {
	if len(name) > nameMaxSize {
		return Student{}, errors.New(fmt.Sprintf("Name length greater than the max value (%d)", nameMaxSize))
	}
	// Pads with empty space if Name does not fill all length
	name = fmt.Sprintf("%-"+strconv.Itoa(nameMaxSize)+"s", name)

	std := Student{
		Key:   int32(key),
		Grade: float32(grade),
	}

	// "casting" Name from string to "char"
	copy(std.Name[:], name)
	return std, nil
}

// reads Student from specified file
// returns true when EOF is found
func ReadStudent(file io.Reader) (Student, bool) {
	std := Student{}
	if utils.Read(file, &std.Key) || utils.Read(file, &std.Name) || utils.Read(file, &std.Grade) {
		return std, true
	}
	return std, false
}

// reads Student from specified file and position
// returns true when EOF is found
func SeekStudent(file *os.File, pos int) (Student, bool) {
	file.Seek(int64(pos)*33, 0)

	std := Student{}
	if utils.Read(file, &std.Key) || utils.Read(file, &std.Name) || utils.Read(file, &std.Grade) {
		return std, true
	}
	return std, false
}

// reads Student from specified file using index key
// returns true when EOF is found
func FindStudent(file *os.File, index *os.File, findKey int) (Student, bool) {
	std := Student{}
	var key, counter int32

	for key != int32(findKey) {
		if utils.Read(index, &key) || utils.Read(index, &counter) {
			return std, true
		}
	}

	return SeekStudent(file, int(counter))
}

// writes Student into specified file
//  |   Key   |   Name   |  Grade  |
//  | 4 bytes | 25 bytes | 4 bytes |
func WriteStudent(file io.Writer, std Student) {
	utils.Write(file, std.Key)
	utils.Write(file, std.Name)
	utils.Write(file, std.Grade)
}
