package postgres

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)


// Query For AdminInsert
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

// Query For AdminList with searchTerm
const listAdminQuery = `SELECT * FROM admin WHERE deleted_at IS NULL AND (username ILIKE '%%' || $1 || '%%' or email ILIKE '%%' || $1 || '%%' or first_name ILIKE '%%' || $1 || '%%' or last_name ILIKE '%%' || $1 || '%%')`

func (p PostgresStorage) ListAdmin(uf storage.AdminFilter) ([]storage.LoginAdmin, error) {

	var admin []storage.LoginAdmin
	if err := p.DB.Select(&admin, listAdminQuery, uf.SearchTerm); err != nil {
		log.Println(err)
		return nil, err
	}
	return admin, nil
}


// Query For Admin Update
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


// Query For Admin Delete
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

// Query For Get Admin by Username
const getAdminByUsernameQuery = `SELECT * FROM admin WHERE username=$1`

func (s PostgresStorage) GetAdminByUsername(username string) (*storage.LoginAdmin, error) {
	var u storage.LoginAdmin
	if err := s.DB.Get(&u, getAdminByUsernameQuery, username); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

// Query For Check Admin Already Exists
func (p PostgresStorage) CheckAdminExists(username, email string) (bool, error) {
	var exists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM admin WHERE username = $1 or email = $2 AND deleted_at IS NULL)`, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Query For Get Admin By ID
const getAdminByIDQuery = `SELECT * FROM admin WHERE id=$1 AND deleted_at IS NULL`

func (p PostgresStorage) GetAdminByID(id int) (*storage.LoginAdmin, error) {
	var s storage.LoginAdmin
	if err := p.DB.Get(&s, getAdminByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}