package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
)

// Query For StudentSubject Insert
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


// Query For StudentSubjectMark Update
const updateMarkQuery = `UPDATE student_subject
		SET student_id = :student_id, 
		subject_id = :subject_id,
		marks = :marks
		WHERE student_id = $1 AND subject_id = $2
		RETURNING *;
	`

func (p PostgresStorage) UpdateMark(s storage.StudentSubject) (*storage.StudentSubject, error) {
	stmt, err := p.DB.PrepareNamed(updateMarkQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}

// Query For GetSubjectByClassID Insert Function StudentSubject
const getSubjectByClassIDQuery = `SELECT * FROM subjects WHERE class_id=$1`

func (s PostgresStorage) GetSubjectByClassID(class_id int) ([]storage.Subjects, error) {

	var u []storage.Subjects
	if err := s.DB.Select(&u, getSubjectByClassIDQuery, class_id); err != nil {
		log.Println(err)
		return nil, err
	}
	return u, nil
}


const getFixedStudentSubjectByIDQuery = `SELECT stsub.*, subjects.subjectname, admitstudent.username, admitstudent.first_name, admitstudent.last_name, admitstudent.email, admitstudent.roll, classes.classname
FROM student_subject AS stsub
INNER JOIN subjects ON stsub.subject_id = subjects.id
INNER JOIN admitstudent ON stsub.student_id = admitstudent.id
INNER JOIN classes ON subjects.class_id = classes.id
WHERE stsub.student_id = $1 AND admitstudent.deleted_at IS NULL`

func (p PostgresStorage) GetFixedStudentSubjectByID(id int) ([]storage.StudentSubject, error) {
	var s []storage.StudentSubject
	if err := p.DB.Select(&s, getFixedStudentSubjectByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return s, nil
}

const deleteMarkByIDQuery = `UPDATE student_subject SET deleted_at = CURRENT_TIMESTAMP WHERE student_id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) DeleteMarkByID(StudentID int) error {
	res, err := p.DB.Exec(deleteMarkByIDQuery, StudentID)
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