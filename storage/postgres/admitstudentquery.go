package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Query For AdmitSTudent Insert
const insertAdmitStudentQuery = `
	INSERT INTO admitstudent(
		class_id,
		username,
		first_name,
		last_name,
		email,
		roll,
		password
		)  
	VALUES(
		:class_id,
		:username,
		:first_name,
		:last_name,
		:email,
		:roll,
		:password
		)RETURNING *;
	`

func (p PostgresStorage) AdmitStudentCreate(s storage.AdmitStudents) (*storage.AdmitStudents, error) {

	stmt, err := p.DB.PrepareNamed(insertAdmitStudentQuery)
	if err != nil {
		log.Fatalln(err)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	s.Password = string(hashPass)
	if err := stmt.Get(&s, s); err != nil {
		return nil, err
	}

	if s.ID == 0 {
		return nil, fmt.Errorf("unable to insert user into db")
	}

	return &s, nil
}


// Query For AdmitStudent by INNER JOIN with search term
const listAdmitStudentQuery = `SELECT admitstudent.id, admitstudent.class_id, username, first_name, last_name, email, roll, status, classes.classname
FROM admitstudent
INNER JOIN classes ON admitstudent.class_id = classes.id
WHERE classes.deleted_at IS NULL AND admitstudent.deleted_at IS NULL AND (username ILIKE '%%' || $1 || '%%' or first_name ILIKE '%%' || $1 || '%%' or last_name ILIKE '%%' || $1 || '%%' or email ILIKE '%%' || $1 || '%%' or classname ILIKE '%%' || $1 || '%%')
ORDER BY classes.classname ASC;`

func (p PostgresStorage) AdmitStudentList(uf storage.AdmitStudentFilter) ([]storage.AdmitStudents, error) {

	var admitstudent []storage.AdmitStudents
	if err := p.DB.Select(&admitstudent, listAdmitStudentQuery, uf.SearchTerm); err != nil {
		log.Println(err)
		return nil, err
	}
	return admitstudent, nil
}

// Query For AdmitStudent Update
const updateAdmitStudentQuery = `UPDATE admitstudent
		SET class_id = :class_id, 
		username = :username,
		first_name = :first_name,
		last_name = :last_name,
		email = :email,
		roll = :roll,
		password =:password,
		status = :status
		WHERE id = :id
		RETURNING *;
	`

func (p PostgresStorage) UpdateAdmitStudent(s storage.AdmitStudents) (*storage.AdmitStudents, error) {
	stmt, err := p.DB.PrepareNamed(updateAdmitStudentQuery)
	if err != nil {
		log.Fatalln(err)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	s.Password = string(hashPass)
	if err := stmt.Get(&s, s); err != nil {
		return nil, err
	}

	if s.ID == 0 {
		return nil, fmt.Errorf("unable to insert user into db")
	}

	return &s, nil
}


// Query For AdmitStudent Delete
const deleteAdmitStudentByIDQuery = `UPDATE admitstudent SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) DeleteAdmitStudentByID(id int) error {
	res, err := p.DB.Exec(deleteAdmitStudentByIDQuery, id)
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

// Query For Get AdmitStudent By ID
const getAdmitStudentByIDQuery = `SELECT * FROM admitstudent WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetAdmitStudentByID(id int) (*storage.AdmitStudents, error) {
	var s storage.AdmitStudents
	if err := p.DB.Get(&s, getAdmitStudentByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}


// Query For AdmitStudent Username Already Exists
func (p PostgresStorage) CheckAdmitStudentUsernameExists(username string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admitstudent WHERE username = $1 AND deleted_at IS NULL)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}


// Query For AdmitStudent Email Already Exists
func (p PostgresStorage) CheckAdmitStudenEmailtExists(email string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admitstudent WHERE email = $1 AND deleted_at IS NULL)`, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}


// Query For AdmitStudent Roll Already Exists
func (p PostgresStorage) CheckAdmitStudentRollExists(roll int) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admitstudent WHERE roll = $1 AND deleted_at IS NULL AND deleted_at IS NULL)`, roll).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}


const getFixedStudentSubjectByIDQuery = `SELECT * FROM student_subject WHERE student_id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetFixedStudentSubjectByID(id int) ([]storage.StudentSubject, error) {
	var s []storage.StudentSubject
	if err := p.DB.Select(&s, getFixedStudentSubjectByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return s, nil
}

const getFixedSubjectByIDQuery = `SELECT subjectname FROM subjects WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetFixedSubjectByID(id int) (storage.Subjects, error) {
	var s storage.Subjects
	if err := p.DB.Get(&s, getFixedSubjectByIDQuery, id); err != nil {
		log.Println(err)
		return storage.Subjects{}, err
	}

	return s, nil
}