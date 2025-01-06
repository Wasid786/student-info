package storage

import "github.com/Wasid786/student-info/internal/types"

type Storage interface {
	CreateStudent(name string, age int, email string) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
}
