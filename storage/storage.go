package storage

import (
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AdminFilter struct {
	SearchTerm string
}

type ClassFilter struct {
	SearchTerm string
}

type SubjectFilter struct {
	SearchTerm string
}

type AdmitStudentFilter struct {
	SearchTerm string
}


type AdminLogin struct {
	ID         int              `db:"id" form:"-"`
	Username   string           `db:"username" form:"username"`
	Password   string           `db:"password" form:"password"`
	FormError  map[string]error `db:"-"`
}

type LoginAdmin struct {
	ID         int              `db:"id" form:"-"`
	Username   string           `db:"username" form:"username"`
	Email      string           `db:"email" form:"email"`
	First_name string           `db:"first_name" form:"first_name"`
	Last_name  string           `db:"last_name" form:"last_name"`
	Password   string           `db:"password" form:"password"`
	Status     bool             `db:"status" form:"status"`
	CreatedAt  time.Time        `db:"created_at" form:"-"`
	UpdatedAt  time.Time        `db:"updated_at" form:"-"`
	DeletedAt  sql.NullTime     `db:"deleted_at" form:"-"`
	CSRFToken  string           `db:"-" form:"csrf_token"`
	FormError  map[string]error `db:"-"`
}

type Classes struct {
	ID         int              `db:"id" form:"-"`
	Class_name string           `db:"classname" form:"classname"`
	CreatedAt  time.Time        `db:"created_at" form:"-"`
	UpdatedAt  time.Time        `db:"updated_at" form:"-"`
	DeletedAt  sql.NullTime     `db:"deleted_at" form:"-"`
	CSRFToken  string           `db:"-" form:"csrf_token"`
	FormError  map[string]error `db:"-"`
}

type Subjects struct {
	ID           int              `db:"id" form:"-"`
	Class_ID     int              `db:"class_id" form:"class_id"`
	Subject_name string           `db:"subjectname" form:"subjectname"`
	Class_name   string           `db:"classname" form:"-"`
	CreatedAt    time.Time        `db:"created_at" form:"-"`
	UpdatedAt    time.Time        `db:"updated_at" form:"-"`
	DeletedAt    sql.NullTime     `db:"deleted_at" form:"-"`
	CSRFToken    string           `db:"-" form:"csrf_token"`
	FormError    map[string]error `db:"-"`
}

type AdmitStudents struct {
	ID         int              `db:"id" form:"-"`
	Class_ID   int              `db:"class_id" form:"class_id"`
	Username   string           `db:"username" form:"username"`
	First_name string           `db:"first_name" form:"first_name"`
	Last_name  string           `db:"last_name" form:"last_name"`
	Email      string           `db:"email" form:"email"`
	Class_name string           `db:"classname" form:"-"`
	Roll       int              `db:"roll" form:"roll"`
	Marks      int              `db:"marks" form:"marks"`
	Password   string           `db:"password" form:"password"`
	Status     bool             `db:"status" form:"status"`
	CreatedAt  time.Time        `db:"created_at" form:"-"`
	UpdatedAt  time.Time        `db:"updated_at" form:"-"`
	DeletedAt  sql.NullTime     `db:"deleted_at" form:"-"`
	FormError  map[string]error `db:"-"`
	CSRFToken  string           `db:"-" form:"csrf_token"`
}

type StudentSubject struct {
	ID           int    `db:"id" form:"-"`
	StudentID    int    `db:"student_id" form:"student_id"`
	SubjectID    int    `db:"subject_id" form:"subject_id"`
	Marks        int    `db:"marks" form:"marks"`
	Subject_name string `db:"subjectname" form:"subjectname"`
	Username     string `db:"username" form:"username"`
	First_name   string `db:"first_name" form:"first_name"`
	Last_name    string `db:"last_name" form:"last_name"`
	Email        string `db:"email" form:"email"`
	Class_name   string `db:"classname" form:"-"`
	Roll         int    `db:"roll" form:"roll"`
	Mark         map[int]int
	CreatedAt    time.Time    `db:"created_at" form:"-"`
	UpdatedAt    time.Time    `db:"updated_at" form:"-"`
	DeletedAt    sql.NullTime `db:"deleted_at" form:"-"`
}

func (s AdmitStudents) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Class_ID,
			validation.Required.Error("Select class field is required."),
		),
		validation.Field(&s.Username,
			validation.Required.Error("The name field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&s.First_name,
			validation.Required.Error("The name field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&s.Last_name,
			validation.Required.Error("The name field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&s.Email,
			validation.Required.Error("The email field is required."),
			is.Email.Error("This email is not valid."),
		),
		validation.Field(&s.Roll,
			validation.Required.Error("Student roll start from 1"),
			validation.Min(1).Error("Student roll start from 1"),
			validation.Max(200).Error("Only 200 Students are allowed"),
		),
		validation.Field(&s.Password,
			validation.Required.When(s.ID == 0).Error("The password field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
	)
}

func (a LoginAdmin) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Username,
			validation.Required.Error("The username field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&a.First_name,
			validation.Required.Error("The name field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&a.Last_name,
			validation.Required.Error("The name field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&a.Email,
			validation.Required.Error("The email field is required."),
			is.Email.Error("This email is not valid."),
		),
		validation.Field(&a.Password,
			validation.Required.When(a.ID == 0).Error("The password field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
	)
}

func (a AdminLogin) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Username,
			validation.Required.Error("The username field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&a.Password,
			validation.Required.Error("The password field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
	)
}

func (c Classes) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Class_name,
			validation.Required.Error("The classname field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
	)
}

func (s Subjects) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Class_ID,
			validation.Required.Error("Select class field is required."),
		),
		validation.Field(&s.Subject_name,
			validation.Required.Error("Subject Name field is required"),
		),
	)
}
