package main

import (
	"context"
	"fmt"
	"net/http/httptest"
	"sync"

	"ogen-task/api"
)

type bookStore struct {
	mu     sync.Mutex
	nextID int64
	books  map[int64]api.Book
}

func newBookStore() *bookStore {
	return &bookStore{
		nextID: 1,
		books:  make(map[int64]api.Book),
	}
}

func (s *bookStore) ListBooks(ctx context.Context) ([]api.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	books := make([]api.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books, nil
}

func (s *bookStore) CreateBook(ctx context.Context, req *api.NewBook) (*api.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book := api.Book{
		ID:     s.nextID,
		Title:  req.Title,
		Author: req.Author,
	}
	s.books[book.ID] = book
	s.nextID++
	return &book, nil
}

func (s *bookStore) GetBook(ctx context.Context, params api.GetBookParams) (*api.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[params.ID]
	if !ok {
		return nil, fmt.Errorf("book %d not found", params.ID)
	}
	return &book, nil
}

func (s *bookStore) UpdateBook(ctx context.Context, req *api.NewBook, params api.UpdateBookParams) (*api.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[params.ID]; !ok {
		return nil, fmt.Errorf("book %d not found", params.ID)
	}
	book := api.Book{
		ID:     params.ID,
		Title:  req.Title,
		Author: req.Author,
	}
	s.books[book.ID] = book
	return &book, nil
}

func (s *bookStore) DeleteBook(ctx context.Context, params api.DeleteBookParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[params.ID]; !ok {
		return fmt.Errorf("book %d not found", params.ID)
	}
	delete(s.books, params.ID)
	return nil
}

func main() {
	ctx := context.Background()

	handler, err := api.NewServer(newBookStore())
	if err != nil {
		panic(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := api.NewClient(server.URL)
	if err != nil {
		panic(err)
	}

	created, err := client.CreateBook(ctx, &api.NewBook{
		Title:  "The Go Programming Language",
		Author: "Alan A. A. Donovan, Brian W. Kernighan",
	})
	if err != nil {
		panic(err)
	}

	updated, err := client.UpdateBook(ctx, &api.NewBook{
		Title:  "The Go Programming Language",
		Author: "Donovan and Kernighan",
	}, api.UpdateBookParams{ID: created.ID})
	if err != nil {
		panic(err)
	}

	got, err := client.GetBook(ctx, api.GetBookParams{ID: updated.ID})
	if err != nil {
		panic(err)
	}

	books, err := client.ListBooks(ctx)
	if err != nil {
		panic(err)
	}

	if err := client.DeleteBook(ctx, api.DeleteBookParams{ID: got.ID}); err != nil {
		panic(err)
	}

	fmt.Printf("created=%d got=%q books_before_delete=%d\n", created.ID, got.Title, len(books))
}
