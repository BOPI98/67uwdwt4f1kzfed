package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/BOPI98/67uwdwt4f1kzfed/books"
	"github.com/BOPI98/67uwdwt4f1kzfed/borrowers"
	"github.com/BOPI98/67uwdwt4f1kzfed/server"
	_ "github.com/go-sql-driver/mysql"
)

func startServer(ctx context.Context) {
	db, err := sql.Open("mysql", "root:asd@tcp(localhost)/LIBRARY_MANAGEMENT?parseTime=true")
	if err != nil {
		log.Fatal("Connecting to the database failed:", err)
	}
	defer db.Close()
	log.Println("Database connected!")

	srv := server.NewServer("8889", db)
	srv.Run(ctx)
}

func TestLibraryManagement(t *testing.T) {
	go startServer(t.Context())

	client := http.Client{
		Timeout: time.Second * 10}
	//Create book
	inputBook := books.Book{
		Title:      "Lord of the rings",
		Author:     "Tolkien",
		PageCount:  2000,
		BorrowedBy: nil,
	}
	requestBodyBuf, err := json.Marshal(inputBook)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Sending request: %s %s", "POST", "http://localhost:8889/books")
	resp, err := client.Post("http://localhost:8889/books", "application/json", bytes.NewBuffer(requestBodyBuf))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Unexpected status:", resp.Status)
	}
	var bookIdResponse struct {
		BookId int `json:"bookId"`
	}
	err = json.Unmarshal(buf, &bookIdResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	t.Logf("Book Created! id:%v", bookIdResponse.BookId)

	//Create Borrower
	inputBorrower := borrowers.Borrower{
		Fullname:    "Borrower Name",
		Email:       "email@a.com",
		PhoneNumber: "+361849661321",
		Age:         30,
	}
	requestBodyBuf, err = json.Marshal(inputBorrower)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Sending request: %s %s", "POST", "http://localhost:8889/borrowers")
	resp, err = client.Post("http://localhost:8889/borrowers", "application/json", bytes.NewBuffer(requestBodyBuf))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Unexpected status:", resp.Status)
	}
	var borrowerIdResponse struct {
		BorrowerId int `json:"borrowerId"`
	}
	err = json.Unmarshal(buf, &borrowerIdResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	t.Logf("Borrower Created! id:%v", borrowerIdResponse.BorrowerId)

	//Borrow book
	borrowRequestBody := server.BookBorrowRequest{
		BorrowerId: borrowerIdResponse.BorrowerId,
	}
	borrowRequestBodyJson, err := json.Marshal(borrowRequestBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	t.Logf("Sending request: %s %s", "PATCH", "http://localhost:8889/books/"+strconv.Itoa(bookIdResponse.BookId)+"/borrow")
	request, err := http.NewRequest(http.MethodPatch, "http://localhost:8889/books/"+strconv.Itoa(bookIdResponse.BookId)+"/borrow", bytes.NewBuffer(borrowRequestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("Unexpected status:", resp.Status)
	}

	t.Log("Book borrowed successfully!")

	//Get borrower
	t.Logf("Sending request: %s %s", "GET", "http://localhost:8889/borrowers/"+strconv.Itoa(borrowerIdResponse.BorrowerId))
	resp, err = client.Get("http://localhost:8889/borrowers/" + strconv.Itoa(borrowerIdResponse.BorrowerId))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("Unexpected status:", resp.Status)
	}

	var borrower borrowers.Borrower
	err = json.Unmarshal(buf, &borrower)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	if borrower != inputBorrower {
		t.Error("The recieved borrower data does not match the inputed data! ")
	}

	//List books
	t.Logf("Sending request: %s %s", "GET", "http://localhost:8889/books")
	resp, err = client.Get("http://localhost:8889/books")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("Unexpected status:", resp.Status)
	}

	booklist := []books.BookListItem{}
	err = json.Unmarshal(buf, &booklist)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	t.Logf("Book list: %v", booklist)

	//Find new book borrowed by the new borrower
	found := false
	for _, book := range booklist {
		if book.BookId != bookIdResponse.BookId {
			continue
		}

		found = true
		if book.BorrowedBy == nil {
			t.Fatalf("The new book's borrower is not specified! Expected to be: %d", borrowerIdResponse.BorrowerId)
		}
		if *book.BorrowedBy != borrowerIdResponse.BorrowerId {
			t.Fatal("The new book's borrower does not match the new borrower!")
		}
		t.Log("The new book was found in the book list, with the new borrower!")
	}
	if !found {
		t.Error("The created book was not found in the book list!")
	}
}
