package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
)

const listQuery = `SELECT * FROM students WHERE deleted_at IS NULL AND (name ILIKE '%%' || $1 || '%%' or email ILIKE '%%' || $1 || '%%')`

func (p PostgresStorage) ListStudent(uf storage.StudentFilter) ([]storage.Student, error) {
	var student []storage.Student

	if err := p.DB.Select(&student, listQuery, uf.SearchTerm); err != nil {
		log.Println(err)
		return nil, err
	}
	return student, nil
}

const insertStudentQuery = `
	INSERT INTO students(
		name,
		email,
		roll,
		english,
		bangla,		
		mathematics,
		grade,
		gpa 
		)  
	VALUES(
		:name,
		:email,
		:roll,
		:english,
		:bangla,		
		:mathematics,
	    :grade,
	    :gpa
		)RETURNING *;
	`

func (p PostgresStorage) CreateStudent(s storage.Student) (*storage.Student, error) {

	stmt, err := p.DB.PrepareNamed(insertStudentQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const updateQuery = `UPDATE students
		SET name = :name, 
			email = :email,
			roll = :roll,
			english = :english,
			bangla = :bangla,
			mathematics = :mathematics,
			grade = :grade,
			gpa = :gpa,
			status = :status
		WHERE id = :id
		RETURNING *;
	`

func (p PostgresStorage) UpdateStudent(s storage.Student) (*storage.Student, error) {
	stmt, err := p.DB.PrepareNamed(updateQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}

const getStudentByIDQuery = `SELECT * FROM students WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetStudentByID(id int) (*storage.Student, error) {
	var s storage.Student
	if err := p.DB.Get(&s, getStudentByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const deleteStudentByIDQuery = `UPDATE students SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) DeleteStudentByID(id int) error {
	res, err := p.DB.Exec(deleteStudentByIDQuery, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if rowCount <= 0 {
		return fmt.Errorf("unable to delete user")
	}

	return nil
}


func (p PostgresStorage) CheckStudentExists(email string, roll int) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM students WHERE email = $1 or roll = $2)`, email, roll).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}


