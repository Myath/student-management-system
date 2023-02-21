package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) AdmitStudentEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	editStudents, err := h.storage.GetAdmitStudentByID(uID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.pharseEditStudents(w, AdmitStudentForm{
		Classlist: classlist,
		Student:   *editStudents,
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) AdmitStudentUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	student := storage.AdmitStudents{}
	student.ID = uID
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		h.pharseEditStudents(w, AdmitStudentForm{
			Classlist: classlist,
			Student:   student,
			CSRFToken: nosurf.Token(r),
			FormError: student.FormError,
		})
		return
	}

	singlestudent, err := h.storage.GetAdmitStudentByID(student.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if singlestudent.Username != student.Username {
		checkUsernameAlreadyExist, err := h.IsAdmitStudentUsername(w, r, student.Username)
		if err != nil {
			fmt.Println(err)
			return
		}

		if checkUsernameAlreadyExist {
			h.pharseEditStudents(w, AdmitStudentForm{
				Classlist: classlist,
				Student:   student,
				Class:     storage.Classes{},
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Username": fmt.Errorf("The username already Exist."),
				}})
			return
		}

	}

	if singlestudent.Email != student.Email {
		checkEmailAlreadyExist, err := h.IsAdmitStudentEmail(w, r, student.Email)
		if err != nil {
			fmt.Println(err)
			return
		}

		if checkEmailAlreadyExist {
			h.pharseEditStudents(w, AdmitStudentForm{
				Classlist: classlist,
				Student:   student,
				Class:     storage.Classes{},
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Email": fmt.Errorf("The email already Exist."),
				}})
			return
		}

	} 
	
	
	if singlestudent.Roll != student.Roll {
		checkRollAlreadyExist, err := h.IsAdmitStudentRoll(w, r, student.Roll)

		if err != nil {
			fmt.Println(err)
			return
		}

		if checkRollAlreadyExist {
			h.pharseEditStudents(w, AdmitStudentForm{
				Classlist: classlist,
				Student:   student,
				Class:     storage.Classes{},
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Roll": fmt.Errorf("The roll already Exist."),
				}})
			return
		}
	}

	_, eRr := h.storage.UpdateAdminStudent(student)
	if eRr != nil {
		log.Println(eRr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admitstudentlist", http.StatusSeeOther)
}