package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/BOPI98/67uwdwt4f1kzfed/books"
	"github.com/BOPI98/67uwdwt4f1kzfed/borrowers"
)

type LibraryHandler struct {
	db *sql.DB
}

func NewLibraryHandler(db *sql.DB) LibraryHandler {
	return LibraryHandler{
		db: db,
	}
}

func (this *LibraryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/books", this.Handle_Books)
	mux.HandleFunc("/books/{id}/borrow", this.Handle_Books_Id_Borrow)
	mux.HandleFunc("/borrowers", this.Handle_Borrowers)
	mux.HandleFunc("/borrowers/{id}", this.Handle_Borrowers_Id)
	mux.HandleFunc("/borrowers/{id}/borrowedbooks", this.Handle_Borrowers_Id_BorrowedBooks)
}

func (this *LibraryHandler) Handle_Books(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		Result, err := books.ListBooks(this.db, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		respondWithJson(w, Result)
	case http.MethodPost:
		//unmarshal body
		buff, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		var req books.Book
		err = json.Unmarshal(buff, &req)
		if err != nil {
			http.Error(w, "Failed to unmarshal request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		bookId, err := books.AddBooks(this.db, ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		respondWithJson(w, struct {
			BookId int `json:"bookId"`
		}{bookId})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type BookBorrowRequest struct {
	BorrowerId int `json:"borrowerId"`
}

func (this *LibraryHandler) Handle_Books_Id_Borrow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	bookId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid bookId in request URL: "+err.Error(), http.StatusBadRequest)
		return
	}

	//unmarshal body
	buff, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	var req BookBorrowRequest
	err = json.Unmarshal(buff, &req)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = books.BorrowBook(this.db, ctx, bookId, req.BorrowerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondWithJson(w, struct{}{})
}

func (this *LibraryHandler) Handle_Borrowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//unmarshal body
	buff, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	var req borrowers.Borrower
	err = json.Unmarshal(buff, &req)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	borrowerId, err := borrowers.CreateBorrower(this.db, ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondWithJson(w, struct {
		BorrowerId int `json:"borrowerId"`
	}{borrowerId})
}

func (this *LibraryHandler) Handle_Borrowers_Id(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	borrowerId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid borrowerId in request URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	Result, err := borrowers.GetBorrower(this.db, ctx, borrowerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondWithJson(w, Result)
}

func (this *LibraryHandler) Handle_Borrowers_Id_BorrowedBooks(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	borrowerId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid borrowerId in request URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	books, err := borrowers.BorrowedBooks(this.db, ctx, borrowerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondWithJson(w, books)
}

func respondWithJson(w http.ResponseWriter, data any) {
	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
