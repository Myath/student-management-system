package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"log"
	"net/http"
	"sort"
)

type AdminFilterList struct {
	Alladmin []storage.LoginAdmin
	SearchTerm string
}

func (h Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	ad := r.FormValue("SearchTerm")
	uf := storage.AdminFilter{
		SearchTerm: ad,
	}

	admin, err := h.storage.ListAdmin(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	t := h.Templates.Lookup("adminlist.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	sort.SliceStable(admin, func(i, j int) bool {
		return admin[i].ID < admin[j].ID
	})

	data := AdminFilterList{
		Alladmin:   admin,
		SearchTerm: ad,
	}

	t.Execute(w, data)
}
