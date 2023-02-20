package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"log"
	"net/http"
)

type SubjectFilterList struct {
	AllSubject []storage.Subjects
	SearchTerm string
}

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
