package pages

import (
	"bookshelf/internal/dto"
	"bookshelf/internal/render"
)

type quotesPageData struct {
	Quotes []dto.Quote
}

func RenderQuotesPage(renderer *render.TemplateRenderer, bookshelf *dto.Bookshelf) error {
	data := quotesPageData{
		Quotes: bookshelf.BookQuotes(),
	}

	return renderer.RenderToFile("quotes", data, "quotes")
}
