package dto

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
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
	booksByStatus := b.booksByStatus()
	delete(booksByStatus, StatusFinished) // Finished books are not part of the upcoming books
	hasUpcomingBooks := b.hasUpcomingBooks(booksByStatus)

	b.sortBooksByRank(booksByStatus[StatusWishlisted])

	if limit <= 0 {
		return booksByStatus, hasUpcomingBooks
	}

	upcomingBooks := b.getUpcomingBooksWithLimit(booksByStatus, limit)

	return upcomingBooks, hasUpcomingBooks
}

func (b *Bookshelf) hasUpcomingBooks(booksByStatus map[string][]Book) bool {
	for _, books := range booksByStatus {
		if len(books) > 0 {
			return true
		}
	}

	return false
}

func (b *Bookshelf) getUpcomingBooksWithLimit(booksByStatus map[string][]Book, limit int) map[string][]Book {
	upcomingBooks := make(map[string][]Book, len(booksByStatus))
	total := 0

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

	return upcomingBooks
}

func (b *Bookshelf) ShelvedBooks() map[string][]Book {
	shelvedBooks := b.booksByStatus()
	delete(shelvedBooks, StatusWishlisted) // Wishlisted books are not considered as shelved books

	for _, books := range shelvedBooks {
		b.sortBooksAlphabetically(books)
	}

	return shelvedBooks
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

func (b *Bookshelf) BookQuotes() []Quote {
	var quotes []Quote

	for _, book := range b.Books {
		for _, quote := range book.Quotes {
			quotes = append(quotes, Quote{
				Quote:     quote,
				Authors:   book.Authors,
				BookTitle: book.Title,
				Id:        book.Id,
			})
		}
	}

	// Sort quotes based on the hash of the quote text for a pseudo random (deterministic) order
	sort.SliceStable(quotes, func(i, j int) bool {
		hashI := md5.Sum([]byte(quotes[i].Quote))
		hashJ := md5.Sum([]byte(quotes[j].Quote))
		return hex.EncodeToString(hashI[:]) < hex.EncodeToString(hashJ[:])
	})

	return quotes
}

func (b *Bookshelf) WishlistedBooks() []Book {
	wishlistedBooks := b.booksByStatus()[StatusWishlisted]
	b.sortBooksByRank(wishlistedBooks)

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

		stats.TotalBooks++

		b.updateStatsForFinishedBook(&stats, book, currentYear)
		b.updateStatsForPages(&stats, book, currentYear)

		totalPages = b.updateTotalPages(totalPages, book)
		totalRating, ratedBooks = b.updateTotalRating(totalRating, ratedBooks, book)

		genreCount = b.updateGenreCount(genreCount, book)
		languageCount = b.updateLanguageCount(languageCount, book)
	}

	stats.AveragePages = b.calculateAverage(float64(totalPages), stats.TotalBooks)
	stats.AverageRating = b.calculateAverage(totalRating, ratedBooks)

	stats.TopGenres = b.topGenres(genreCount, 3)
	stats.BooksByLanguage = b.mapToStatCountSlice(languageCount)
	stats.BooksByStatus = b.mapToStatCountSlice(statusCount)

	return stats
}

func (b *Bookshelf) sortBooksAlphabetically(books []Book) {
	sort.SliceStable(books, func(i, j int) bool {
		return books[i].Title < books[j].Title
	})
}

func (b *Bookshelf) sortBooksByRank(books []Book) {
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

func (b *Bookshelf) sortByCount(s []StatCount) {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Count > s[j].Count
	})
}

func (b *Bookshelf) getYearFromDate(date string) int {
	// Try to parse the date in "yyyy-mm-dd" format first
	parsedDate, err := time.Parse("2006-01-02", date)
	if err == nil {
		return parsedDate.Year()
	}

	// If the full year parse fails, try just the year "yyyy" format
	parsedDate, err = time.Parse("2006", date)
	if err != nil {
		return 0
	}

	return parsedDate.Year()
}

func (b *Bookshelf) updateStatsForFinishedBook(stats *Stats, book Book, currentYear int) {
	if book.Status == StatusFinished {
		stats.BooksFinished++

		finishedYear := b.getYearFromDate(book.Progress.DateFinished)
		if finishedYear == currentYear {
			stats.BooksFinishedThisYear++
		}
	}
}

func (b *Bookshelf) updateStatsForPages(stats *Stats, book Book, currentYear int) {
	stats.PagesRead += b.calculatePagesRead(book)

	if book.Status != StatusToRead {
		startYear := b.getYearFromDate(book.Progress.DateStarted)
		if startYear == currentYear {
			stats.PagesReadThisYear += b.calculatePagesRead(book)
		}
	}
}

func (b *Bookshelf) calculatePagesRead(book Book) int {
	if book.Status == StatusFinished {
		return book.Pages
	}

	return book.Progress.PagesRead
}

func (b *Bookshelf) updateTotalPages(totalPages int, book Book) int {
	if book.Pages > 0 {
		return totalPages + book.Pages
	}

	return totalPages
}

func (b *Bookshelf) updateTotalRating(totalRating float64, ratedBooks int, book Book) (float64, int) {
	if book.Rating > 0 {
		totalRating += book.Rating
		ratedBooks++
	}

	return totalRating, ratedBooks
}

func (b *Bookshelf) updateGenreCount(genreCount map[string]int, book Book) map[string]int {
	if book.Genre != "" {
		genreCount[book.Genre]++
	}

	return genreCount
}

func (b *Bookshelf) updateLanguageCount(languageCount map[string]int, book Book) map[string]int {
	if book.Language != "" {
		languageCount[book.Language]++
	}

	return languageCount
}

func (b *Bookshelf) calculateAverage(total float64, fraction int) float64 {
	if fraction > 0 {
		return math.Round(total/float64(fraction)*100) / 100
	}

	return 0
}

func (b *Bookshelf) topGenres(genreCount map[string]int, n int) []StatCount {
	statCount := b.mapToStatCountSlice(genreCount)

	if len(statCount) > n {
		return statCount[:n]
	}

	return statCount
}

func (b *Bookshelf) mapToStatCountSlice(m map[string]int) []StatCount {
	stats := make([]StatCount, 0, len(m))
	for k, v := range m {
		stats = append(stats, StatCount{Value: k, Count: v})
	}
	b.sortByCount(stats)

	return stats
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
