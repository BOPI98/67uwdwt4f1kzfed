package books

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	ErrUnknownBook     = errors.New("Unknown book!")
	ErrUnknownBorrower = errors.New("Unknown borrower!")

	bookTracer = otel.Tracer("book-tracer")
	bookMeter  = otel.Meter("book-meter")
	bookCount  metric.Int64Counter
)

func init() {
	var err error
	bookCount, err = bookMeter.Int64Counter("books.count")
	if err != nil {
		log.Fatal(err)
	}
}

type BookListItem struct {
	BookId     int       `json:"bookId"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	PageCount  int       `json:"pageCount"`
	BorrowedBy *int      `json:"borrowedBy"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Book struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	PageCount  int    `json:"pageCount"`
	BorrowedBy *int   `json:"borrowedBy"`
}

func ListBooks(db *sql.DB, ctx context.Context) ([]BookListItem, error) {
	Result := []BookListItem{}
	rows, err := db.QueryContext(ctx, "SELECT `BookId`,`Title`,`Author`,`PageCount`,`BorrowedBy`,`CreatedAt` FROM BOOKS;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book BookListItem
		err = rows.Scan(&book.BookId, &book.Title, &book.Author, &book.PageCount, &book.BorrowedBy, &book.CreatedAt)
		if err != nil {
			return nil, err
		}
		Result = append(Result, book)
	}
	return Result, nil
}

func AddBooks(db *sql.DB, ctx context.Context, book Book) (int, error) {
	bookCount.Add(ctx, 1, metric.WithAttributes(attribute.String("book-title", book.Title)))

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

	res, err := tx.ExecContext(ctx, "INSERT INTO BOOKS(`Title`,`Author`,`PageCount`,`BorrowedBy`) VALUES( ?,?,?,? );", book.Title, book.Author, book.PageCount, book.BorrowedBy)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func BorrowBook(db *sql.DB, ctx context.Context, bookId, borrowerId int) error {
	ctx, span := bookTracer.Start(ctx, "borrow-book")
	defer span.End()

	span.SetAttributes(attribute.Int("borrowerId", borrowerId), attribute.Int("bookId", bookId))

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var temp int
	err = db.QueryRow("SELECT 1 FROM BORROWERS WHERE `BorrowerId`=?", borrowerId).Scan(&temp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUnknownBorrower
		}
		return err
	}

	err = db.QueryRow("SELECT 1 FROM BOOKS WHERE `BookId`=?", bookId).Scan(&temp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUnknownBook
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE BOOKS SET `BorrowedBy`=? WHERE `BookId`=?", borrowerId, bookId)
	if err != nil {
		return err
	}

	return nil
}
