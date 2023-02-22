package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

type StudentSubjectForm struct {
	ID int
	FixedStudentSubject []storage.StudentSubject
	FixedSubject []storage.Subjects
	StudentSubjectList storage.StudentSubject
	CSRFToken string
	FormError map[string]error
}

func (h Handler) AddMark(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}


	fixedStudentSubject, err := h.storage.GetFixedStudentSubjectByID(uID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var fixedSubjects []storage.Subjects
	for _, stu := range fixedStudentSubject{
		p, eRR := h.storage.GetFixedSubjectByID(stu.SubjectID)
		if eRR != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		fixedSubjects = append(fixedSubjects, p)
	}

	h.pharseMarksAdd(w, StudentSubjectForm{
		ID : uID,
		FixedStudentSubject: fixedStudentSubject,
		FixedSubject: fixedSubjects,
		CSRFToken:    nosurf.Token(r),
		FormError:    map[string]error{},
	})
}


func (h Handler) Markstore(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	studentSubject := storage.StudentSubject{}
	if err := h.decoder.Decode(&studentSubject, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	fmt.Println(studentSubject)

	// classlist, err := h.storage.ListOfClassName()
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "internal server error", http.StatusInternalServerError)
	// }

	// if err := subject.Validate(); err != nil {
	// 	if vErr, ok := err.(validation.Errors); ok {
	// 		subject.FormError = vErr
	// 	}
	// 	h.pareseSubjectTemplate(w, SubjectForm{
	// 		Subject:     subject,
	// 		Classlist: classlist,
	// 		CSRFToken: nosurf.Token(r),
	// 		FormError: subject.FormError,
	// 	})
	// 	return
	// }

	// checkAlreadyExist, err := h.IsSubject(w, r, subject.Subject_name, subject.Class_ID)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// if checkAlreadyExist {
	// 	h.pareseSubjectTemplate(w, SubjectForm{
	// 		Subject:     subject,
	// 		Classlist: classlist,
	// 		CSRFToken: nosurf.Token(r),
	// 		FormError: map[string]error{
	// 			"Subject_name": fmt.Errorf("The Subject already Exist."),
	// 		}})
	// 	return
	// }

	// _, eRr := h.storage.CreateSubject(subject)
	// if eRr != nil {
	// 	log.Println(err)
	// 	http.Error(w, "internal server error", http.StatusInternalServerError)
	// }
	http.Redirect(w, r, "/subjectlist", http.StatusSeeOther)
}