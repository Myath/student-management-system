package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"log"
	"net/http"
)

type AdmitStudentFilterList struct {
	AllAdmitStudent []storage.AdmitStudents
	SearchTerm      string
}

func (h Handler) AdmitStudentlist(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	cl := r.FormValue("SearchTerm")
	uf := storage.AdmitStudentFilter{
		SearchTerm: cl,
	}

	admitstudent, err := h.storage.AdmitStudentList(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	t := h.Templates.Lookup("admitstudentlist.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	data := AdmitStudentFilterList{
		AllAdmitStudent: admitstudent,
		SearchTerm:      cl,
	}

	t.Execute(w, data)
}
