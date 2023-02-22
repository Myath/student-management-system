package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
)

const insertClassQuery = `
	INSERT INTO classes(
		classname
		)  
	VALUES(
		:classname
		)RETURNING *;
	`
// Query For ClassInsert
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

// Query For ClassList with searchTerm
const listClassQuery = `SELECT * FROM classes WHERE deleted_at IS NULL AND (classname ILIKE '%%' || $1 || '%%')`

func (p PostgresStorage) ListClass(uf storage.ClassFilter) ([]storage.Classes, error) {

	var class []storage.Classes
	if err := p.DB.Select(&class, listClassQuery, uf.SearchTerm); err != nil {
		log.Println(err)
		return nil, err
	}
	return class, nil
}

// Query For ClassUpdate
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

// Query For ClassDelete
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

// Query For Get CLass
const getClasstByIDQuery = `SELECT * FROM Classes WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetClassByID(id int) (*storage.Classes, error) {
	var s storage.Classes
	if err := p.DB.Get(&s, getClasstByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

// Query For Get ClassName
const listClassNameQuery = `SELECT id, classname FROM classes WHERE deleted_at IS NULL ORDER BY classname ASC;`

func (p PostgresStorage) ListOfClassName() ([]storage.Classes, error) {

	var class []storage.Classes
	if err := p.DB.Select(&class, listClassNameQuery); err != nil {
		log.Println(err)
		return nil, err
	}
	return class, nil
}

// Query For Class ALready Exists Check
func (p PostgresStorage) CheckClassExists(classname string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM classes WHERE classname = $1 AND deleted_at IS NULL)`, classname).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}