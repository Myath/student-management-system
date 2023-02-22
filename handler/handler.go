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
	// For Admin
	CreateAdmin(storage.LoginAdmin) (*storage.LoginAdmin, error)
	ListAdmin(storage.AdminFilter) ([]storage.LoginAdmin, error)
	UpdateAdmin(s storage.LoginAdmin) (*storage.LoginAdmin, error)
	DeleteAdminByID(id int) error
	GetAdminByID(id int) (*storage.LoginAdmin, error)
	GetAdminByUsername(string) (*storage.LoginAdmin, error)
	CheckAdminExists(string, string) (bool, error)

	// For Class
	CreateClass(s storage.Classes) (*storage.Classes, error)
	ListClass(uf storage.ClassFilter) ([]storage.Classes, error)
	UpdateClasses(s storage.Classes) (*storage.Classes, error)
	DeleteSClassByID(id int) error
	GetClassByID(id int) (*storage.Classes, error)
	ListOfClassName() ([]storage.Classes, error)
	CheckClassExists(classname string) (bool, error)
	
	// For Subject
	CreateSubject(s storage.Subjects) (*storage.Subjects, error)
	SubjectList(uf storage.SubjectFilter) ([]storage.Subjects, error)
	UpdateSubjects(s storage.Subjects) (*storage.Subjects, error)
	DeleteSubjectByID(id int) error
	GetSubjectByID(id int) (*storage.Subjects, error)
	CheckSubjectExists(subjectname string, class_id int) (bool, error)


	// For AdmitSTudent
	AdmitStudentCreate(s storage.AdmitStudents) (*storage.AdmitStudents, error)
	AdmitStudentList(uf storage.AdmitStudentFilter) ([]storage.AdmitStudents, error)
	UpdateAdmitStudent(s storage.AdmitStudents) (*storage.AdmitStudents, error)
	DeleteAdmitStudentByID(id int) error
	CheckAdmitStudentUsernameExists(username string) (bool, error)
	CheckAdmitStudenEmailtExists(email string) (bool, error)
	CheckAdmitStudentRollExists(roll int) (bool, error)
	GetAdmitStudentByID(id int) (*storage.AdmitStudents, error)


	// For StudentSubject
	GetFixedStudentSubjectByID(id int) ([]storage.StudentSubject, error)
	InsertStudentSubject(s storage.StudentSubject) (*storage.StudentSubject, error)
	GetSubjectByClassID(class_id int) ([]storage.Subjects, error)
	UpdateMark(s storage.StudentSubject) (*storage.StudentSubject, error)
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

	//Make Template prefix variable
	assetsPrefixForSubjectEdit := "/subject/edit/static/"
	assetsPrefixForSubjectUpdate := "/subject/update/static/"
	assetsPrefixForStudentEdit := "/admitstudent/edit/static/"
	assetsPrefixForStudentUpdate := "/admitstudent/update/static/"
	assetsPrefixForClassEdit := "/class/edit/static/"
	assetsPrefixForClassUpdate := "/class/update/static/"
	assetsPrefixForAdminEdit := "/admin/edit/static/"
	assetsPrefixForAddMark := "/addmark/static/"

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

	// For Template Asset Prefixes
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
	r.Handle(assetsPrefixForAddMark+"*", http.StripPrefix(assetsPrefixForAddMark, http.FileServer(filesDir)))
	

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Use(h.Authentication)
		r.Get("/dashboard", h.Dashboard)
		r.Get("/adminlist", h.AdminList)

		// Class Route
		r.Get("/classcreate", h.ClassCreate)
		r.Post("/classcreate", h.ClassCreateProcess)
		r.Get("/classlist", h.ClassList)
		r.Get("/class/edit/{id:[0-9]+}", h.EditClass)
		r.Put("/class/update/{id:[0-9]+}", h.ClassUpdate)
		r.Get("/class/delete/{id:[0-9]+}", h.DeleteClass)

		// Subject Route
		r.Get("/subjectcreate", h.SubjectCreate)
		r.Post("/subjectcreate", h.SubjectCreateProcess)
		r.Get("/subjectlist", h.SubjectList)
		r.Get("/subject/edit/{id:[0-9]+}", h.SubjectEdit)
		r.Put("/subject/update/{id:[0-9]+}", h.SubjectUpdate)
		r.Get("/subject/delete/{id:[0-9]+}", h.DeleteSubject)

		// Admit Student
		r.Get("/admitstudent", h.AdmitStudent)
		r.Post("/admitstudent", h.AdmitStudentProcess)
		r.Get("/admitstudentlist", h.AdmitStudentlist)
		r.Get("/admitstudent/edit/{id:[0-9]+}", h.AdmitStudentEdit)
		r.Put("/admitstudent/update/{id:[0-9]+}", h.AdmitStudentUpdate)
		r.Get("/admitstudent/delete/{id:[0-9]+}", h.DeleteAdmitStudent)

		r.Get("/addmark/{id:[0-9]+}", h.AddMark)
		r.Post("/markstore", h.Markstore)

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
	tmpl := template.Must(templates.ParseFS(newFS, "*/*.html", "*.html"))
	if tmpl == nil {
		log.Fatalln("unable to parse templates")
	}

	h.Templates = tmpl
	return nil
}
