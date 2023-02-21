package handler

import (
	"STUDENT-MANAGEMENT-PROJECT/storage"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/form"
)

type Handler struct {
	sessionManager *scs.SessionManager
	decoder        *form.Decoder
	storage        dbstorage
	Templates      *template.Template
}

type dbstorage interface {
	CreateStudent(storage.Student) (*storage.Student, error)
	ListStudent(storage.StudentFilter) ([]storage.Student, error)
	UpdateStudent(storage.Student) (*storage.Student, error)
	GetStudentByID(int) (*storage.Student, error)
	DeleteStudentByID(int) error
	GetAdminByUsername(string) (*storage.LoginAdmin, error)
	CreateAdmin(storage.LoginAdmin) (*storage.LoginAdmin, error)
	CheckStudentExists(string, int) (bool, error)
	CheckAdminExists(string, string) (bool, error)
	ListAdmin(storage.AdminFilter) ([]storage.LoginAdmin, error)
	CreateClass(s storage.Classes) (*storage.Classes, error)
	CheckClassExists(classname string) (bool, error)
	ListClass(uf storage.ClassFilter) ([]storage.Classes, error)
	ListOfClassName() ([]storage.Classes, error)
	CreateSubject(s storage.Subjects) (*storage.Subjects, error)
	CheckSubjectExists(subjectname string, class_id int) (bool, error)
	SubjectList(uf storage.SubjectFilter) ([]storage.Subjects, error)
	AdmitStudentCreate(s storage.AdmitStudents) (*storage.AdmitStudents, error)
	CheckAdmitStudentExists(username, email string) (bool, error)
	CheckAdmitStudentRollExists(roll int) (bool, error)
	AdmitStudentList(uf storage.AdmitStudentFilter) ([]storage.AdmitStudents, error)
	InsertStudentSubject(s storage.StudentSubject) (*storage.StudentSubject, error)
	GetSubjectByClassID(class_id int) ([]storage.Subjects, error)
	DeleteSClassByID(id int) error
	GetSubjectByID(id int) (*storage.Subjects, error)
	UpdateSubjects(s storage.Subjects) (*storage.Subjects, error)
	DeleteSubjectByID(id int) error
	GetAdmitStudentByID(id int) (*storage.AdmitStudents, error)
	UpdateAdminStudent(s storage.AdmitStudents) (*storage.AdmitStudents, error)
	DeleteAdmitStudentByID(id int) error
	GetClassByID(id int) (*storage.Classes, error)
	UpdateClasses(s storage.Classes) (*storage.Classes, error)
	GetAdminByID(id int) (*storage.LoginAdmin, error)
	UpdateAdmin(s storage.LoginAdmin) (*storage.LoginAdmin, error)
	DeleteAdminByID(id int) error
}

const (
	adminLoginPath = "/adminLogin"
)

