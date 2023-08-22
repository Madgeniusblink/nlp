package main

import (
	"encoding/json"
	"expvar"
	_ "expvar"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/madgeniusblink/nlp"
	"github.com/madgeniusblink/nlp/stemmer"
)

var (
	numTok = expvar.NewInt("tokenize.calls")
)

func main() {
	// create server (dependency injection)
	logger := log.New(log.Writer(), "[nlp]", log.LstdFlags|log.Lshortfile)
	s := &Server{
		logger: logger, // dependency injection
	}
	// routing
	// health is an exact match
	// /health/ is a prefix match
	// http.HandleFunc("/health", healthHandler)
	// http.HandleFunc("/tokenize", tokenizeHandler)

	r := mux.NewRouter()
	r.HandleFunc("/health", s.healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/tokenize", s.tokenizeHandler).Methods(http.MethodPost)
	r.HandleFunc("/stem/{word}", s.stemHandler).Methods(http.MethodGet)

	http.Handle("/", r)

	//run server
	addr := ":8080"
	s.logger.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

}

type Server struct {
	logger *log.Logger
}

func (s *Server) stemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]
	stem := stemmer.Stem(word)
	fmt.Fprintf(w, "%s -> %s\n", word, stem)
}

// exercise: write a tokenizeHandler that will read the text from the request
// body and return a JSON in the format: "{"tokens": ["who", "on", "first""]}"
func (s *Server) tokenizeHandler(w http.ResponseWriter, r *http.Request) {
	//Before Gorilla/mux
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	numTok.Add(1)
	// Read the request body
	// good idea to not read the entire body from the internet
	rdr := io.LimitReader(r.Body, 1_000_000)

	body, err := io.ReadAll(rdr)
	if err != nil {
		s.logger.Printf("ERROR: %s", err)
		http.Error(w, "Error reading request body",
			http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// ensure the body is closed
	// would the server do that? perhaps we don't need to do this
	defer r.Body.Close()
	// convert bite slice to string
	text := string(body)
	// Assume tokenize function exists to tokenize the text
	tokens := nlp.Tokenize(text)

	// create a struct to format the response
	response := map[string]any{
		"tokens": tokens,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling JSON (cant encode)", http.StatusInternalServerError)
		return
	}

	// write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: run a health check
	fmt.Fprintln(w, "Okay")

}
