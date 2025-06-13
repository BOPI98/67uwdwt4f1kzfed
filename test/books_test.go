package test

import (
	"main/books"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %s", err)
	}
	defer db.Close()

	b := books.Book{
		Title:      "Test Book",
		Author:     "Test Name",
		PageCount:  300,
		BorrowedBy: nil,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO BOOKS(`Title`,`Author`,`PageCount`,`BorrowedBy`) VALUES( ?,?,?,? );")).
		WithArgs(b.Title, b.Author, b.PageCount, b.BorrowedBy).
		WillReturnResult(sqlmock.NewResult(33, 1))
	mock.ExpectCommit()

	bookId, err := books.AddBooks(db, t.Context(), b)
	if err != nil {
		t.Error(err)
	}

	if bookId != 33 {
		t.Error("Incorrect bookId!")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
