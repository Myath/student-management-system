package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) StudentEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	editStudent, err := h.storage.GetStudentByID(uID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form := StudentForm{}
	form.Student = *editStudent

	form.CSRFToken = nosurf.Token(r)

	h.pharseEditStudent(w, form)
}

func (h Handler) StudentUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form := StudentForm{}
	student := storage.Student{ID: uID}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Fatal(err)
	}

	form.Student = student
	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		h.pharseEditStudent(w, form)
		return
	}

	grade, gpa := Grade(student.English, student.Bangla, student.Mathematics)

	student.Grade = grade
	student.GPA = gpa

	_, nErr := h.storage.UpdateStudent(student)
	if nErr != nil {
		log.Println(nErr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)
}
