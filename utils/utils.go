package utils

import (
	"encoding/binary"
	"io"
	"os"
)

// kills the execution on error
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// writes the given data into the given file
func Write(w io.Writer, data interface{}) {
	err := binary.Write(w, binary.LittleEndian, data)
	CheckErr(err)
}

// reads from the given file into the specified data pointer
// returns true when EOF is found
func Read(r io.Reader, data interface{}) bool {
	err := binary.Read(r, binary.LittleEndian, data)
	switch err {
	case io.EOF:
		return true
	case nil:
		return false
	default:
		CheckErr(err)
	}
	return false
}

// checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
