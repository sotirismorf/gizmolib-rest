// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package database

import ()

type Author struct {
	ID   int64
	Name string
	Bio  string
}

type Book struct {
	ID              int64
	Title           string
	AuthorID        int64
	Description     string
	YearPublished   int16
	CopiesTotal     int32
	CopiesAvailable int32
}

type User struct {
	ID       int64
	Username string
	Password string
}
