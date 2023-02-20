package handler

import (
	// "STUDENT-MANAGEMENT-PROJECT/storage"
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	// "fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) EditClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	editClass, err := h.storage.GetClassByID(uID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.pharseEditClass(w, ClassForm{
		Class:     *editClass,
		CSRFToken: nosurf.Token(r),
		FormError: map[string]error{},
	})
}

func (h Handler) ClassUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	class := storage.Classes{}
	class.ID = uID

	if err := h.decoder.Decode(&class, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := class.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			class.FormError = vErr
		}
		h.pharseEditClass(w, ClassForm{
			Class:     class,
			CSRFToken: nosurf.Token(r),
			FormError: class.FormError,
		})
		return
	}

	classlist, err := h.storage.GetClassByID(class.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if classlist.Class_name != class.Class_name {
		
		checkAlreadyExist, err := h.IsClass(w, r, class.Class_name)

		if err != nil {
			fmt.Println(err)
			return
		}
		if checkAlreadyExist {
			h.pharseEditClass(w, ClassForm{
				Class:     class,
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Class_name": fmt.Errorf("The Class already Exist."),
				}})
			return
		}
	}

	_, eRr := h.storage.UpdateClasses(class)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/classlist", http.StatusSeeOther)
}

func GetClassByID(i int) {
	panic("unimplemented")
}
