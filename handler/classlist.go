package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"log"
	"net/http"
	"sort"
)

type ClassFilterList struct {
	Allclass []storage.Classes
	SearchTerm string
}


func (h Handler) ClassList(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	cl := r.FormValue("SearchTerm")
	uf := storage.ClassFilter{
		SearchTerm: cl,
	}

	class, err := h.storage.ListClass(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	t := h.Templates.Lookup("classlist.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	sort.SliceStable(class, func(i, j int) bool {
		return class[i].ID < class[j].ID
	})

	data := ClassFilterList{
		Allclass:   class,
		SearchTerm: cl,
	}

	t.Execute(w, data)
}