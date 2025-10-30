package dto

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const (
	StatusFinished   = "finished"
	StatusReading    = "reading"
	StatusToRead     = "to read"
	StatusWishlisted = "wishlisted"
)

func LoadBookshelfFromFile(path string) (*Bookshelf, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading JSON file: %w", err)
	}

	var bookshelf Bookshelf
	if err := json.Unmarshal(data, &bookshelf); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return &bookshelf, nil
}

func (b *Bookshelf) BooksByStatus() map[string][]Book {
	booksByStatus := make(map[string][]Book)
	books := append([]Book(nil), b.Books...) // clone to preserve original order

	sortBooksAlphabetically(books)
	for _, book := range books {
		booksByStatus[book.Status] = append(booksByStatus[book.Status], book)
	}

	return booksByStatus
}

func (b *Bookshelf) UpcomingBooks(limit int) map[string][]Book {
	upcoming := map[string][]Book{
		StatusReading: {},
		StatusToRead:  {},
	}

	for _, book := range b.Books {
		if book.Status == StatusReading || book.Status == StatusToRead {
			upcoming[book.Status] = append(upcoming[book.Status], book)
		}
	}

	if limit > 0 {
		for key, books := range upcoming {
			if len(books) > limit {
				upcoming[key] = books[:limit]
			}
		}
	}

	return upcoming
}

func (b *Bookshelf) BookshelfedBooks() map[string][]Book {
	bookshelfedBooks := b.BooksByStatus()
	delete(bookshelfedBooks, StatusWishlisted) // Wishlisted books are not part of the bookshelf

	return bookshelfedBooks
}

func (b *Bookshelf) BookCollections() []ResolvedCollection {
	bookByID := make(map[string]Book, len(b.Books))
	for _, book := range b.Books {
		bookByID[book.Id] = book
	}

	resolved := make([]ResolvedCollection, 0, len(b.Collections))
	for _, c := range b.Collections {
		var books []Book
		for i, id := range c.Books {
			if book, ok := bookByID[id]; ok {
				book.Rank = i + 1
				books = append(books, book)
			}
		}

		resolved = append(resolved, ResolvedCollection{
			Name:        c.Name,
			Description: c.Description,
			Books:       books,
		})
	}

	return resolved
}

func (b *Bookshelf) WishlistedBooks() []Book {
	wishlistedBooks := b.BooksByStatus()[StatusWishlisted]
	sortBooksByRank(wishlistedBooks)

	return wishlistedBooks
}

func (b *Bookshelf) Stats() Stats {
	var stats Stats
	if len(b.Books) == 0 {
		return stats
	}

	var totalPages int
	var totalRating float64
	var ratedBooks int
	statusCount := make(map[string]int)
	languageCount := make(map[string]int)
	genreCount := make(map[string]int)

	for _, book := range b.Books {
		stats.TotalBooks++

		if book.Status == StatusFinished {
			stats.BooksFinished++
		}
		stats.PagesRead += book.Progress.PagesRead
		if book.Rating > 0 {
			totalRating += book.Rating
			ratedBooks++
		}
		if book.Pages > 0 {
			totalPages += book.Pages
		}

		statusCount[book.Status]++
		if book.Language != "" {
			languageCount[book.Language]++
		}
		if book.Genre != "" {
			genreCount[book.Genre]++
		}
	}

	if stats.TotalBooks > 0 {
		stats.AveragePages = float64(totalPages) / float64(stats.TotalBooks)
	}
	if ratedBooks > 0 {
		stats.AverageRating = totalRating / float64(ratedBooks)
	}

	stats.TopGenres = topN(mapToStatCountSlice(genreCount), 3)
	stats.BooksByStatus = mapToStatCountSlice(statusCount)
	stats.BooksByLanguage = mapToStatCountSlice(languageCount)

	return stats
}

func sortBooksAlphabetically(books []Book) {
	sort.SliceStable(books, func(i, j int) bool {
		return books[i].Title < books[j].Title
	})
}

func sortBooksByRank(books []Book) {
	sort.SliceStable(books, func(i, j int) bool {
		return books[i].Rank < books[j].Rank
	})
}

func mapToStatCountSlice(m map[string]int) []StatCount {
	stats := make([]StatCount, 0, len(m))
	for k, v := range m {
		stats = append(stats, StatCount{Value: k, Count: v})
	}
	sortByCount(stats)
	return stats
}

func sortByCount(s []StatCount) {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Count > s[j].Count
	})
}

func topN(s []StatCount, n int) []StatCount {
	if len(s) > n {
		return s[:n]
	}
	return s
}
