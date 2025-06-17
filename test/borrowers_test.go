package test

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/BOPI98/67uwdwt4f1kzfed/borrowers"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateBorrower(t *testing.T) {
	type testCase struct {
		Name         string
		Borrower     borrowers.Borrower
		LastInsertId int
		RowsAffected int
		BorrowerId   int
	}
	cases := []testCase{
		testCase{
			Name: "Create a borrower",
			Borrower: borrowers.Borrower{
				Fullname:    "Borrower Name",
				Email:       "test@a.com",
				PhoneNumber: "",
				Age:         0,
			},
			LastInsertId: 4,
			RowsAffected: 1,
			BorrowerId:   4,
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
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO BORROWERS(`Fullname`,`Email`,`PhoneNumber`,`Age`) VALUES(?,?,?,?);")).
				WithArgs(tc.Borrower.Fullname, tc.Borrower.Email, tc.Borrower.PhoneNumber, tc.Borrower.Age).
				WillReturnResult(sqlmock.NewResult(int64(tc.LastInsertId), int64(tc.RowsAffected)))
			mock.ExpectCommit()

			borrowerId, err := borrowers.CreateBorrower(db, t.Context(), tc.Borrower)
			if err != nil {
				t.Fatal(err)
			}

			if borrowerId != tc.BorrowerId {
				t.Errorf("Expected borrowerId: %d, but got %d", tc.BorrowerId, borrowerId)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGetBorrower(t *testing.T) {
	type testCase struct {
		Name             string
		BorrowerId       int
		ExpectedRowsData [][]driver.Value
		ExpectedResult   *borrowers.Borrower
		ExpectedError    error
	}
	cases := []testCase{
		testCase{
			Name:       "Borrower found",
			BorrowerId: 2,
			ExpectedRowsData: [][]driver.Value{
				{"Borrower Name", "test@example.com", "+336556849", 22},
			},
			ExpectedResult: &borrowers.Borrower{
				Fullname:    "Borrower Name",
				Email:       "test@example.com",
				PhoneNumber: "+336556849",
				Age:         22,
			},
			ExpectedError: nil,
		},
		testCase{
			Name:             "Borrower not found",
			BorrowerId:       2,
			ExpectedRowsData: [][]driver.Value{},
			ExpectedResult:   nil,
			ExpectedError:    borrowers.ErrUnknownBorrower,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %s", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"Fullname", "Email", "PhoneNumber", "Age"})
			for _, r := range tc.ExpectedRowsData {
				rows.AddRow(r...)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT `Fullname`,`Email`,`PhoneNumber`,`Age` FROM BORROWERS WHERE `BorrowerId`=?")).
				WithArgs(tc.BorrowerId).
				WillReturnRows(rows).WillReturnError(tc.ExpectedError)

			result, err := borrowers.GetBorrower(db, t.Context(), tc.BorrowerId)
			if err != nil && !errors.Is(err, tc.ExpectedError) {
				t.Fatal(err)
			}

			if ((result == nil) != (tc.ExpectedResult == nil)) || (result != nil && *result != *tc.ExpectedResult) {
				t.Errorf("Expected %v result, but recieved %v", tc.ExpectedResult, result)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBorrowedBooks(t *testing.T) {
	type testCase struct {
		Name             string
		BorrowerId       int
		ExpectedRowsData [][]driver.Value
		ExpectedError    error
	}
	cases := []testCase{
		testCase{
			Name:       "Get borrowed books",
			BorrowerId: 2,
			ExpectedRowsData: [][]driver.Value{
				{1, "Test title", "Test author", 450, 2, time.Now()},
			},
			ExpectedError: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %s", err)
			}
			defer db.Close()

			rows := mock.NewRows([]string{"BookId", "Title", "Author", "PageCount", "BorrowedBy", "CreatedAt"})
			for _, r := range tc.ExpectedRowsData {
				rows.AddRow(r...)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT `BookId`,`Title`,`Author`,`PageCount`,`BorrowedBy`,`CreatedAt` FROM BOOKS WHERE `BorrowedBy`=?;")).
				WithArgs(tc.BorrowerId).
				WillReturnRows(rows).
				WillReturnError(tc.ExpectedError)

			books, err := borrowers.BorrowedBooks(db, t.Context(), tc.BorrowerId)
			if err != nil {
				t.Fatal(err)
			}

			if len(books) != len(tc.ExpectedRowsData) {
				t.Error(fmt.Errorf("Result list doesnt match the expected list length! result:%d expected:%d", len(books), len(tc.ExpectedRowsData)))
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
