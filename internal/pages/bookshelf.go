package pages

import (
	"bookshelf/internal/dto"
	"bookshelf/internal/render"
)

type bookshelfPageData struct {
	Books map[string][]dto.Book
}

func RenderBookshelfPage(renderer *render.TemplateRenderer, bookshelf *dto.Bookshelf) error {
	data := bookshelfPageData{
		Books: bookshelf.ShelvedBooks(),
	}

	return renderer.RenderToFile("bookshelf", data, "bookshelf")
}
