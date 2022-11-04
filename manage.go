package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// show subjects for each school term
func show(courses *[]Course) {
	term, year := 1, 1
	header(term, year)
	for _, course := range *courses {
		if course.Term == term && course.Year == year {
			prettify(course)
		} else {
			if term < 3 {
				term++
			} else if term == 3 {
				term = 1
				year++
			}
			// prompt()
			header(year, term)
		}
	}
}

// show courses subroutine for formatting
func prettify(course Course) {
	// fmt.Println(course.Id, '\t', course.Title, '\t', course.Units, '\t', course.Grade)
	if course.Grade == 0 {
		fmt.Printf("%-10.10s %-35.35s %10d %10s\n", course.Id, course.Title, course.Units, "-")
	} else {
		fmt.Printf("%-10.10s %-35.35s %10d %10d\n", course.Id, course.Title, course.Units, course.Grade)
	}
}

func header(year, term int) {
	var yearStr, termStr string
	switch year {
	case 1:
		yearStr = "First Year"
	case 2:
		yearStr = "Second Year"
	case 3:
		yearStr = "Third Year"
	case 4:
		yearStr = "Fourth Year"
	}

	switch term {
	case 1:
		termStr = "First Semester"
	case 2:
		termStr = "Second Semester"
	case 3:
		termStr = "Short Term"

	}
	fmt.Println()
	fmt.Println("======================================================================")
	fmt.Printf("%s, %s\n", yearStr, termStr)
	fmt.Printf("%-10s %-35s %11s %10s\n", "Course ID", "Descriptive Title", "Units", "Grade")
	fmt.Println("======================================================================")
}

func menu() {
	fmt.Println("\n1. My curriculum checklist")
	fmt.Println("2. Show courses with grades")
	fmt.Println("3. Input / Edit grades")
	fmt.Println("4. Edit a course")
	fmt.Println("5. Manage electives")
	fmt.Println("6. Shift from a program, to the BSCS curriculum\n")
}

// possible demo for debugging
// input from bufreader returns
// a string with \r\n at the end
func input(prompt string) string {
	fmt.Print(prompt, ": ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input = strings.Replace(input, "\r\n", "", -1)
	return input
}
