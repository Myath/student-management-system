package handler

import (
	"log"
	"net/http"
)


//Phrase For Admin Create Template
func (h Handler) pareseAdminCreateTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("admincreate.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Admin Login Template
func (h Handler) pareseLoginTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("adminlogin.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Dashboard Template
func (h Handler) pareseDashboardTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("dashboard.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Class Create Template
func (h Handler) pareseClassTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("classcreate.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Subject Create Template
func (h Handler) pareseSubjectTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("createsubject.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For AdmitStudent Create Template
func (h Handler) pareseAdmitStudentTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("admitstudent.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Subject Edit Template
func (h Handler) pharseEditSubject(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("subjectedit.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Student Edit Template
func (h Handler) pharseEditStudents(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("admitstudentedit.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Class Edit Template
func (h Handler) pharseEditClass(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("classedit.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Admin Edit Template
func (h Handler) pharseEditAdmin(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("adminedit.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For Mark Add Template
func (h Handler) pharseMarksAdd(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("mark.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

//Phrase For AddAdmin Template
func (h Handler) pareseAddAdminCreateTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("addadmin.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}