package pages

import (
	"bookshelf/internal/dto"
	"bookshelf/internal/render"
)

type bookPageData struct {
	Books map[string][]dto.Book
}

func RenderBookPages(renderer *render.TemplateRenderer, bookshelf *dto.Bookshelf) error {
	for _, book := range bookshelf.Books {
		err := renderer.RenderToFile("book", book, book.Id)

		if err != nil {
			return err
		}
	}

	return nil
}
