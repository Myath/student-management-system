package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) StudentStore(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	student := storage.Student{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		h.pharseCreateStudent(w, StudentForm{
			Student:   student,
			FormError: student.FormError,
			CSRFToken: nosurf.Token(r),
		})
		return
	}

	checkAlreadyExist, err := h.IsStudent(w, r, student.Email, student.Roll)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkAlreadyExist {
		h.pharseCreateStudent(w, StudentForm{
			Student:   student,
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"AlreadyExistCheck": fmt.Errorf("User Already Exists"),
			}})
		return
	}

	grade, gpa := Grade(student.English, student.Bangla, student.Mathematics)

	student.Grade = grade
	student.GPA = gpa

	_, eRr := h.storage.CreateStudent(student)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)
}

func (h Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	h.pharseCreateStudent(w, StudentForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) IsStudent(w http.ResponseWriter, r *http.Request, email string, roll int) (bool, error) {
	ad, err := h.storage.CheckStudentExists(email, roll)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}
