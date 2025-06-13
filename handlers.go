package main

import (
	"encoding/json"
	"io"
	"main/books"
	"main/borrowers"
	"net/http"
	"strconv"
)

func Handle_Books(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		Result, err := books.ListBooks(db, ctx)
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
		bookId, err := books.AddBooks(db, ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		respondWithJson(w, struct {
			BookId int `json:"bookId"`
		}{bookId})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

type BookBorrowRequest struct {
	BorrowerId int `json:"borrowerId"`
}

func Handle_Books_Id_Borrow(w http.ResponseWriter, r *http.Request) {
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
	err = books.BorrowBook(db, ctx, bookId, req.BorrowerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondWithJson(w, struct{}{})
}

func Handle_Borrowers(w http.ResponseWriter, r *http.Request) {
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
	borrowerId, err := borrowers.CreateBorrower(db, ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondWithJson(w, struct {
		BorrowerId int `json:"borrowerId"`
	}{borrowerId})
}

func Handle_Borrowers_Id(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	borrowerId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid borrowerId in request URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	Result, err := borrowers.GetBorrower(db, ctx, borrowerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondWithJson(w, Result)
}

func Handle_Borrowers_Id_BorrowedBooks(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	borrowerId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid borrowerId in request URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	books, err := borrowers.BorrowedBooks(db, ctx, borrowerId)
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
