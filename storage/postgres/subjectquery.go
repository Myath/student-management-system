package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
)


// Query For SubjectInsert
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

// Query For SubjectList by INNER JOIN with search term
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


// Query For SUbjectUpdate
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

// Query For SubjectDelete
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

// Query For Get Subject By ID
const getSubjectByIDQuery = `SELECT * FROM subjects WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetSubjectByID(id int) (*storage.Subjects, error) {
	var s storage.Subjects
	if err := p.DB.Get(&s, getSubjectByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

// Query For Subject ALready Exists Check
func (p PostgresStorage) CheckSubjectExists(subjectname string, class_id int) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM subjects WHERE subjectname = $1 AND class_id =$2 AND deleted_at IS NULL)`, subjectname, class_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}