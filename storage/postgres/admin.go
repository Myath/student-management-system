package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const getAdminByUsernameQuery = `SELECT * FROM admin WHERE username=$1`

func (s PostgresStorage) GetAdminByUsername(username string) (*storage.LoginAdmin, error) {
	var u storage.LoginAdmin
	if err := s.DB.Get(&u, getAdminByUsernameQuery, username); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

const insertAdminQuery = `
	INSERT INTO admin(
		username,
		email,
		first_name,
		last_name,
		password
		)  
	VALUES(
		:username,
		:email,
		:first_name,
		:last_name,
		:password
		)RETURNING *;
	`

func (p PostgresStorage) CreateAdmin(s storage.LoginAdmin) (*storage.LoginAdmin, error) {

	stmt, err := p.DB.PrepareNamed(insertAdminQuery)
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

func (p PostgresStorage) CheckAdminExists(username, email string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admin WHERE username = $1 or email = $2 AND deleted_at IS NULL)`, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const listAdminQuery = `SELECT * FROM admin WHERE deleted_at IS NULL AND (username ILIKE '%%' || $1 || '%%' or email ILIKE '%%' || $1 || '%%' or first_name ILIKE '%%' || $1 || '%%' or last_name ILIKE '%%' || $1 || '%%')`

func (p PostgresStorage) ListAdmin(uf storage.AdminFilter) ([]storage.LoginAdmin, error) {

	var admin []storage.LoginAdmin
	if err := p.DB.Select(&admin, listAdminQuery, uf.SearchTerm); err != nil {
		log.Println(err)
		return nil, err
	}
	return admin, nil
}

const insertClassQuery = `
	INSERT INTO classes(
		classname
		)  
	VALUES(
		:classname
		)RETURNING *;
	`

func (p PostgresStorage) CreateClass(s storage.Classes) (*storage.Classes, error) {

	stmt, err := p.DB.PrepareNamed(insertClassQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

func (p PostgresStorage) CheckClassExists(classname string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM classes WHERE classname = $1 AND deleted_at IS NULL)`, classname).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const listClassQuery = `SELECT * FROM classes WHERE deleted_at IS NULL AND (classname ILIKE '%%' || $1 || '%%')`

func (p PostgresStorage) ListClass(uf storage.ClassFilter) ([]storage.Classes, error) {

	var class []storage.Classes
	if err := p.DB.Select(&class, listClassQuery, uf.SearchTerm); err != nil {
		log.Println(err)
		return nil, err
	}
	return class, nil
}

const getClasstByIDQuery = `SELECT * FROM Classes WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetClassByID(id int) (*storage.Classes, error) {
	var s storage.Classes
	if err := p.DB.Get(&s, getClasstByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const deleteClassByIDQuery = `UPDATE classes SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) DeleteSClassByID(id int) error {
	res, err := p.DB.Exec(deleteClassByIDQuery, id)
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

const listClassNameQuery = `SELECT id, classname FROM classes WHERE deleted_at IS NULL ORDER BY classname ASC;`

func (p PostgresStorage) ListOfClassName() ([]storage.Classes, error) {

	var class []storage.Classes
	if err := p.DB.Select(&class, listClassNameQuery); err != nil {
		log.Println(err)
		return nil, err
	}
	return class, nil
}

const insertSubjectQuery = `
	INSERT INTO subjects(
		class_id,
		subjectname
		)  
	VALUES(
		:class_id,
		:subjectname
		)RETURNING *;
	`

func (p PostgresStorage) CreateSubject(s storage.Subjects) (*storage.Subjects, error) {

	stmt, err := p.DB.PrepareNamed(insertSubjectQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

func (p PostgresStorage) CheckSubjectExists(subjectname string, class_id int) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM subjects WHERE subjectname = $1 AND class_id =$2 AND deleted_at IS NULL)`, subjectname, class_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const listSubjectQuery = `SELECT subjects.id, subjects.class_id, subjectname, classes.classname
FROM subjects
INNER JOIN classes ON subjects.class_id = classes.id
WHERE classes.deleted_at IS NULL AND subjects.deleted_at IS NULL AND (subjectname ILIKE '%%' || $1 || '%%' or classname ILIKE '%%' || $1 || '%%')
ORDER BY classes.classname ASC;`

func (p PostgresStorage) SubjectList(uf storage.SubjectFilter) ([]storage.Subjects, error) {

	var subject []storage.Subjects
	if err := p.DB.Select(&subject, listSubjectQuery, uf.SearchTerm); err != nil {
		log.Println(err)
		return nil, err
	}
	return subject, nil
}

const getSubjectByIDQuery = `SELECT * FROM subjects WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetSubjectByID(id int) (*storage.Subjects, error) {
	var s storage.Subjects
	if err := p.DB.Get(&s, getSubjectByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const updateSubjectsQuery = `UPDATE subjects
		SET class_id = :class_id, 
		    subjectname = :subjectname
		WHERE id = :id
		RETURNING *;
	`

func (p PostgresStorage) UpdateSubjects(s storage.Subjects) (*storage.Subjects, error) {
	stmt, err := p.DB.PrepareNamed(updateSubjectsQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}

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

const getSubjectByClassIDQuery = `SELECT * FROM subjects WHERE class_id=$1`

func (s PostgresStorage) GetSubjectByClassID(class_id int) ([]storage.Subjects, error) {

	var u []storage.Subjects
	if err := s.DB.Select(&u, getSubjectByClassIDQuery, class_id); err != nil {
		log.Println(err)
		return nil, err
	}
	return u, nil
}

const insertStudentSubjectMQuery = `
		INSERT INTO student_subject (
			student_id, 
			subject_id, 
			marks
		) VALUES (
			:student_id,
			:subject_id,
			:marks
		)
		RETURNING *;
	`

func (p PostgresStorage) InsertStudentSubject(s storage.StudentSubject) (*storage.StudentSubject, error) {

	stmt, err := p.DB.PrepareNamed(insertStudentSubjectMQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}



const updateStudentSubjectsMarkQuery = `UPDATE student_subject
		SET student_id = :student_id, 
		subject_id = :subject_id,
		marks = :marks
		WHERE id = :id
		RETURNING *;
	`

func (p PostgresStorage) UpdateStudentSubjectsMark(s storage.StudentSubject) (*storage.StudentSubject, error) {
	stmt, err := p.DB.PrepareNamed(updateStudentSubjectsMarkQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}

func (p PostgresStorage) CheckAdmitStudentUsernameExists(username string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admitstudent WHERE username = $1 AND deleted_at IS NULL)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p PostgresStorage) CheckAdmitStudenEmailtExists(email string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admitstudent WHERE email = $1 AND deleted_at IS NULL)`, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p PostgresStorage) CheckAdmitStudentRollExists(roll int) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admitstudent WHERE roll = $1 AND deleted_at IS NULL AND deleted_at IS NULL)`, roll).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

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

const deleteSubjectByIDQuery = `UPDATE subjects SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) DeleteSubjectByID(id int) error {
	res, err := p.DB.Exec(deleteSubjectByIDQuery, id)
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

const getAdmitStudentByIDQuery = `SELECT * FROM admitstudent WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetAdmitStudentByID(id int) (*storage.AdmitStudents, error) {
	var s storage.AdmitStudents
	if err := p.DB.Get(&s, getAdmitStudentByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const updateAdminStudentQuery = `UPDATE admitstudent
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

func (p PostgresStorage) UpdateAdminStudent(s storage.AdmitStudents) (*storage.AdmitStudents, error) {
	stmt, err := p.DB.PrepareNamed(updateAdminStudentQuery)
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

const updateClassesQuery = `UPDATE classes
		SET classname = :classname
		WHERE id = :id
		RETURNING *;
	`

func (p PostgresStorage) UpdateClasses(s storage.Classes) (*storage.Classes, error) {
	stmt, err := p.DB.PrepareNamed(updateClassesQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}

const getAdminByIDQuery = `SELECT * FROM admin WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetAdminByID(id int) (*storage.LoginAdmin, error) {
	var s storage.LoginAdmin
	if err := p.DB.Get(&s, getAdminByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const updateAdminQuery = `UPDATE admin
		SET username = :username,
		email = :email,
		first_name = :first_name,
		last_name = :last_name,
		password =:password,
		status = :status
		WHERE id = :id
		RETURNING *;
	`

func (p PostgresStorage) UpdateAdmin(s storage.LoginAdmin) (*storage.LoginAdmin, error) {
	stmt, err := p.DB.PrepareNamed(updateAdminQuery)
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


const deleteAdminByIDQuery = `UPDATE admin SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) DeleteAdminByID(id int) error {
	res, err := p.DB.Exec(deleteAdminByIDQuery, id)
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
