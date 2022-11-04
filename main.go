package main

import (
	"fmt"
	"strconv"
)

func main() {
	// file := Open("res/bscs.csv")
	// courses := Parse(file)
	connect()
	for {
		menu()
		choice, _ := strconv.Atoi(input("Choice"))
		switch choice {
		case 1:
			curriculum(1, 1)
		case 2:
			gradedCourses()
		case 3:
			grade()
		case 4:
			edit()
		case 5:
			elective()
		case 6:
			shift()
		default:
			fmt.Println("Invalid input.")
			return
		}

	}

}
