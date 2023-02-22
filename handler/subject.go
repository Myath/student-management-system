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

type SubjectForm struct {
	Classlist []storage.Classes
	Subject   storage.Subjects
	Class storage.Classes
	CSRFToken string
	FormError map[string]error
}
type SubjectFilterList struct {
	AllSubject []storage.Subjects
	SearchTerm string
}

// For Subject Create
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

// For Subject Insert
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

// For Show Subject List
func (h Handler) SubjectList(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	cl := r.FormValue("SearchTerm")
	uf := storage.SubjectFilter{
		SearchTerm: cl,
	}

	subject, err := h.storage.SubjectList(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	t := h.Templates.Lookup("subjectlist.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}


	data := SubjectFilterList{
		AllSubject: subject,
		SearchTerm: cl,
	}

	t.Execute(w, data)
}

// For Subject Edit
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

// For Subject Update
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

// For Subject Delete
func (h Handler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	if err := h.storage.DeleteSubjectByID(uID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/subjectlist", http.StatusSeeOther)

}

// For Subject Already
func (h Handler) IsSubject(w http.ResponseWriter, r *http.Request, subjectname string, class_id int) (bool, error) {
	ad, err := h.storage.CheckSubjectExists(subjectname, class_id)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}