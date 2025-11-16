package dto

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func createTestBookshelf() *Bookshelf {
	return &Bookshelf{
		Books: []Book{
			{
				Id:        "book-1",
				Title:     "Book One",
				Language:  "en",
				Genre:     "fiction",
				Pages:     300,
				Status:    StatusFinished,
				Rank:      2,
				Rating:    4.5,
				Progress:  Progress{DateStarted: "2025-11-08", DateFinished: "2025-11-14"},
				DateAdded: "2025-11-01",
			},
			{
				Id:        "book-2",
				Title:     "Book Two",
				Language:  "de",
				Genre:     "non-fiction",
				Pages:     150,
				Status:    StatusFinished,
				Rank:      5,
				Rating:    3.8,
				Progress:  Progress{DateStarted: "2025-11-01", DateFinished: "2025-11-07"},
				DateAdded: "2025-11-01",
			},
			{
				Id:        "book-3",
				Title:     "Book Three",
				Language:  "en",
				Genre:     "fiction",
				Pages:     200,
				Status:    StatusReading,
				Rank:      3,
				Rating:    0,
				Progress:  Progress{DateStarted: "2025-11-15", DateFinished: "", PagesRead: 50},
				DateAdded: "2025-11-01",
			},
			{
				Id:        "book-4",
				Title:     "Book Four",
				Language:  "en",
				Genre:     "science-fiction",
				Pages:     350,
				Status:    StatusToRead,
				Rank:      1,
				Rating:    0,
				Progress:  Progress{DateStarted: "", DateFinished: "", PagesRead: 0},
				DateAdded: "2025-11-01",
			},
			{
				Id:        "book-5",
				Title:     "Book Five",
				Language:  "en",
				Genre:     "science-fiction",
				Pages:     400,
				Status:    StatusWishlisted,
				Rank:      6,
				Rating:    0,
				Progress:  Progress{DateStarted: "", DateFinished: "", PagesRead: 0},
				DateAdded: "2025-11-01",
			},
			{
				Id:        "book-6",
				Title:     "Book Six",
				Language:  "de",
				Genre:     "literature",
				Pages:     100,
				Status:    StatusFinished,
				Rank:      4,
				Rating:    4.0,
				Progress:  Progress{DateStarted: "2024", DateFinished: "2024"},
				DateAdded: "2025-11-01",
			},
		},
	}
}

func TestLoadBookshelfFromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "data-*.json")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	bookshelf := createTestBookshelf()

	data, err := json.Marshal(bookshelf)
	if err != nil {
		t.Fatalf("could not marshal bookshelf data: %v", err)
	}

	if _, err := tmpFile.Write(data); err != nil {
		t.Fatalf("could not write to temp file: %v", err)
	}

	loadedBookshelf, err := LoadBookshelfFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("could not load bookshelf from file: %v", err)
	}

	if len(loadedBookshelf.Books) != len(bookshelf.Books) {
		t.Errorf("expected %d books, got %d", len(bookshelf.Books), len(loadedBookshelf.Books))
	}
}

func TestUpcomingBooks_LimitZero(t *testing.T) {
	bookshelf := createTestBookshelf()

	upcomingBooks, _ := bookshelf.UpcomingBooks(0)

	if len(upcomingBooks) != 3 {
		t.Errorf("expected 3 upcoming books, got %d", len(upcomingBooks))
	}
}

func TestUpcomingBooks_Limited(t *testing.T) {
	bookshelf := createTestBookshelf()

	upcomingBooks, _ := bookshelf.UpcomingBooks(2)

	if len(upcomingBooks[StatusFinished]) != 0 {
		t.Errorf("finished books should not be included in upcoming books")
	}

	if len(upcomingBooks[StatusReading])+len(upcomingBooks[StatusToRead])+len(upcomingBooks[StatusWishlisted]) != 2 {
		t.Errorf("expected 2 upcoming books, got %d", len(upcomingBooks[StatusReading])+len(upcomingBooks[StatusToRead]))
	}
}

