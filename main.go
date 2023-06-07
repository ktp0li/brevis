package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Link string
	URL  string
}

func GenerateLink() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	const zero, z = 48, 122 // ASCII

	stroke := [12]string{}

	for i := range stroke {
		rand := 48 + r1.Intn(z-zero+1)
		// look only for digit and alphabetic chars
		for 58 <= rand && rand <= 64 || 91 <= rand && rand <= 96 {
			rand = 48 + r1.Intn(z-zero+1)
		}

		stroke[i] = string(rune(rand))
	}

	return strings.Join(stroke[:], "")
}

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "urls.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func CreateDB(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS url (id INTEGER PRIMARY KEY, link TEXT NOT NULL UNIQUE, url TEXT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully created table!")
	}
	statement.Exec()
}

func AddEntry(db *sql.DB, entry Entry) {
	statement, _ := db.Prepare("INSERT INTO url (link, url) VALUES (?, ?)")
	statement.Exec(entry.Link, entry.URL)
	log.Printf("Inserted %s into database!", entry.URL)
}

func RedirectTo(mux *http.ServeMux) {
	mux.Handle("/", http.RedirectHandler("https://ktp0li.su", http.StatusSeeOther))
}

func main() {
	db := Connect()
	defer db.Close()

	CreateDB(db)
	AddEntry(db, Entry{"puk11", "kak"})

	mux := http.NewServeMux()
	http.ListenAndServe(":8080", mux)
	RedirectTo(mux)
}
