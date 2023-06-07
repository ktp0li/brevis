package main

import (
	"database/sql"
	"fmt"
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

func generateLink() string {
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

func connect() *sql.DB {
	db, err := sql.Open("sqlite3", "urls.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createDB() {
	db := connect()
	defer db.Close()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS url (id INTEGER PRIMARY KEY, link TEXT NOT NULL UNIQUE, url TEXT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully created table!")
	}
	statement.Exec()
}

func addEntry(entry Entry) {
	db := connect()
	defer db.Close()

	statement, _ := db.Prepare("INSERT INTO url (link, url) VALUES (?, ?)")
	statement.Exec(entry.Link, entry.URL)
	log.Printf("Inserted %s into database!", entry.URL)
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var url string

	if r.URL.Path != "/" {
		db.QueryRow("SELECT url FROM url WHERE link = ?", r.URL.Path[1:]).Scan(&url)
		http.Redirect(w, r, url, http.StatusSeeOther)
	} else {
		fmt.Fprintf(w, "pupupu")
	}
}

func main() {
	createDB()
	addEntry(Entry{"puk11", "kak"})
	addEntry(Entry{"poli", "https://ktp0li.su"})

	mux := http.NewServeMux()

	mux.HandleFunc("/", redirectTo)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
