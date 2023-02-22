package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
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


// Query For StudentSubject Update
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
