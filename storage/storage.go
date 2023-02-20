package storage

import (
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type StudentFilter struct {
	SearchTerm string
}

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

type Student struct {
	ID          int              `db:"id" form:"-"`
	Name        string           `db:"name" form:"name"`
	Email       string           `db:"email" form:"email"`
	Roll        int              `db:"roll" form:"roll"`
	English     int              `db:"english" form:"eng"`
	Bangla      int              `db:"bangla" form:"ban"`
	Mathematics int              `db:"mathematics" form:"math"`
	Grade       string           `db:"grade" form:"-"`
	GPA         float64          `db:"gpa" form:"-"`
	Status      bool             `db:"status" form:"status"`
	CreatedAt   time.Time        `db:"created_at" form:"-"`
	UpdatedAt   time.Time        `db:"updated_at" form:"-"`
	DeletedAt   sql.NullTime     `db:"deleted_at" form:"-"`
	FormError   map[string]error `db:"-"`
	CSRFToken   string           `db:"-" form:"csrf_token"`
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
	Password   string           `db:"password" form:"password"`
	Status     bool             `db:"status" form:"status"`
	CreatedAt  time.Time        `db:"created_at" form:"-"`
	UpdatedAt  time.Time        `db:"updated_at" form:"-"`
	DeletedAt  sql.NullTime     `db:"deleted_at" form:"-"`
	FormError  map[string]error `db:"-"`
	CSRFToken  string           `db:"-" form:"csrf_token"`
}

type StudentSubject struct {
	ID        int          `db:"id" form:"-"`
	StudentID int       `db:"student_id" form:"student_id"`
	SubjectID int          `db:"subject_id" form:"subject_id"`
	Marks     int          `db:"marks" form:"marks"`
	CreatedAt time.Time    `db:"created_at" form:"-"`
	UpdatedAt time.Time    `db:"updated_at" form:"-"`
	DeletedAt sql.NullTime `db:"deleted_at" form:"-"`
}

func (s Student) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error("The name field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&s.Email,
			validation.Required.When(s.ID == 0).Error("The email field is required."),
			is.Email.Error("This email is not valid."),
		),
		validation.Field(&s.Roll,
			validation.Required.When(s.ID == 0).Error("Student roll start from 1"),
			validation.Min(1).Error("Student roll start from 1"),
			validation.Max(200).Error("Only 200 Students are allowed"),
		),
		validation.Field(&s.English,
			validation.Min(0).Error("The lowest mark is 0."),
			validation.Max(100).Error("The highest mark is 100."),
		),
		validation.Field(&s.Bangla,
			validation.Min(0).Error("The lowest mark is 0."),
			validation.Max(100).Error("The highest mark is 100."),
		),
		validation.Field(&s.Mathematics,
			validation.Min(0).Error("The lowest mark is 0."),
			validation.Max(100).Error("The highest mark is 100."),
		),
	)
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
			validation.Required.When(s.ID == 0).Error("The email field is required."),
			is.Email.Error("This email is not valid."),
		),
		validation.Field(&s.Roll,
			validation.Required.When(s.ID == 0).Error("Student roll start from 1"),
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
		// validation.Field(&a.Email,
		// 	validation.Required.Error("The email field is required."),
		// 	is.Email.Error("This email is not valid."),
		// ),
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

// func rollAlreadyExists(value any) error {
// 	roll, ok := value.(int)
// 	if !ok {
// 		return errors.New("unsupported data given")
// 	}

// 	db, err := sqlx.Connect("postgres", "user=postgres password=secret dbname=students_management sslmode=disable")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	var student []Student

// 	if err := db.Select(&student, `SELECT * FROM students`); err != nil {
// 		log.Fatal(err)
// 	}

// 	// var editUser User
// 	for _, user := range student {
// 		if user.Roll == roll {
// 			return errors.New("the roll already exists")
// 		}
// 	}
// 	return nil
// }

// func emailAlreadyExists(value any) error {
// 	email, ok := value.(string)
// 	if !ok {
// 		return errors.New("unsupported data given")
// 	}

// 	db, err := sqlx.Connect("postgres", "user=postgres password=secret dbname=students_management sslmode=disable")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	var student []Student

// 	if err := db.Select(&student, `SELECT * FROM students`); err != nil {
// 		log.Fatal(err)
// 	}

// 	// var editUser User
// 	for _, user := range student {
// 		if user.Email == email {
// 			return errors.New("the email already exists")
// 		}
// 	}
// 	return nil
// }

// func (s Student) storeAlreadyExistValidate() error {
// 	return validation.ValidateStruct(&s,
// 		validation.Field(&s.Email,
// 			validation.By(emailAlreadyExists),
// 		),
// 		validation.Field(&s.Roll,
// 			validation.By(rollAlreadyExists),
// 		),
// 	)
// }

// type AdmitStudentFilter struct {
// 	SearchTerm string
// }
