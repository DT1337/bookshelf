package pages

import (
	"bookshelf/internal/dto"
	"bookshelf/internal/render"
)

type indexPageData struct {
	HasUpcomingBooks      bool
	UpcomingBooks         map[string][]dto.Book
	HasLatestReviewedBook bool
	LatestReviewedBook    dto.Book
	HasFavoriteQuote      bool
	FavoriteQuote         dto.Quote
	Stats                 dto.Stats
}

func RenderIndexPage(renderer *render.TemplateRenderer, bookshelf *dto.Bookshelf) error {
	upcomingBooks, hasUpcomingBooks := bookshelf.UpcomingBooks(3)
	latestReviewedBook, hasLatestReviewedBook := bookshelf.LatestReviewedBook()
	favoriteQuote, hasFavoriteQuote := bookshelf.FavoriteQuote()

	data := indexPageData{
		HasUpcomingBooks:      hasUpcomingBooks,
		UpcomingBooks:         upcomingBooks,
		HasLatestReviewedBook: hasLatestReviewedBook,
		LatestReviewedBook:    latestReviewedBook,
		HasFavoriteQuote:      hasFavoriteQuote,
		FavoriteQuote:         favoriteQuote,
		Stats:                 bookshelf.Stats(),
	}

	return renderer.RenderToFile("index", data, "index")
}