func NewHandler(sm *scs.SessionManager, formdecoder *form.Decoder, storage dbstorage) *chi.Mux {
	h := &Handler{
		sessionManager: sm,
		decoder:        formdecoder,
		storage:        storage,
	}

	h.ParseTemplates()
	r := chi.NewRouter()

	assetsPrefixForSubjectEdit := "/subject/edit/static/"
	assetsPrefixForSubjectUpdate := "/subject/update/static/"
	assetsPrefixForStudentEdit := "/admitstudent/edit/static/"
	assetsPrefixForStudentUpdate := "/admitstudent/update/static/"
	assetsPrefixForClassEdit := "/class/edit/static/"
	assetsPrefixForClassUpdate := "/class/update/static/"
	assetsPrefixForAdminEdit := "/admin/edit/static/"

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(VerbMethod)

	r.Get("/admincreate", h.AdminCreate)
	r.Post("/admincreate", h.AdminCreateProcess)

	r.Get("/admin/edit/{id:[0-9]+}", h.AdminEdit)
	r.Put("/admin/update/{id:[0-9]+}", h.AdminUpdate)
	r.Get("/admin/delete/{id:[0-9]+}", h.DeleteAdmin)

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Use(h.AuthenticationForLogin)
		r.Get(adminLoginPath, h.AdminLogin)
		r.Post(adminLoginPath, h.AdminLoginProcess)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets/src"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))
	r.Handle(assetsPrefixForSubjectEdit+"*", http.StripPrefix(assetsPrefixForSubjectEdit, http.FileServer(filesDir)))
	r.Handle(assetsPrefixForSubjectUpdate+"*", http.StripPrefix(assetsPrefixForSubjectUpdate, http.FileServer(filesDir)))
	r.Handle(assetsPrefixForStudentEdit+"*", http.StripPrefix(assetsPrefixForStudentEdit, http.FileServer(filesDir)))
	r.Handle(assetsPrefixForStudentUpdate+"*", http.StripPrefix(assetsPrefixForStudentUpdate, http.FileServer(filesDir)))
	r.Handle(assetsPrefixForClassEdit+"*", http.StripPrefix(assetsPrefixForClassEdit, http.FileServer(filesDir)))
	r.Handle(assetsPrefixForClassUpdate+"*", http.StripPrefix(assetsPrefixForClassUpdate, http.FileServer(filesDir)))
	r.Handle(assetsPrefixForAdminEdit+"*", http.StripPrefix(assetsPrefixForAdminEdit, http.FileServer(filesDir)))

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Use(h.Authentication)
		r.Get("/dashboard", h.Dashboard)
		r.Get("/adminlist", h.AdminList)
		r.Get("/classcreate", h.ClassCreate)
		r.Post("/classcreate", h.ClassCreateProcess)
		r.Get("/classlist", h.ClassList)

		r.Get("/class/edit/{id:[0-9]+}", h.EditClass)
		r.Put("/class/update/{id:[0-9]+}", h.ClassUpdate)
		r.Get("/class/delete/{id:[0-9]+}", h.DeleteClass)

		r.Get("/subjectcreate", h.SubjectCreate)
		r.Post("/subjectcreate", h.SubjectCreateProcess)
		r.Get("/subjectlist", h.SubjectList)

		r.Get("/subject/edit/{id:[0-9]+}", h.SubjectEdit)
		r.Put("/subject/update/{id:[0-9]+}", h.SubjectUpdate)
		r.Get("/subject/delete/{id:[0-9]+}", h.DeleteSubject)

		r.Get("/admitstudent", h.AdmitStudent)
		r.Post("/admitstudent", h.AdmitStudentProcess)
		r.Get("/admitstudentlist", h.AdmitStudentlist)

		r.Get("/admitstudent/edit/{id:[0-9]+}", h.AdmitStudentEdit)
		r.Put("/admitstudent/update/{id:[0-9]+}", h.AdmitStudentUpdate)
		r.Get("/admitstudent/delete/{id:[0-9]+}", h.DeleteAdmitStudent)

		r.Route("/student", func(r chi.Router) {

			r.Get("/list", h.StudentsList)

			r.Get("/create", h.CreateStudent)

			r.Post("/store", h.StudentStore)

			r.Get("/{id:[0-9]+}/edit", h.StudentEdit)

			r.Put("/{id:[0-9]+}/update", h.StudentUpdate)

			r.Get("/{id:[0-9]+}/delete", h.DeleteStudent)
		})
		r.Get("/adminlogout", h.AdminLogOut)
	})

	return r
}

func VerbMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := h.sessionManager.GetString(r.Context(), "username")
		if username == "" {
			http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) AuthenticationForLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := h.sessionManager.GetString(r.Context(), "username")
		if username != "" {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) ParseTemplates() error {
	templates := template.New("web-templates").Funcs(template.FuncMap{
		"golabalFunc": func(n string) string {
			return strings.Title(n)
		},
	}).Funcs(sprig.FuncMap())

	newFS := os.DirFS("assets/templates")
	tmpl := template.Must(templates.ParseFS(newFS, "*/*/*.html", "*.html"))
	if tmpl == nil {
		log.Fatalln("unable to parse templates")
	}

	h.Templates = tmpl
	return nil
}
