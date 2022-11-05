package crsmgmt

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Old methods for CSV parsing
// This file will be deleted in the demo

// Open
func Open(name string) *os.File {
	csvfile, err := os.Open(name)
	if err != nil {
		fmt.Println("Error")
	}
	return csvfile
}

// Parse
func Parse(csvFile *os.File) []Course {
	var courses []Course
	reader := csv.NewReader(csvFile)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		year, _ := strconv.Atoi(line[0])
		term, _ := strconv.Atoi(line[1])
		id := line[2]
		title := line[3]
		units, _ := strconv.Atoi(line[4])

		if len(line) == 6 {
			grade, _ := strconv.Atoi(line[5])
			courses = append(courses,
				*newCourse(year, term, id, title, units, grade, "BSCS"))
		}

		courses = append(courses,
			*newCourse(year, term, id, title, units, 0, "BSCS"))
	}
	return courses
}
