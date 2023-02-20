package handler

import (
	"log"
	"net/http"
)

func (h Handler) pharseCreateStudent(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-students.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h Handler) pharseEditStudent(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("edit-student.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

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

func (h Handler) pareseClassTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("class-create.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h Handler) pareseSubjectTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-subject.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

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