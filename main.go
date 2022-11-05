package main

import (
	"curriculum_mgmt/crsmgmt"
	"fmt"
	"strconv"
)

func main() {
	// file := Open("res/bscs.csv")
	// courses := Parse(file)
	crsmgmt.Connect()
	for {
		crsmgmt.Menu()
		choice, _ := strconv.Atoi(crsmgmt.Input("Choice"))
		switch choice {
		case 1:
			crsmgmt.Curriculum(1, 1)
		case 2:
			crsmgmt.GradedCourses()
		case 3:
			crsmgmt.Grade()
		case 4:
			crsmgmt.Edit()
		case 5:
			crsmgmt.Elective()
		case 6:
			crsmgmt.Shift()
		default:
			fmt.Println("Invalid input.")
			return
		}

	}

}
