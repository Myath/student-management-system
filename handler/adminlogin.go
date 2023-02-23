package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
	"golang.org/x/crypto/bcrypt"
)

type AdminForm struct {
	Admin     storage.LoginAdmin
	Logintime storage.AdminLogin
	CSRFToken string
	FormError map[string]error
}

func (h Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	h.pareseLoginTemplate(w, AdminForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) AdminLoginProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	admin := storage.AdminLogin{}
	if err := h.decoder.Decode(&admin, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := admin.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			admin.FormError = vErr
		}
		h.pareseLoginTemplate(w, AdminForm{
			Logintime:     admin,
			CSRFToken: nosurf.Token(r),
			FormError: admin.FormError,
		})
		return
	}

	ad, err := h.storage.GetAdminByUsername(admin.Username)

	if err != nil {
		h.pareseLoginTemplate(w, AdminForm{
			Logintime:     admin,
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Username": fmt.Errorf("credentials does not match"),
			}})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(ad.Password), []byte(admin.Password)); err != nil {
		h.pareseLoginTemplate(w, AdminForm{
			Logintime:     admin,
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Username": fmt.Errorf("credentials does not match"),
			}})
		return
	}

	h.sessionManager.Put(r.Context(), "username", ad.Username)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func (h Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	h.pareseDashboardTemplate(w, nil)
}

