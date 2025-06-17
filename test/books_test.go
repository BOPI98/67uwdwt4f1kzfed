package test

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/BOPI98/67uwdwt4f1kzfed/books"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListBooks(t *testing.T) {
	type testCase struct {
		Name             string
		ExpectedRowsData [][]driver.Value
	}
	cases := []testCase{}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %s", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"BookId", "Title", "Author", "PageCount", "BorrowedBy", "CreatedAt"})
			for _, r := range tc.ExpectedRowsData {
				rows.AddRow(r...)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT `BookId`,`Title`,`Author`,`PageCount`,`BorrowedBy`,`CreatedAt` FROM BOOKS;")).
				WillReturnRows(rows).
				RowsWillBeClosed()

			list, err := books.ListBooks(db, t.Context())
			if err != nil {
				t.Fatal(err)
			}

			if len(list) != len(tc.ExpectedRowsData) {
				t.Error(fmt.Errorf("Result list doesnt match the expected list length! result:%d expected:%d", len(list), len(tc.ExpectedRowsData)))
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAddBook(t *testing.T) {
	type testCase struct {
		Name         string
		Book         books.Book
		LastInsertId int
		RowsAffected int
	}

	cases := []testCase{
		testCase{
			Name: "test1",
			Book: books.Book{
				Title:      "Test book1",
				Author:     "Test author",
				PageCount:  300,
				BorrowedBy: nil,
			},
			LastInsertId: 1,
			RowsAffected: 1,
		},
		testCase{
			Name: "test2",
			Book: books.Book{
				Title:      "Test book2",
				Author:     "Test author2",
				PageCount:  500,
				BorrowedBy: new(int),
			},
			LastInsertId: 2,
			RowsAffected: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %s", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO BOOKS(`Title`,`Author`,`PageCount`,`BorrowedBy`) VALUES( ?,?,?,? );")).
				WithArgs(tc.Book.Title, tc.Book.Author, tc.Book.PageCount, tc.Book.BorrowedBy).
				WillReturnResult(sqlmock.NewResult(int64(tc.LastInsertId), int64(tc.RowsAffected)))
			mock.ExpectCommit()

			bookId, err := books.AddBooks(db, t.Context(), tc.Book)
			if err != nil {
				t.Fatal(err)
			}

			if bookId != tc.LastInsertId {
				t.Error("Incorrect bookId!")
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}

}

func TestBorrowBook(t *testing.T) {
	type testCase struct {
		Name                 string
		BorrowerId           int
		ExpectedBorrowerRows [][]driver.Value
		BookId               int
		ExpectedBookRows     [][]driver.Value
		ExpectedError        error
	}

	cases := []testCase{
		testCase{
			Name:       "borrow book",
			BorrowerId: 1,
			ExpectedBorrowerRows: [][]driver.Value{
				[]driver.Value{1},
			},
			BookId: 2,
			ExpectedBookRows: [][]driver.Value{
				[]driver.Value{1},
			},
			ExpectedError: nil,
		},
		testCase{
			Name:       "borrow unknown book",
			BorrowerId: 1,
			ExpectedBorrowerRows: [][]driver.Value{
				[]driver.Value{1},
			},
			BookId:           2,
			ExpectedBookRows: [][]driver.Value{},
			ExpectedError:    books.ErrUnknownBook,
		},
		testCase{
			Name:                 "borrow with unknown borrower",
			BorrowerId:           1,
			ExpectedBorrowerRows: [][]driver.Value{},
			BookId:               2,
			ExpectedBookRows: [][]driver.Value{
				[]driver.Value{2},
			},
			ExpectedError: books.ErrUnknownBorrower,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %s", err)
			}
			defer db.Close()

			mock.ExpectBegin()

			borrowerRows := sqlmock.NewRows([]string{"1"})
			for _, r := range tc.ExpectedBorrowerRows {
				borrowerRows.AddRow(r...)
			}
			mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM BORROWERS WHERE `BorrowerId`=?")).
				WithArgs(tc.BorrowerId).
				WillReturnRows(borrowerRows)

			if tc.Name != "borrow with unknown borrower" {
				bookRows := sqlmock.NewRows([]string{"1"})
				for _, r := range tc.ExpectedBookRows {
					bookRows.AddRow(r...)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM BOOKS WHERE `BookId`=?")).
					WithArgs(tc.BookId).
					WillReturnRows(bookRows)

				if tc.Name != "borrow unknown book" {
					mock.ExpectExec(regexp.QuoteMeta("UPDATE BOOKS SET `BorrowedBy`=? WHERE `BookId`=?")).
						WithArgs(tc.BorrowerId, tc.BookId).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
			}

			err = books.BorrowBook(db, t.Context(), tc.BookId, tc.BorrowerId)
			fmt.Println(errors.Is(err, tc.ExpectedError))
			if err != nil && !errors.Is(err, tc.ExpectedError) {
				t.Error(err)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
