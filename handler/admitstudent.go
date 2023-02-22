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

type AdmitStudentForm struct {
	Classlist []storage.Classes
	Student   storage.AdmitStudents
	Class     storage.Classes
	CSRFToken string
	FormError map[string]error
}

type AdmitStudentFilterList struct {
	AllAdmitStudent []storage.AdmitStudents
	SearchTerm      string
}

// For AdmitStudent Create
func (h Handler) AdmitStudent(w http.ResponseWriter, r *http.Request) {
	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.pareseAdmitStudentTemplate(w, AdmitStudentForm{
		Classlist: classlist,
		CSRFToken: nosurf.Token(r),
	})
}

// For AdmitStudent Insert
func (h Handler) AdmitStudentProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	student := storage.AdmitStudents{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		h.pareseAdmitStudentTemplate(w, AdmitStudentForm{
			Classlist: classlist,
			Student:   student,
			CSRFToken: nosurf.Token(r),
			FormError: student.FormError,
		})
		return
	}

	checkUsernameAlreadyExist, err := h.IsAdmitStudentUsername(w, r, student.Username)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkUsernameAlreadyExist {
		h.pareseAdmitStudentTemplate(w, AdmitStudentForm{
			Classlist: classlist,
			Student:   student,
			Class:     storage.Classes{},
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Username": fmt.Errorf("The username already Exist."),
			}})
		return
	}

	checkEmailAlreadyExist, eRr := h.IsAdmitStudentEmail(w, r, student.Email)

	if eRr != nil {
		fmt.Println(eRr)
		return
	}
	if checkEmailAlreadyExist {
		h.pareseAdmitStudentTemplate(w, AdmitStudentForm{
			Classlist: classlist,
			Student:   student,
			Class:     storage.Classes{},
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Email": fmt.Errorf("The username already Exist."),
			}})
		return
	}

	checkRollAlreadyExist, eRR := h.IsAdmitStudentRoll(w, r, student.Roll)

	if eRR != nil {
		fmt.Println(eRR)
		return
	}
	if checkRollAlreadyExist {
		h.pareseAdmitStudentTemplate(w, AdmitStudentForm{
			Classlist: classlist,
			Student:   student,
			Class:     storage.Classes{},
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Roll": fmt.Errorf("The roll already Exist."),
			}})
		return
	}

	data, er := h.storage.AdmitStudentCreate(student)
	if er != nil {
		log.Println(eRr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	eRRor := h.StudentSubjectHandler(w, r, student.Class_ID, data.ID)
	if eRRor != nil {
		log.Println(eRR)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admitstudentlist", http.StatusSeeOther)
}

// For AdmitStudent List
func (h Handler) AdmitStudentlist(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	cl := r.FormValue("SearchTerm")
	uf := storage.AdmitStudentFilter{
		SearchTerm: cl,
	}

	admitstudent, err := h.storage.AdmitStudentList(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	t := h.Templates.Lookup("admitstudentlist.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	data := AdmitStudentFilterList{
		AllAdmitStudent: admitstudent,
		SearchTerm:      cl,
	}

	t.Execute(w, data)
}

// For AdmitStudent Edit
func (h Handler) AdmitStudentEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	editStudents, err := h.storage.GetAdmitStudentByID(uID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.pharseEditStudents(w, AdmitStudentForm{
		Classlist: classlist,
		Student:   *editStudents,
		CSRFToken: nosurf.Token(r),
	})
}

// For AdmitStudent Update
func (h Handler) AdmitStudentUpdate(w http.ResponseWriter, r *http.Request) {
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

	student := storage.AdmitStudents{}
	student.ID = uID
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	classlist, err := h.storage.ListOfClassName()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		h.pharseEditStudents(w, AdmitStudentForm{
			Classlist: classlist,
			Student:   student,
			CSRFToken: nosurf.Token(r),
			FormError: student.FormError,
		})
		return
	}

	singlestudent, err := h.storage.GetAdmitStudentByID(student.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if singlestudent.Username != student.Username {
		checkUsernameAlreadyExist, err := h.IsAdmitStudentUsername(w, r, student.Username)
		if err != nil {
			fmt.Println(err)
			return
		}

		if checkUsernameAlreadyExist {
			h.pharseEditStudents(w, AdmitStudentForm{
				Classlist: classlist,
				Student:   student,
				Class:     storage.Classes{},
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Username": fmt.Errorf("The username already Exist."),
				}})
			return
		}

	}

	if singlestudent.Email != student.Email {
		checkEmailAlreadyExist, err := h.IsAdmitStudentEmail(w, r, student.Email)
		if err != nil {
			fmt.Println(err)
			return
		}

		if checkEmailAlreadyExist {
			h.pharseEditStudents(w, AdmitStudentForm{
				Classlist: classlist,
				Student:   student,
				Class:     storage.Classes{},
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Email": fmt.Errorf("The email already Exist."),
				}})
			return
		}

	} 
	
	
	if singlestudent.Roll != student.Roll {
		checkRollAlreadyExist, err := h.IsAdmitStudentRoll(w, r, student.Roll)

		if err != nil {
			fmt.Println(err)
			return
		}

		if checkRollAlreadyExist {
			h.pharseEditStudents(w, AdmitStudentForm{
				Classlist: classlist,
				Student:   student,
				Class:     storage.Classes{},
				CSRFToken: nosurf.Token(r),
				FormError: map[string]error{
					"Roll": fmt.Errorf("The roll already Exist."),
				}})
			return
		}
	}

	_, eRr := h.storage.UpdateAdmitStudent(student)
	if eRr != nil {
		log.Println(eRr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admitstudentlist", http.StatusSeeOther)
}

// For AdmitStudent Delete
func (h Handler) DeleteAdmitStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	if err := h.storage.DeleteAdmitStudentByID(uID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admitstudentlist", http.StatusSeeOther)

}

// For AdmitStudent Username Already Exist Check
func (h Handler) IsAdmitStudentUsername(w http.ResponseWriter, r *http.Request, username string) (bool, error) {
	ad, err := h.storage.CheckAdmitStudentUsernameExists(username)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

// For AdmitStudent Email Already Exist Check
func (h Handler) IsAdmitStudentEmail(w http.ResponseWriter, r *http.Request, email string) (bool, error) {
	ad, err := h.storage.CheckAdmitStudentUsernameExists(email)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

// For AdmitStudent Roll Already Exist Check
func (h Handler) IsAdmitStudentRoll(w http.ResponseWriter, r *http.Request, roll int) (bool, error) {
	ad, err := h.storage.CheckAdmitStudentRollExists(roll)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

// For StudentSubject Insert Function
func (h Handler) StudentSubjectHandler(w http.ResponseWriter, r *http.Request, classID int, studentID int) error {

	subject, err := h.storage.GetSubjectByClassID(classID)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	for _, s := range subject {
		b := storage.StudentSubject{
			StudentID: studentID,
			SubjectID: s.ID,
			Marks:     0,
		}

		_, err := h.storage.InsertStudentSubject(b)
		if err != nil {
			log.Fatalf("%v", err)
			return err
		}
	}
	return nil
}