func TestShelvedBooks(t *testing.T) {
	bookshelf := createTestBookshelf()

	shelvedBooks := bookshelf.ShelvedBooks()

	if len(shelvedBooks[StatusWishlisted]) > 0 {
		t.Error("wishlisted books should not be included in shelved books")
	}

	numberOfShelvedBooks := len(shelvedBooks[StatusFinished]) + len(shelvedBooks[StatusReading]) + len(shelvedBooks[StatusToRead])
	if numberOfShelvedBooks != 5 {
		t.Errorf("expected 5 shelved books, got %d", numberOfShelvedBooks)
	}
}

func TestUpdateStatsForFinishedBook(t *testing.T) {
	bookshelf := createTestBookshelf()

	stats := Stats{}
	currentYear := time.Now().Year()

	for _, book := range bookshelf.Books {
		bookshelf.updateStatsForFinishedBook(&stats, book, currentYear)
	}

	if stats.BooksFinished != 3 {
		t.Errorf("Expected 3 finished books, got %d", stats.BooksFinished)
	}

	if stats.BooksFinishedThisYear != 2 {
		t.Errorf("Expected 2 finished books this year, got %d", stats.BooksFinishedThisYear)
	}
}

func TestUpdateStatsForPages(t *testing.T) {
	bookshelf := createTestBookshelf()

	stats := Stats{}
	currentYear := time.Now().Year()

	for _, book := range bookshelf.Books {
		bookshelf.updateStatsForPages(&stats, book, currentYear)
	}

	if stats.PagesRead != 600 {
		t.Errorf("Expected total pages read to be 600, got %d", stats.PagesRead)
	}

	if stats.PagesReadThisYear != 500 {
		t.Errorf("Expected pages read this year to be 500, got %d", stats.PagesReadThisYear)
	}
}

func TestUpdateTotalPages(t *testing.T) {
	bookshelf := createTestBookshelf()

	totalPages := 0
	for _, book := range bookshelf.Books {
		totalPages = bookshelf.updateTotalPages(totalPages, book)
	}

	if totalPages != 1500 {
		t.Errorf("Expected total pages to be 1500, got %d", totalPages)
	}
}

func TestUpdateTotalRating(t *testing.T) {
	bookshelf := createTestBookshelf()

	totalRating := 0.0
	ratedBooks := 0
	for _, book := range bookshelf.Books {
		totalRating, ratedBooks = bookshelf.updateTotalRating(totalRating, ratedBooks, book)
	}

	if totalRating != 12.3 {
		t.Errorf("Expected total rating to be 12.3, got %.2f", totalRating)
	}

	if ratedBooks != 3 {
		t.Errorf("Expected 3 rated books, got %d", ratedBooks)
	}
}

func TestUpdateGenreCount(t *testing.T) {
	bookshelf := createTestBookshelf()

	genreCount := make(map[string]int)
	for _, book := range bookshelf.Books {
		genreCount = bookshelf.updateGenreCount(genreCount, book)
	}

	if genreCount["fiction"] != 2 {
		t.Errorf("Expected 2 books in the 'fiction' genre, got %d", genreCount["fiction"])
	}

	if genreCount["non-fiction"] != 1 {
		t.Errorf("Expected 1 book in the 'non-fiction' genre, got %d", genreCount["non-fiction"])
	}

	if genreCount["science-fiction"] != 2 {
		t.Errorf("Expected 2 book in the 'science-fiction' genre, got %d", genreCount["science-fiction"])
	}

	if genreCount["literature"] != 1 {
		t.Errorf("Expected 1 book in the 'literature' genre, got %d", genreCount["literature"])
	}
}

func TestUpdateLanguageCount(t *testing.T) {
	bookshelf := createTestBookshelf()

	languageCount := make(map[string]int)
	for _, book := range bookshelf.Books {
		languageCount = bookshelf.updateLanguageCount(languageCount, book)
	}

	if languageCount["en"] != 4 {
		t.Errorf("Expected 4 books in en, got %d", languageCount["en"])
	}

	if languageCount["de"] != 2 {
		t.Errorf("Expected 2 book in de, got %d", languageCount["de"])
	}
}

