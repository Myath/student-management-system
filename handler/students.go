package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"log"
	"net/http"
	"sort"
)

type StudentFilterList struct {
	Students []storage.Student
	SearchTerm string
}

type StudentForm struct {
	Student   storage.Student
	FormError map[string]error
	CSRFToken string
}

func (h Handler) StudentsList(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	st := r.FormValue("SearchTerm")
	uf := storage.StudentFilter{
		SearchTerm: st,
	}

	student, err := h.storage.ListStudent(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	t := h.Templates.Lookup("students-list.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	sort.SliceStable(student, func(i, j int) bool {
		return student[i].ID < student[j].ID
	})

	data := StudentFilterList{
		Students:   student,
		SearchTerm: st,
	}

	t.Execute(w, data)
}
