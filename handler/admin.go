package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

type AdminFilterList struct {
	Alladmin []storage.LoginAdmin
	SearchTerm string
}

//For AdminCreate
func (h Handler) AdminCreate(w http.ResponseWriter, r *http.Request) {
	h.pareseAdminCreateTemplate(w, AdminForm{
		CSRFToken: nosurf.Token(r),
	})
}

//For AdminCreate Insert
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

//For SHow AdminList
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

//For AdminEdit
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

//For AdminUpdate
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

//For AdminDelete
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

//For Admin Already exists
func (h Handler) IsAdmin(w http.ResponseWriter, r *http.Request, username, email string) (bool, error) {
	ad, err := h.storage.CheckAdminExists(username, email)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

func (h Handler) AdminAdd(w http.ResponseWriter, r *http.Request) {
	h.pareseAddAdminCreateTemplate(w, AdminForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) AddAdminCreateProcess(w http.ResponseWriter, r *http.Request) {
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
		h.pareseAddAdminCreateTemplate(w, AdminForm{
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
		h.pareseAddAdminCreateTemplate(w, AdminForm{
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
	http.Redirect(w, r, "/adminlist", http.StatusSeeOther)
}

