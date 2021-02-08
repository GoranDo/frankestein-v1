package main

import (
	"time"
)

//Book object
type Book struct{
	ID int
	Name string
	Author string
	PublicationDate time.Time
	Pages int


}
//User object
type User struct {
	UserName string
	Password []byte
	First    string
	Last     string
}

//IndexPage list of books
type IndexPage struct{
	AllBooks []Book
}
//BookPage single book
type BookPage struct{
	TargetBook Book
}
//ErrorPage error
type ErrorPage struct{
	ErrorMsg string
}
//PublicationDateStr returns a sanitized Publication Date in the format YYYY-MM-DD
func (b Book) PublicationDateStr() string {
	return b.PublicationDate.Format("2006-01-02")
}