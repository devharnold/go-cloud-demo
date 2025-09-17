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
	ID 		int		`json:"id"`
	Title 	string `json:"title"`
	Content string		`json:"content"`
}

func addNotes(db *sql.DB, Title, Content string) error {
	query := `INSERT INTO notes (title, content) VALUES (?, ?)`
	_, err := db.Exec(query, Title, Content)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}
	return nil
}

func getNotes(db *sql.DB) ([]Note, error) {
	query := `SELECT id, title, content FROM notes`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		notes = append(notes, n)
	}
	return notes, nil
}

// add note handler
func addNoteHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var n Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if n.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	err := addNotes(db, n.Title, n.Content)
	if err != nil {
		http.Error(w, "failed to insert note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Note added successfully")
}

// get notes handler
func getNotesHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	notes, err := getNotes(db)
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
		content TEXT
	);
	`
	fmt.Println("Creating a table with sql: \n", sqlStmt)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Error creating the table", err)
	}

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addNoteHandler(db, w, r)
		} else if r.Method == http.MethodGet {
			getNotesHandler(db, w, r)
		} else {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}