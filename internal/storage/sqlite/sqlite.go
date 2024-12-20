package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ananddub/students-api/internal/config"
	"github.com/ananddub/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	result, err := s.Db.Exec("INSERT INTO students (name, email, age) VALUES (?, ?, ?)", name, email, age)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS
        students (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            email TEXT,
            age INTEGER)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? limit 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with id %d not found", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("failed to get student with id %d: %v", id, err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	rows, err := s.Db.Query("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
