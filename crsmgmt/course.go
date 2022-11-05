// Allows for the manipulation of courses via a management interface.
package crsmgmt

// A representation of the Course data structure
type Course struct {
	Year   int
	Term   int
	Id     string
	Title  string
	Units  int
	Grade  int
	Degree string
}

// Factory method for creating course objects
func NewCourse(year int, term int, id string, title string, units int, grade int, degree string) *Course {
	return &Course{Year: year, Term: term, Id: id, Title: title, Units: units, Grade: grade, Degree: degree}
}

// Generate a constructor here ^
// struct == Go's Data structure
// struct != classes
// NewCourse is a factory method
// Returns a pointer the newly created course.
// Builder pattern
// A factory is a builder that returns objects of some type
// pre-initialized to some state
