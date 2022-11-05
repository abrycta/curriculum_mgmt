package crsmgmt

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	// _ "github.com/go-sql-driver/mysql"
)

// Stores a reference to the database handle
var db *sql.DB

// Errors returned by functions are assigned to this variable
var err error

// Connects to the MariaDB back-end and stores a reference to the database handle for use with SQL transactions.
func Connect() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "curriculum",
		AllowNativePasswords: true,
	}

	// get a database handle
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}

// returns a list of all courses
// for the demo, show a redundant for loop
// that takes the Input of the parsed csv entries
// for the profiling part
func courses() ([]Course, error) {
	// A course slice to hold data from the return values of the query
	var courses []Course

	rows, err := db.Query("SELECT * FROM course")

	if err != nil {
		return nil, fmt.Errorf("courses: %v", err)
	}

	// wrap via closure
	// closures are anonymous nested functions
	// that have bindings to variables
	// outside its lexical scope
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var course Course
		if err = rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, &course.Grade, &course.Degree); err != nil {
			// alternative to coalce, handle null values
			// Scan returns an error
			// we pass pointers
			rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, 0, &course.Degree)
			// return nil, fmt.Errorf("courses: %v", err)
		}
		courses = append(courses, course)
		// fmt.Println(course)
	}

	return courses, nil
}

// Interface for editing a course object
func Edit() {
	var choice int
	var stmt *sql.Stmt
	// verify if row exists
	key := Input("Search Key")
	if !searchByID(key) {
		fmt.Println("Invalid search parameters.")
		return
	}

	// show menu to edit what particular column
	choice, input := editMenu()
	// id, title, units, back
	switch choice {
	case 1:
		stmt, err = db.Prepare("UPDATE course" +
			" SET id = ? " +
			"WHERE id = ?")
	case 2:
		stmt, err = db.Prepare("UPDATE course" +
			" SET title = ? " +
			"WHERE id = ?")
	case 3:
		stmt, err = db.Prepare("UPDATE course" +
			" SET units = ? " +
			"WHERE id = ?")
	case 4:
		stmt, err = db.Prepare("UPDATE course" +
			" SET degree = ? " +
			" WHERE id = ?")
	default:
		return
	}

	// Execute
	if err != nil {
		fmt.Println(err)
	}

	_, err = stmt.Exec(input, key)

	if err != nil {
		fmt.Println(err)
	}
	// call an input function for input verification, return an error if malformed input is detected
	// use a prepared statement to reflect changes into the server
	// call show()
}

// Validator function to check if a course with a matching ID exists
func searchByID(id string) bool {
	var stmt, err = db.Prepare("SELECT * FROM course WHERE id = ?")
	var banana Course // ???
	if err != nil {
		log.Fatal(err)
	}
	// multiple cursor
	// err = stmt.QueryRow(id).Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, &course.Grade)
	err = stmt.QueryRow(id).Scan(&banana.Year, &banana.Term, &banana.Id, &banana.Title, &banana.Units, &banana.Grade, &banana.Degree)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	header(banana.Year, banana.Term)
	prettify(banana)
	return true
}

// Helper subroutine for the edit() function
func editMenu() (int, string) {
	var val int
	var inputstr string
	fmt.Println(
		"\n" +
			"Edit?\n" +
			"1. Course ID\n" +
			"2. Course Name\n" +
			"3. Units\n" +
			"4. Degree\n" +
			"5. Back\n")
	inputstr = Input("Choice")
	val, _ = strconv.Atoi(inputstr)
	if val != 5 {
		inputstr = Input("Value")
	}
	return val, inputstr
}

// Provides an interface for appending grade properties to course objects.
func Grade() {
	var grade int
	var tmp string
	// verify if row exists
	// can be used as another case for the multiple cursor function
	key := Input("Search Key")
	if !searchByID(key) {
		fmt.Println("Invalid search parameters")
		return
	}
	tmp = Input("Grade")
	grade, _ = strconv.Atoi(tmp)

	stmt, err := db.Prepare("UPDATE course" +
		" SET grade = ? " +
		"WHERE id = ?")

	if err != nil {
		fmt.Println(err)
	}

	_, err = stmt.Exec(grade, key)

	if err != nil {
		fmt.Println(err)
	}
}

// Displays the course checklist
func Curriculum(year, term int) {
	stmt, err := db.Prepare("SELECT * FROM course WHERE year = ? AND term = ?")

	if err != nil {
		fmt.Println(err)
	}

	var courses []Course
	rows, err := stmt.Query(year, term)
	for rows.Next() {
		var course Course
		if err = rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, &course.Grade, &course.Degree); err != nil {
			// alternative to coalce, handle null values
			_ = rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, 0, &course.Degree)
		}
		courses = append(courses, course)
	}

	// display
	header(year, term)
	for _, crs := range courses {
		prettify(crs)
	}

	var tmp string
	tmp = Input("\nContinue? [y]es | custom [r]ange | [n]o")
	if tmp == "r" {
		tmp = Input("Term")
		term, _ := strconv.Atoi(tmp)
		if term == 1 || term == 2 || term == 3 {
			tmp = Input("Year")
			year, _ := strconv.Atoi(tmp)
			if year == 1 || year == 2 || year == 3 || year == 4 {
				Curriculum(year, term)
			}
		} else {
			return
		}
	} else if tmp == "y" {
		// do not increment year yet
		if term < 3 {
			// add additional condition because 4th year
			// does not have a short term
			if term == 2 && year == 4 {
				return
			}
			Curriculum(year, term+1)
		} else {
			// reset term, start of new year
			if year != 4 {
				Curriculum(year+1, 1)
			} else {
				return
			}
		}
	} else {
		return
	}
}

