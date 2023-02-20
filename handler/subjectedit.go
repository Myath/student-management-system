package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	// "fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) SubjectEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	editSubjects, err := h.storage.GetSubjectByID(uID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.pharseEditSubject(w, SubjectForm{
		Classlist: classlist,
		Subject:   *editSubjects,
		CSRFToken: nosurf.Token(r),
		FormError: map[string]error{},
	})
}

func (h Handler) SubjectUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	subject := storage.Subjects{}
	subject.ID = uID

	if err := h.decoder.Decode(&subject, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := subject.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			subject.FormError = vErr
		}
		h.pharseEditSubject(w, SubjectForm{
			Subject:   subject,
			Classlist: classlist,
			CSRFToken: nosurf.Token(r),
			FormError: subject.FormError,
		})
		return
	}

	singlesubject, err := h.storage.GetSubjectByID(subject.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if singlesubject.Subject_name != subject.Subject_name || singlesubject.Class_ID != subject.Class_ID {

		checkAlreadyExist, err := h.IsSubject(w, r, subject.Subject_name, subject.Class_ID)

		if err != nil {
			fmt.Println(err)
			return
		}
		if checkAlreadyExist {
			h.pharseEditSubject(w, SubjectForm{
				Subject:   subject,
				Classlist: classlist,
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Subject_name": fmt.Errorf("The Subject already Exist."),
				}})
			return
		}
	}

	_, eRr := h.storage.UpdateSubjects(subject)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/subjectlist", http.StatusSeeOther)
}
