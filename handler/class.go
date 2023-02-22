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

type ClassForm struct {
	Class     storage.Classes
	CSRFToken string
	FormError map[string]error
}

type ClassFilterList struct {
	Allclass []storage.Classes
	SearchTerm string
}

// For Class Create
func (h Handler) ClassCreate(w http.ResponseWriter, r *http.Request) {
	h.pareseClassTemplate(w, ClassForm{
		CSRFToken: nosurf.Token(r),
	})
}

// For Class Insert
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

// For Show Class List
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

// For Class Edit
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

// For Class Update
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

// For Class Delete
func (h Handler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	if err := h.storage.DeleteSClassByID(uID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/classlist", http.StatusSeeOther)

}

// For Class Already Exists Check By CLassname
func (h Handler) IsClass(w http.ResponseWriter, r *http.Request, classname string) (bool, error) {
	ad, err := h.storage.CheckClassExists(classname)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

