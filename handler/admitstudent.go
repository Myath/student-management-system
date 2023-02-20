package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"fmt"
	"log"
	"net/http"

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

	checkAlreadyExist, err := h.IsAdmitStudent(w, r, student.Username, student.Email)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkAlreadyExist {
		h.pareseAdmitStudentTemplate(w, AdmitStudentForm{
			Classlist: classlist,
			Student:   student,
			Class:     storage.Classes{},
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"Username": fmt.Errorf("The username/email already Exist."),
			}})
		return
	}

	checkRollAlreadyExist, err := h.IsAdmitStudentRoll(w, r, student.Roll)

	if err != nil {
		fmt.Println(err)
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

	data, eRr := h.storage.AdmitStudentCreate(student)
	if eRr != nil {
		log.Println(eRr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	eRR := h.MarksHandler(w, r, student.Class_ID, data.ID)
	if eRR != nil {
		log.Println(eRR)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admitstudentlist", http.StatusSeeOther)
}

func (h Handler) IsAdmitStudent(w http.ResponseWriter, r *http.Request, username, email string) (bool, error) {
	ad, err := h.storage.CheckAdmitStudentExists(username, email)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

func (h Handler) IsAdmitStudentRoll(w http.ResponseWriter, r *http.Request, roll int) (bool, error) {
	ad, err := h.storage.CheckAdmitStudentRollExists(roll)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ad, nil
}

func (h Handler) MarksHandler(w http.ResponseWriter, r *http.Request, classID int, studentID int) error {

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

		_, err := h.storage.InsertMark(b)
		if err != nil {
			log.Fatalf("%v", err)
			return err
		}
	}
	return nil
}
