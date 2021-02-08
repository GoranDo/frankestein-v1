package main

import (
	"fmt"
	"database/sql"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
)
const (
	host     = "gorandb.ca9rjrw9e8wu.eu-central-1.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "mydb"
  )
var db *sql.DB

var tpl *template.Template
var dbUsers = map[string]User{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID

func init() {
	tmpDB, err := sql.Open("postgres", "host=gorandb.ca9rjrw9e8wu.eu-central-1.rds.amazonaws.com port=5432 user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpDB.Ping()
	if err != nil {
	  panic(err)
	}
	fmt.Println("Successfully connected!")
	fmt.Println(tmpDB)
	db = tmpDB

	tpl = template.Must(template.ParseGlob("templates/*"))
	bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	dbUsers["test@test.com"] = User{"test@test.com", bs, "James", "Bond"}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("www/assets"))))
}

func main(){
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/book", handleListBook)
	http.HandleFunc("/book.html", handleViewBook)
	http.HandleFunc("/save", handleSaveBook)
	http.HandleFunc("/delete", handleDeleteBook)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	tpl.ExecuteTemplate(w, "indexOne.html", u)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/book", http.StatusSeeOther)
		return
	}
	var u User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		// username taken?
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un
		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u = User{un, bs, f, l}
		dbUsers[un] = u
		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.html", u)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/book", http.StatusSeeOther)
		return
	}
	var u User
	// process form submission
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		// is there a username?
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// does the entered password match the stored password?
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.html", u)
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}