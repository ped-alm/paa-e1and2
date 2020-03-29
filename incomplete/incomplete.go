package incomplete

import (
	"fmt"
	"io"
	"os"
	"sort"

	data "github.com/ped-alm/paa-e1and2/student"
	"github.com/ped-alm/paa-e1and2/utils"
)

// This method still incomplete
func sortFileWithExternalMemory() {
	const chunkSize = 1500
	const numFiles = 2

	// opening student file
	dataFile, err := os.Open("files/student")
	utils.CheckErr(err)
	defer dataFile.Close()

	var files []*os.File
	for i := 0; i < numFiles*2; i++ {
		// creating the file to be used on the sort process
		path := fmt.Sprintf("files/file%d", i)
		file, err := os.Create(path)
		utils.CheckErr(err)
		files = append(files, file)
	}

	eof := false
	counter := 0

	for !eof {
		// reads a chunk of student from student file into memory
		var dataArr []data.Student
		for i := 0; i < chunkSize; i++ {
			var d data.Student
			// will stop on EOF
			if d, eof = data.ReadStudent(dataFile); eof {
				break
			}
			dataArr = append(dataArr, d)
		}

		// sorting student on memory
		sort.Slice(dataArr, func(i, j int) bool {
			return dataArr[i].Key < dataArr[j].Key
		})

		// alternates the write between the auxiliary files
		var file io.Writer
		file = files[counter%numFiles]

		// Writes the ordered student chunk into the selected file
		for _, data := range dataArr {
			utils.Write(file, data.Key)
			utils.Write(file, data.Name)
			utils.Write(file, data.Grade)
		}
		counter++
	}

	counter = 0
	var stops []bool
	var readCountArr []int
	var dataArr []data.Student

	for i := 0; i < numFiles; i++ {
		stops = append(stops, false)
		readCountArr = append(readCountArr, 0)
		dataArr = append(dataArr, data.Student{})
	}

	if dataArr[0], eof = data.ReadStudent(files[0]); eof {
		//TODO treat when eof
	}
	if dataArr[1], eof = data.ReadStudent(files[1]); eof {
		//TODO treat when eof
	}

	// sets the file that will be written the chunk
	writeFile := files[counter%numFiles+numFiles]
	// will iterate while the current chunks still unmerged
	// and unsorted on the set file
	for {
		// all chunk student have been read
		if stops[0] && stops[1] {
			break
		}

		// chunk student from first file have ended
		// writes remaining chunk student from second file
		if stops[0] {
			data.WriteStudent(writeFile, dataArr[1])
			if dataArr[1], eof = data.ReadStudent(files[1]); eof {
				// end of second file chunk student
				stops[1] = true
			}
			readCountArr[1]++
		}

		// chunk student from second file have ended
		// writes remaining chunk student from first file
		if stops[1] {
			data.WriteStudent(writeFile, dataArr[0])
			if dataArr[0], eof = data.ReadStudent(files[0]); eof {
				// end of first file chunk student
				stops[0] = true
			}
			readCountArr[0]++
		}

		// both files have chunk student
		// sorts the student and writes
		if !stops[0] && !stops[1] {
			if dataArr[0].Key < dataArr[1].Key {
				data.WriteStudent(writeFile, dataArr[0])
				if dataArr[0], eof = data.ReadStudent(files[0]); eof {
					// end of first file
					stops[0] = true
				}
				readCountArr[0]++
			} else {
				data.WriteStudent(writeFile, dataArr[1])
				if dataArr[1], eof = data.ReadStudent(files[1]); eof {
					// end of second file
					stops[1] = true
				}
				readCountArr[1]++
			}
		}

		// end of first file chunk student
		if readCountArr[0] == chunkSize {
			stops[0] = true
		}
		// end of second file chunk student
		if readCountArr[1] == chunkSize {
			stops[1] = true
		}
	}
	counter++

	// deletes all files used on sort process
	for _, f := range files {
		os.Remove(f.Name())
	}
}
