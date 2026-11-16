// Package grpcsvc implements the KVStore gRPC service.
//
// STEP 4 — In Part I every RPC returns a DUMMY response: the goal is only to
// prove that the .proto compiles and that two processes can talk over gRPC.
// In Part II this type will apply operations to the state machine, and in
// Part III it will submit them to the replicated log.
package grpcsvc

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	booklibraryv1 "github.com/your-team/kvs/gen/booklibrary/v1"
)

// Server implements kvsv1.KVStoreServer.
//
// Embedding UnimplementedKVStoreServer makes every RPC not written yet return
// codes.Unimplemented (HTTP 501 after translation): when you add RPCs for your
// own service, the project keeps compiling while you work through them.
type Server struct {
	booklibraryv1.UnimplementedBooklibraryServer

	NodeID string // shown in dummy values, so you can see WHICH server answered
}

func New(nodeID string) *Server { return &Server{NodeID: nodeID} }

func (s *Server) dummyBooks() []*booklibraryv1.Book {
	return []*booklibraryv1.Book{
		{
			Booktitle:       "Dummy 1",
			Author:          fmt.Sprintf("sara nasra (dummy from %s)", s.NodeID),
			Year:            2019,
			TotalCopies:     10,
			AvailableCopies: 7,
		},
		{
			Booktitle:       "Dummy 2",
			Author:          fmt.Sprintf("suheir suleiman (dummy from %s)", s.NodeID),
			Year:            2017,
			TotalCopies:     8,
			AvailableCopies: 8,
		},
		{
			Booktitle:       "Dummy 3",
			Author:          fmt.Sprintf("Dummy author from %s", s.NodeID),
			Year:            2018,
			TotalCopies:     5,
			AvailableCopies: 0,
		},
	}
}

// ---------------------------------------------------------------------------
// The example service (a key-value store). Replace it with your own service.
// ---------------------------------------------------------------------------

// GetAllBooks is server streaming , one GetAllBooksResponse per book, in sort order.
func (s *Server) GetAllBooks(req *booklibraryv1.GetAllBooksRequest, stream booklibraryv1.Booklibrary_GetAllBooksServer) error {
	if req.GetLimit() < 0 {
		return status.Error(codes.InvalidArgument, "limit must be postive")
	}
	var books []*booklibraryv1.Book

	for _, book := range s.dummyBooks() {
		if req.Author != nil && !strings.Contains(strings.ToLower(book.GetAuthor()), strings.ToLower(req.GetAuthor())) {
			continue
		}
		if req.Year != nil && book.GetYear() != req.GetYear() {
			continue
		}
		books = append(books, book)
	}

	// Sorting by book title.
	for i := 0; i < len(books); i++ {
		for j := i + 1; j < len(books); j++ {
			if books[i].GetBooktitle() > books[j].GetBooktitle() {
				books[i], books[j] = books[j], books[i]
			}
		}
	}

	if n := int(req.GetLimit()); n > 0 && n < len(books) {
		books = books[:n]
	}

	for _, book := range books {
		if err := stream.Send(&booklibraryv1.GetAllBooksResponse{Book: book}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) GetBook(_ context.Context, req *booklibraryv1.GetBookRequest) (*booklibraryv1.GetBookResponse, error) {
	for _, book := range s.dummyBooks() {
		if book.GetBooktitle() == req.GetBooktitle() {
			return &booklibraryv1.GetBookResponse{Book: book}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetBooktitle())
}
func (s *Server) AddBook(_ context.Context, req *booklibraryv1.AddBookRequest) (*booklibraryv1.AddBookResponse, error) {
	if req.GetBooktitle() == "" || req.GetAuthor() == "" {
		return nil, status.Error(codes.InvalidArgument,
			"book title and author are required")
	}
	if req.GetCopies() <= 0 {
		return nil, status.Error(codes.InvalidArgument,
			"the book should have copies ")
	}
	for _, book := range s.dummyBooks() {
		if book.GetBooktitle() == req.GetBooktitle() {
			return nil, status.Errorf(codes.AlreadyExists, "book %q already exists", req.GetBooktitle())
		}
	}
	book := &booklibraryv1.Book{
		Booktitle:       req.GetBooktitle(),
		Author:          fmt.Sprintf("%s (dummy from %s)", req.GetAuthor(), s.NodeID),
		Year:            req.GetYear(),
		TotalCopies:     req.GetCopies(),
		AvailableCopies: req.GetCopies(),
	}

	return &booklibraryv1.AddBookResponse{Book: book}, nil
}

func (s *Server) AddCopies(_ context.Context, req *booklibraryv1.AddCopiesRequest) (*booklibraryv1.AddCopiesResponse, error) {
	if req.GetCopies() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "copies must be positive")
	}
	for _, book := range s.dummyBooks() {
		if book.GetBooktitle() == req.GetBooktitle() {
			book.TotalCopies += req.GetCopies()
			book.AvailableCopies += req.GetCopies()

			return &booklibraryv1.AddCopiesResponse{Book: book}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetBooktitle())
}

func (s *Server) RemoveBook(_ context.Context, req *booklibraryv1.RemoveBookRequest) (*booklibraryv1.RemoveBookResponse, error) {
	for _, book := range s.dummyBooks() {
		if book.GetBooktitle() == req.GetBooktitle() {
			if book.GetTotalCopies() != book.GetAvailableCopies() {
				return nil, status.Error(codes.FailedPrecondition, "can't remove a book with borrowed copies")
			}
			return &booklibraryv1.RemoveBookResponse{}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetBooktitle())
}

func (s *Server) BuyBook(_ context.Context, req *booklibraryv1.BuyBookRequest) (*booklibraryv1.BuyBookResponse, error) {
	for _, book := range s.dummyBooks() {
		if book.GetBooktitle() == req.GetBooktitle() {

			if book.GetAvailableCopies() == 0 {
				return nil, status.Error(codes.FailedPrecondition, "no available copies")
			}
			book.TotalCopies--
			book.AvailableCopies--
			return &booklibraryv1.BuyBookResponse{Book: book}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetBooktitle())
}

func (s *Server) BorrowBook(_ context.Context, req *booklibraryv1.BorrowBookRequest) (*booklibraryv1.BorrowBookResponse, error) {
	if req.GetUserid() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	for _, book := range s.dummyBooks() {
		if book.GetBooktitle() == req.GetBooktitle() {
			if book.GetAvailableCopies() == 0 {
				return nil, status.Error(codes.FailedPrecondition, "no available copies")
			}
			if req.GetUserid() == "user1" && book.GetTotalCopies() > book.GetAvailableCopies() {
				return nil, status.Error(codes.FailedPrecondition, "user already borrowed this book")
			}
			book.AvailableCopies--
			return &booklibraryv1.BorrowBookResponse{Book: book}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetBooktitle())
}

func (s *Server) ReturnBook(_ context.Context, req *booklibraryv1.ReturnBookRequest) (*booklibraryv1.ReturnBookResponse, error) {

	if req.GetUserid() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	for _, book := range s.dummyBooks() {
		if book.GetBooktitle() == req.GetBooktitle() {
			if req.GetUserid() != "user1" || book.GetTotalCopies() == book.GetAvailableCopies() {
				return nil, status.Error(codes.FailedPrecondition, "user did not borrow this book")
			}
			book.AvailableCopies++
			return &booklibraryv1.ReturnBookResponse{Book: book}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetBooktitle())
}
