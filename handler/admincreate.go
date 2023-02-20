package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)


func (h Handler) AdminCreate(w http.ResponseWriter, r *http.Request) {
	h.pareseAdminCreateTemplate(w, AdminForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) AdminCreateProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	admin := storage.LoginAdmin{}
	if err := h.decoder.Decode(&admin, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := admin.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			admin.FormError = vErr
		}
		h.pareseAdminCreateTemplate(w, AdminForm{
			Admin:     admin,
			CSRFToken: nosurf.Token(r),
			FormError: admin.FormError,
		})
		return
	}

	checkAlreadyExist, err := h.IsAdmin(w, r, admin.Username, admin.Email)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkAlreadyExist {
		h.pareseAdminCreateTemplate(w, AdminForm{
			Admin:     admin,
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Username": fmt.Errorf("The Username/Email already Exist."),
			}})
		return
	}

	_, eRr := h.storage.CreateAdmin(admin)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
}

func (h Handler) IsAdmin(w http.ResponseWriter, r *http.Request, username, email string) (bool, error) {
	ad, err := h.storage.CheckAdminExists(username, email)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}


