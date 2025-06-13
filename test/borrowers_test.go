package test

import (
	"main/borrowers"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBorrowedBooks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %s", err)
	}
	defer db.Close()

	borrowerId := 3

	rows := mock.NewRows([]string{"BookId", "Title", "Author", "PageCount", "BorrowedBy", "CreatedAt"})
	rows.AddRow(5, "Test book", "Test author", 300, &borrowerId, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `BookId`,`Title`,`Author`,`PageCount`,`BorrowedBy`,`CreatedAt` FROM BOOKS WHERE `BorrowedBy`=?;")).
		WithArgs(borrowerId).
		WillReturnRows(rows)

	books, err := borrowers.BorrowedBooks(db, t.Context(), borrowerId)
	if err != nil {
		t.Error(err)
	}

	t.Log(books)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
