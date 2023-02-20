package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h Handler) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	if err := h.storage.DeleteAdminByID(uID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/adminlist", http.StatusSeeOther)

}
