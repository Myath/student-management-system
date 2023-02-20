package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

type ClassForm struct {
	Class     storage.Classes
	CSRFToken string
	FormError map[string]error
}



func (h Handler) ClassCreate(w http.ResponseWriter, r *http.Request) {
	h.pareseClassTemplate(w, ClassForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) ClassCreateProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	class := storage.Classes{}
	if err := h.decoder.Decode(&class, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := class.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			class.FormError = vErr
		}
		h.pareseClassTemplate(w, ClassForm{
			Class:     class,
			CSRFToken: nosurf.Token(r),
			FormError: class.FormError,
		})
		return
	}

	checkAlreadyExist, err := h.IsClass(w, r, class.Class_name)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkAlreadyExist {
		h.pareseClassTemplate(w, ClassForm{
			Class:     class,
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Class_name": fmt.Errorf("The Class already Exist."),
			}})
		return
	}

	_, eRr := h.storage.CreateClass(class)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/classlist", http.StatusSeeOther)
}

func (h Handler) IsClass(w http.ResponseWriter, r *http.Request, classname string) (bool, error) {
	ad, err := h.storage.CheckClassExists(classname)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

