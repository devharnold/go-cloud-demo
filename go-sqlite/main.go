package main

import (
	"fmt"
	"encoding/json"
	"database/sql"
	"log"
	"net/http"

	_"github.com/mattn/go-sqlite3"
)

type Note struct {
	ID int		`json:"id"`
	Title string `json:"string"`
	Description string	`json:"string"`
	Body string		`json:"string"`
}

func addNotes(db *sql.DB, Title, Description, Body string) error {
	query := `INSERT INTO notes (title, description, body) VALUES (?, ?, ?)`
	_, err := db.Exec(query, Title, Description, Body)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}
	return nil
}

func getNotes(db *sql.DB) ([]Note, error) {
	query := `SELECT id, title, description, body FROM notes`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Description, &n.Body); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		notes = append(notes, n)
	}
	return notes, nil
}

// add note handler
func addNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var n Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	if n.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
	}

	err := addNotes(n.Title, n.Description, n.Body)
	if err != nil {
		http.Error(w, "failed to insert note", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Println(w, "Note added successfully")
}

// get notes handler
func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	notes, err := getNotes()
	if err != nil {
		http.Error(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func main() {
	var err error
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Error creating the table", err)
	}

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addNoteHandler(w, r)
		} else if r.Method == http.MethodGet {
			getNotesHandler(w, r)
		} else {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))


}