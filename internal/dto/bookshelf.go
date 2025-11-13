package dto

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
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

func (b *Bookshelf) bookById() map[string]Book {
	bookById := make(map[string]Book, len(b.Books))
	for _, book := range b.Books {
		bookById[book.Id] = book
	}

	return bookById
}

func (b *Bookshelf) booksByStatus() map[string][]Book {
	booksByStatus := make(map[string][]Book)

	for _, book := range b.Books {
		booksByStatus[book.Status] = append(booksByStatus[book.Status], book)
	}

	return booksByStatus
}

func (b *Bookshelf) UpcomingBooks(limit int) (map[string][]Book, bool) {
	hasUpcomingBooks := false
	booksByStatus := b.booksByStatus()
	delete(booksByStatus, StatusFinished) // Finished books are not part of the upcoming books

	for _, books := range booksByStatus {
		if len(books) > 0 {
			hasUpcomingBooks = true
			break
		}
	}

	sortBooksByRank(booksByStatus[StatusWishlisted])

	if limit <= 0 {
		return booksByStatus, hasUpcomingBooks
	}

	total := 0
	upcomingBooks := make(map[string][]Book, len(booksByStatus))

	// Fixed order to ensure deterministic limiting.
	for _, status := range []string{StatusReading, StatusToRead, StatusWishlisted} {
		books := booksByStatus[status]
		if total >= limit {
			break
		}

		remaining := limit - total
		if len(books) > remaining {
			upcomingBooks[status] = books[:remaining]
			total += remaining
		} else {
			upcomingBooks[status] = books
			total += len(books)
		}
	}

	return upcomingBooks, hasUpcomingBooks
}

func (b *Bookshelf) BookshelvedBooks() map[string][]Book {
	bookshelvedBooks := b.booksByStatus()
	delete(bookshelvedBooks, StatusWishlisted) // Wishlisted books are not considered as bookshelved

	for _, books := range bookshelvedBooks {
		sortBooksAlphabetically(books)
	}

	return bookshelvedBooks
}

func (b *Bookshelf) BookCollections() []ResolvedCollection {
	bookByID := b.bookById()
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
	wishlistedBooks := b.booksByStatus()[StatusWishlisted]
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
	genreCount := make(map[string]int)
	languageCount := make(map[string]int)
	statusCount := make(map[string]int)

	currentYear := time.Now().Year()

	for _, book := range b.Books {
		statusCount[book.Status]++

		// Wishlisted books are excluded
		if book.Status == StatusWishlisted {
			continue
		}

		// Total number of books in library
		stats.TotalBooks++

		// Total number of books finished
		if book.Status == StatusFinished {
			stats.BooksFinished++

			// Check if the book was finished this year
			if book.Progress.DateFinished != "" {
				finishedYear, err := getYearFromDate(book.Progress.DateFinished)
				if err == nil && finishedYear == currentYear {
					stats.BooksFinishedThisYear++
				}
			}
		}

		// Total number of pages read
		stats.PagesRead += book.Progress.PagesRead

		// Calculate pages read this year
		if book.Progress.DateStarted != "" && book.Progress.PagesRead > 0 {
			startYear, err := getYearFromDate(book.Progress.DateStarted)
			if err == nil && startYear == currentYear {
				stats.PagesReadThisYear += book.Progress.PagesRead
			}
		}

		// Total number of pages
		if book.Pages > 0 {
			totalPages += book.Pages
		}

		// Total rating
		if book.Rating > 0 {
			totalRating += book.Rating
			ratedBooks++
		}

		// Books by genre
		if book.Genre != "" {
			genreCount[book.Genre]++
		}

		// Books by language
		if book.Language != "" {
			languageCount[book.Language]++
		}
	}

	// Average pages
	if stats.TotalBooks > 0 {
		stats.AveragePages = float64(totalPages) / float64(stats.TotalBooks)
	}

	// Average rating
	if ratedBooks > 0 {
		stats.AverageRating = totalRating / float64(ratedBooks)
	}

	stats.TopGenres = topN(mapToStatCountSlice(genreCount), 3)
	stats.BooksByLanguage = mapToStatCountSlice(languageCount)
	stats.BooksByStatus = mapToStatCountSlice(statusCount)

	return stats
}

func sortBooksAlphabetically(books []Book) {
	sort.SliceStable(books, func(i, j int) bool {
		return books[i].Title < books[j].Title
	})
}

func sortBooksByRank(books []Book) {
	sort.SliceStable(books, func(i, j int) bool {
		ri, rj := books[i].Rank, books[j].Rank

		if ri == 0 && rj != 0 {
			return false
		}
		if ri != 0 && rj == 0 {
			return true
		}
		return ri < rj
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

func getYearFromDate(date string) (int, error) {
	// Try to parse the date in "yyyy-mm-dd" format first
	parsedDate, err := time.Parse("2006-01-02", date)
	if err == nil {
		return parsedDate.Year(), nil
	}

	// If the full year parse fails, try just the year "yyyy" format
	parsedDate, err = time.Parse("2006", date)
	if err != nil {
		return 0, err
	}

	return parsedDate.Year(), nil
}
