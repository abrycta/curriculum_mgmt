package main

type Course struct {
	Year   int
	Term   int
	Id     string
	Title  string
	Units  int
	Grade  int
	Degree string
}

func NewCourse(year int, term int, id string, title string, units int, grade int, degree string) *Course {
	return &Course{Year: year, Term: term, Id: id, Title: title, Units: units, Grade: grade, Degree: degree}
}
