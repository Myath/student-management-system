package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

type StudentSubjectForm struct {
	ID                  int
	FixedStudentSubject []storage.StudentSubject
	StudentSubjecttest storage.StudentSubject
	FixedSubject        []storage.Subjects
	StudentSubjectList  storage.StudentSubject
	CSRFToken           string
	FormError           map[string]error
}

//For Add Mark
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


	t := h.Templates.Lookup("mark.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	data := StudentSubjectForm{
		ID:                  uID,
		FixedStudentSubject: fixedStudentSubject,
		CSRFToken:           nosurf.Token(r),
		FormError:           map[string]error{},
	}

	t.Execute(w, data)
}

// For Mark Store
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

	for a, b := range studentSubject.Mark {
		markinsert := storage.StudentSubject{
			StudentID: studentSubject.StudentID,
			SubjectID: a,
			Marks:     b,
		}
		_, err := h.storage.UpdateMark(markinsert)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

	}

	http.Redirect(w, r, "/admitstudentlist", http.StatusSeeOther)

	// fixedStudentSubject, err := h.storage.GetFixedStudentSubjectByID(studentSubject.StudentID)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "internal server error", http.StatusInternalServerError)
	// }


	// t := h.Templates.Lookup("showadmitstudent.html")
	// if t == nil {
	// 	http.Error(w, "internal server error", http.StatusInternalServerError)
	// }

	// data := StudentSubjectForm{
	// 	FixedStudentSubject: fixedStudentSubject,
	// 	StudentSubjecttest : fixedStudentSubject[0],
	// }

	// t.Execute(w, data)
}