func TestCalculateAverage(t *testing.T) {
	bookshelf := createTestBookshelf()

	averageRating := bookshelf.calculateAverage(12.3, 3)
	if averageRating != 4.10 {
		t.Errorf("Expected average rating to be 4.10, got %.2f", averageRating)
	}

	averagePages := bookshelf.calculateAverage(1100, 5)
	if averagePages != 220 {
		t.Errorf("Expected average pages to be 220, got %f", averagePages)
	}
}

func TestTopGenres(t *testing.T) {
	bookshelf := createTestBookshelf()

	genreCount := map[string]int{
		"fiction":     2,
		"non-fiction": 1,
	}

	topGenres := bookshelf.topGenres(genreCount, 2)

	if len(topGenres) != 2 {
		t.Errorf("Expected 2 top genres, got %d", len(topGenres))
	}

	if topGenres[0].Value != "fiction" {
		t.Errorf("Expected top genre to be 'fiction', got %s", topGenres[0].Value)
	}

	if topGenres[0].Count != 2 {
		t.Errorf("Expected top genre count to be 2, got %d", topGenres[0].Count)
	}
}

func TestStats(t *testing.T) {
	bookshelf := createTestBookshelf()

	stats := bookshelf.Stats()

	if stats.TotalBooks != 5 {
		t.Errorf("Expected total books to be 5, got %d", stats.TotalBooks)
	}

	if stats.BooksFinished != 3 {
		t.Errorf("Expected 3 finished books, got %d", stats.BooksFinished)
	}

	if stats.BooksFinishedThisYear != 2 {
		t.Errorf("Expected 2 finished books this year, got %d", stats.BooksFinishedThisYear)
	}

	if stats.PagesRead != 600 {
		t.Errorf("Expected total pages read to be 600, got %d", stats.PagesRead)
	}

	if stats.PagesReadThisYear != 500 {
		t.Errorf("Expected total pages read to be 500, got %d", stats.PagesReadThisYear)
	}

	if stats.AverageRating != 4.10 {
		t.Errorf("Expected average rating to be 4.10, got %.2f", stats.AverageRating)
	}

	if stats.AveragePages != 220 {
		t.Errorf("Expected average pages to be 150, got %f", stats.AveragePages)
	}

	if len(stats.TopGenres) != 3 {
		t.Errorf("Expected 3 top genres, got %d", len(stats.TopGenres))
	}

	if stats.TopGenres[0].Value != "fiction" {
		t.Errorf("Expected top genre to be fiction, got %s", stats.TopGenres[0].Value)
	}
}

func TestSortBooksAlphabetically(t *testing.T) {
	bookshelf := createTestBookshelf()

	bookshelf.sortBooksAlphabetically(bookshelf.Books)

	if bookshelf.Books[0].Title != "Book Five" {
		t.Errorf("expected 'Book Five' to be first, got %s", bookshelf.Books[0].Title)
	}

	if bookshelf.Books[5].Title != "Book Two" {
		t.Errorf("expected 'Book Two' to be last, got %s", bookshelf.Books[2].Title)
	}
}

func TestSortBooksByRank(t *testing.T) {
	bookshelf := createTestBookshelf()
	bookshelf.sortBooksByRank(bookshelf.Books)

	if bookshelf.Books[0].Title != "Book Four" {
		t.Errorf("expected 'Book Four' to be first, got %s", bookshelf.Books[0].Title)
	}

	if bookshelf.Books[5].Title != "Book Five" {
		t.Errorf("expected 'Book Five' to be last, got %s", bookshelf.Books[5].Title)
	}
}

func TestGetYearFromDate(t *testing.T) {
	bookshelf := Bookshelf{}

	tests := []struct {
		date     string
		expected int
	}{
		{"2025-11-15", 2025},
		{"2025", 2025},
		{"invalid-date", 0},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			year := bookshelf.getYearFromDate(tt.date)
			if year != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, year)
			}
		})
	}
}
