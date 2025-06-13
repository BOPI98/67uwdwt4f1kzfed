package borrowers

import (
	"context"
	"database/sql"
	"errors"
	"main/books"
)

type Borrower struct {
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Age         int    `json:"age"`
}

func CreateBorrower(db *sql.DB, ctx context.Context, borrower Borrower) (int, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	res, err := tx.ExecContext(ctx, "INSERT INTO BORROWERS(`Fullname`,`Email`,`PhoneNumber`,`Age`) VALUES(?,?,?,?);", borrower.Fullname, borrower.Email, borrower.PhoneNumber, borrower.Age)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetBorrower(db *sql.DB, ctx context.Context, borrowerId int) (*Borrower, error) {
	Result := &Borrower{}

	row := db.QueryRow("SELECT `Fullname`,`Email`,`PhoneNumber`,`Age` FROM BORROWERS WHERE `BorrowerId`=?", borrowerId)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, errors.New("Borrower not found!")
		}
		return nil, row.Err()
	}
	err := row.Scan(&Result.Fullname, &Result.Email, &Result.PhoneNumber, &Result.Age)
	if err != nil {
		return nil, err
	}
	return Result, nil
}

func BorrowedBooks(db *sql.DB, ctx context.Context, borrowerId int) ([]books.BookListItem, error) {
	Result := []books.BookListItem{}
	rows, err := db.QueryContext(ctx, "SELECT `BookId`,`Title`,`Author`,`PageCount`,`BorrowedBy`,`CreatedAt` FROM BOOKS WHERE `BorrowedBy`=?;", borrowerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book books.BookListItem
		err = rows.Scan(&book.BookId, &book.Title, &book.Author, &book.PageCount, &book.BorrowedBy, &book.CreatedAt)
		if err != nil {
			return nil, err
		}
		Result = append(Result, book)
	}
	return Result, nil
}