// Interface for managing elective courses.
func Elective() {
	var courses []Course
	var stmt *sql.Stmt
	rows, err := db.Query("SELECT c.year, term, c.id, title, units, grade, d.degree FROM course c JOIN degree d on d.id = c. degree WHERE title LIKE 'Elect%' ORDER BY title")

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var course Course
		if err = rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, &course.Grade, &course.Degree); err != nil {
			// alternative to coalce, handle null values
			rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, 0, &course.Degree)
			// return nil, fmt.Errorf("courses: %v", err)
		}
		courses = append(courses, course)
	}
	// show elective courses
	fmt.Println("Available Elective Courses:")
	for _, crs := range courses {
		prettify(crs)
	}

	choice, course := electiveMenu()

	stmt, err = db.Prepare("UPDATE course SET year = ?, term = ?, id = ?, title = ?, units = ?, grade = ? WHERE id = ?")

	// Execute

	if err != nil {
		fmt.Println(err)
	}

	// Replace this entry
	_, err = stmt.Exec(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, &course.Grade, "CSE "+strconv.Itoa(choice))

	if err != nil {
		fmt.Println(err)
	}

}

// Helper method for the elective() function
func electiveMenu() (int, *Course) {
	var tmp string
	var choice int
	var year int
	var term int
	var id string
	var title string
	var units int
	var grade int
	fmt.Println("\n1. Enroll in an Elective Class\n" +
		"2. Back\n")
	tmp = Input("Choice")
	choice, _ = strconv.Atoi(tmp)

	if choice == 2 {
		return 0, nil
	}

	// select an elective
	tmp = Input("\nSelect an Elective")
	choice, _ = strconv.Atoi(tmp)

	// modify the term
	switch choice {
	case 1:
		year = 4
		term = 1
	case 2:
		year = 4
		term = 1
	case 3:
		year = 4
		term = 2
	case 4:
		year = 4
		term = 2
	}

	id = Input("Course ID")
	title = Input("Descriptive Title")
	tmp = Input("Units")
	units, _ = strconv.Atoi(tmp)
	tmp = Input("Grade")
	grade, _ = strconv.Atoi(tmp)

	return choice, newCourse(year, term, id, title, units, grade, "BSCS")

}

// Returns a list of graded courses
func GradedCourses() {
	var courses []Course

	rows, err := db.Query("SELECT * FROM course WHERE grade IS NOT NULL AND grade > 0")

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var course Course
		if err = rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, &course.Grade, &course.Degree); err != nil {
			// alternative to coalce, handle null values
			rows.Scan(&course.Year, &course.Term, &course.Id, &course.Title, &course.Units, 0, &course.Degree)
			// return nil, fmt.Errorf("courses: %v", err)
		}
		courses = append(courses, course)
		// fmt.Println(course)
	}

	if len(courses) == 0 {
		_ = Input("\nNo graded courses yet.\nPress any key to continue")
	} else {
		for _, crs := range courses {
			prettify(crs)
		}
	}
}

// Allows a user to credit courses taken from their previous degree program
func Shift() {
	var course Course
	var courses []Course
	var dgr int
	degree := Input("From what BS/A Degree?")
	stmt, err := db.Prepare("INSERT INTO degree VALUES(NULL, ?)")

	if err != nil {
		fmt.Println(err)
	}

	_, err = stmt.Exec(degree)

	if err != nil {
		fmt.Println(err)
	}

	// get the integer degree representation
	stmt, err = db.Prepare("SELECT id FROM degree WHERE degree = ?")
	row := stmt.QueryRow(degree)

	if err = row.Scan(&dgr); err != nil {
		fmt.Println(err)
	}

	for {
		course = *enroll()
		courses = append(courses, course)
		if Input("Add more courses [y] | [n]") == "n" {
			break
		}
	}

	for _, crs := range courses {
		// Insert into database
		stmt, err = db.Prepare("INSERT INTO course VALUES(?, ?, ?, ?, ?, ?, ?)")

		if err != nil {
			fmt.Println(err)
		}
		_, err = stmt.Exec(crs.Year, crs.Term, crs.Id, crs.Title, crs.Units, crs.Grade, dgr)

		if err != nil {
			fmt.Println(err)
		}
	}

}

// Helper function for enrolling courses for the shift() function
func enroll() *Course {
	var tmp string
	var year int
	var term int
	var id string
	var title string
	var units int
	var grade int
	var degree string

	fmt.Println("\nCrediting a course from previous degree program: \n")
	tmp = Input("Year")
	year, _ = strconv.Atoi(tmp)
	tmp = Input("Term")
	term, _ = strconv.Atoi(tmp)
	id = Input("Course ID")
	title = Input("Descriptive Title")
	tmp = Input("Units")
	units, _ = strconv.Atoi(tmp)
	tmp = Input("Grade")
	grade, _ = strconv.Atoi(tmp)

	return newCourse(year, term, id, title, units, grade, degree)
}
