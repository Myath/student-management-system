package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

type SubjectForm struct {
	Classlist []storage.Classes
	Subject   storage.Subjects
	Class storage.Classes
	CSRFToken string
	FormError map[string]error
}

func (h Handler) SubjectCreate(w http.ResponseWriter, r *http.Request) {
	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.pareseSubjectTemplate(w, SubjectForm{
		Classlist: classlist,
		CSRFToken: nosurf.Token(r),
	})
}


func (h Handler) SubjectCreateProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	subject := storage.Subjects{}
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
		h.pareseSubjectTemplate(w, SubjectForm{
			Subject:     subject,
			Classlist: classlist,
			CSRFToken: nosurf.Token(r),
			FormError: subject.FormError,
		})
		return
	}

	checkAlreadyExist, err := h.IsSubject(w, r, subject.Subject_name, subject.Class_ID)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkAlreadyExist {
		h.pareseSubjectTemplate(w, SubjectForm{
			Subject:     subject,
			Classlist: classlist,
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Subject_name": fmt.Errorf("The Subject already Exist."),
			}})
		return
	}

	_, eRr := h.storage.CreateSubject(subject)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/subjectlist", http.StatusSeeOther)
}

func (h Handler) IsSubject(w http.ResponseWriter, r *http.Request, subjectname string, class_id int) (bool, error) {
	ad, err := h.storage.CheckSubjectExists(subjectname, class_id)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}
