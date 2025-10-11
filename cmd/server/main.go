package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Adit0507/wiki-search-engine/internal/search"
	"github.com/gorilla/mux"
)

type Server struct {
	engine *search.Engine
	tmpl   *template.Template
}

func main() {
	var (
		indexPath = flag.String("index", "./indexes", "Path to indexes")
		port      = flag.Int("port", 8080, "Server port")
	)
	flag.Parse()

	fmt.Println("Loading search engine")
	engine, err := search.NewEngine(*indexPath)
	if err != nil {
		log.Fatal("Failed to create search engine: ", err)
	}

	tpml, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates: ", err)
	}

	server := &Server{
		engine: engine,
		tmpl:   tpml,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", server.handleHome).Methods("GET")
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "index.html", nil)
}
