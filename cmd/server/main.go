package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	r.HandleFunc("/search", server.handleSearch).Methods("GET")
	r.HandleFunc("/api/search", server.handleApiSearch).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	fmt.Printf("Server starting on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		s.tmpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	results, err := s.engine.Search(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Query   string
		Results []search.Result
	}{
		Query:   query,
		Results: results,
	}

	s.tmpl.ExecuteTemplate(w, "index.html", data)
}

func (s *Server) handleApiSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	query = strings.TrimSpace(query)

	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	results, err := s.engine.Search(query, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
