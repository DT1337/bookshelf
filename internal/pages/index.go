package pages

import (
	"bookshelf/internal/dto"
	"bookshelf/internal/render"
)

type indexPageData struct {
	HasUpcomingBooks bool
	UpcomingBooks    map[string][]dto.Book
	Stats            dto.Stats
}

func RenderIndexPage(renderer *render.TemplateRenderer, bookshelf *dto.Bookshelf) error {
	upcomingBooks := bookshelf.UpcomingBooks(3)

	hasUpcomingBooks := false
	for _, books := range upcomingBooks {
		if len(books) > 0 {
			hasUpcomingBooks = true
			break
		}
	}

	data := indexPageData{
		HasUpcomingBooks: hasUpcomingBooks,
		UpcomingBooks:    upcomingBooks,
		Stats:            bookshelf.Stats(),
	}

	return renderer.RenderToFile("index", data)
}
