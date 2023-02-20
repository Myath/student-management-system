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

func (h Handler) AdminEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	editAdmin, err := h.storage.GetAdminByID(uID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.pharseEditAdmin(w, AdminForm{
		Admin:     *editAdmin,
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	admin := storage.LoginAdmin{}
	admin.ID = uID
	if err := h.decoder.Decode(&admin, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := admin.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			admin.FormError = vErr
		}
		h.pharseEditAdmin(w, AdminForm{
			Admin:     admin,
			CSRFToken: nosurf.Token(r),
			FormError: admin.FormError,
		})
		return
	}

	singleAdmin, err := h.storage.GetAdminByID(admin.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if singleAdmin.Username != admin.Username || singleAdmin.Email != admin.Email {

		checkAlreadyExist, err := h.IsAdmin(w, r, admin.Username, admin.Email)

		if err != nil {
			fmt.Println(err)
			return
		}
		if checkAlreadyExist {
			h.pharseEditAdmin(w, AdminForm{
				Admin:     admin,
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Username": fmt.Errorf("The Username/Email already Exist."),
				}})
			return
		}
	}

	_, eRr := h.storage.UpdateAdmin(admin)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	
	http.Redirect(w, r, "/adminlist", http.StatusSeeOther)
}
