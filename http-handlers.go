package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func handleListBook(w http.ResponseWriter, r*http.Request){
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var books,err = allBooks()
	if err != nil {
		renderErrorPage(w, err)
	}
	buf,err:=ioutil.ReadFile("www/index.html")
	if err != nil {
		renderErrorPage(w, err)
	}
	var page = IndexPage{AllBooks:books}
	indexPage:=string(buf)
	t:=template.Must(template.New("index.html").Parse(indexPage))
	t.Execute(w, page)

}
func handleViewBook(w http.ResponseWriter, r*http.Request){
	params:=r.URL.Query()
	idStr:=params.Get("id")
	var currentBook = Book{}
	currentBook.PublicationDate = time.Now()

	if len(idStr)>0 {
		id,err:=strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w,err)
			return
		}
		currentBook,err = getBook(id)
		if err != nil {
			renderErrorPage(w,err)
			return
		}
	}
	buf,err:= ioutil.ReadFile("www/book.html")
	if err != nil {
		renderErrorPage(w,err)
		return
	}
	var page = BookPage{TargetBook: currentBook}
	bookPage:= string(buf)
	t:=template.Must(template.New("indexPage").Parse(bookPage))
	err=t.Execute(w, page)
	if err != nil {
		renderErrorPage(w,err)
		return
	}
}
func handleSaveBook(w http.ResponseWriter, r *http.Request) {
	var id = 0
	var err error

	r.ParseForm()
	params := r.PostForm
	idStr := params.Get("id")

	if len(idStr) > 0 {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w, err)
			return
		}
	}

	name := params.Get("name")
	author := params.Get("author")

	pagesStr := params.Get("pages")
	pages := 0
	if len(pagesStr) > 0 {
		pages, err = strconv.Atoi(pagesStr)
		if err != nil {
			renderErrorPage(w, err)
			return
		}
	}

	publicationDateStr := params.Get("publicationDate")
	var publicationDate time.Time

	if len(publicationDateStr) > 0 {
		publicationDate, err = time.Parse("2006-01-02", publicationDateStr)
		if err != nil {
			renderErrorPage(w, err)
			return
		}
	}

	if id == 0 {
		_, err = insertBook(name, author, pages, publicationDate)
	} else {
		_, err = updateBook(id, name, author, pages, publicationDate)
	}

	if err != nil {
		renderErrorPage(w, err)
		return
	}

	http.Redirect(w, r, "/book", 302)
}

func handleDeleteBook(w http.ResponseWriter, r*http.Request){
	params:=r.URL.Query()
	idStr:=params.Get("id")
	if len(idStr)>0 {
		id,err:=strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w,err)
			return
		}
		res,err:=removeBook(id)
		if err != nil {
			renderErrorPage(w,err)
			return
		}
		fmt.Printf("removed book with id: %v\n", res)
	}
	http.Redirect(w, r, "/book", 302)
}

func renderErrorPage(w http.ResponseWriter, errorMsg error){
	buff,err:= ioutil.ReadFile("www/error.html")
	if err != nil {
		log.Printf("%v\n", err)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	var page = ErrorPage{ErrorMsg: errorMsg.Error()}
	errorPage:=string(buff)
	t:=template.Must(template.New("errorPage").Parse(errorPage))
	t.Execute(w, page)
}