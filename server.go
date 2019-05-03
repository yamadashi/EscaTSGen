package server

import (
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/yamadashi/EscaTSGen/controller"
	"github.com/yamadashi/EscaTSGen/db"
)

type Server struct {
	dbx    *sqlx.DB
	router *mux.Router
}

func (s *Server) Close() error {
	return s.dbx.Close()
}

func (s *Server) Init(dbconf, env string) {
	cs, err := db.NewConfigsFromFile(dbconf)
	if err != nil {
		log.Fatalf("cannot open database configuration. exit. %s", err)
	}
	dbx, err := cs.Open(env)
	if err != nil {
		log.Fatalf("db initialization failed: %s", err)
	}
	s.dbx = dbx
	s.router = s.Route()
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(addr string) {
	log.Printf("start listening on %s", addr)

	// NOTE: when you serve on TLS, make csrf.Secure(true)
	var csrfProtectKey = []byte(os.Getenv("CSRF_KEY"))
	CSRF := csrf.Protect(
		csrfProtectKey, csrf.Secure(false))

	http.ListenAndServe(addr, context.ClearHandler(CSRF(s.router)))
}

//Tokenはとりあえず保留
func (s *Server) Route() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "pong")
	}).Methods("GET")
	// router.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(map[string]string{
	// 		"token": csrf.Token(r),
	// 	})
	// }).Methods("GET")

	history := &controller.History{DB: s.dbx}
	router.Handle("/api/histories", handler(history.GetAll)).Methods("GET")
	router.Handle("/api/histories/{num}", handler(history.Get)).Methods("GET")
	router.Handle("/api/histories", handler(history.PostOne)).Methods("POST")
	router.Handle("/api/histories", handler(history.Delete)).Methods("DELETE")

	//router.Handle("/", http.FileServer(http.Dir("public")))
	//これがだめな理由がわからん
	router.PathPrefix("/").Handler(
		http.FileServer(http.Dir("public")))

	return router
}
