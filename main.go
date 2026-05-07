package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

// TODO: provide an endpoint to check the grammar of the note.
// TODO: provide an endpoint to save the note that can be passed in as Markdown text.
// TODO: provide an endpoint to list the saved notes (i.e. uploaded markdown files).
// TODO: Return the HTML version of the Markdown note (rendered note) through another endpoint.

type Note struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Mu      sync.RWMutex
}

func newNote() *Note {
	return &Note{
		Name:    "",
		Content: "",
	}
}

func saveNotesHandler(rw http.ResponseWriter, req *http.Request) {
	note := newNote()
	defer req.Body.Close()
	if err := json.NewDecoder(req.Body).Decode(&note); err != nil {
		http.Error(rw, "Failed to parse request body as json", http.StatusInternalServerError)
		return
	}
	file, err := os.Create(note.Name)
	if err != nil {
		http.Error(rw, "Failed to create note file", http.StatusInternalServerError)
		return
	}
	file.WriteString(note.Content)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	msg := `{"message": "successfully saved ` + note.Name + `"}`
	rw.Write([]byte(msg))
}

func listNotes(rw http.ResponseWriter, req *http.Request) {}

func checkGrammar(rw http.ResponseWriter, req *http.Request) {}

func parseHtml(rw http.ResponseWriter, req *http.Request) {}

func homeHandler(rw http.ResponseWriter, req *http.Request) {
	apiDocs := map[string]string{
		"/save-notes":    "save markdown notes",
		"/list-notes":    "list saved markdown notes",
		"/check-grammar": "checks grammar of markdown notes",
		"/parse-html":    "parse markdown note into html",
	}
	payload, err := json.Marshal(apiDocs)
	if err != nil {
		http.Error(rw, "Error: Failed to generate json payload", http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(payload)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/save-notes", saveNotesHandler)
	// mux.HandleFunc("/list-notes", listNotes)
	// mux.HandleFunc("/check-grammar", checkGrammar)
	// mux.HandleFunc("/parse-html", parseHtml)
	http.ListenAndServe(":3000", mux)
}
