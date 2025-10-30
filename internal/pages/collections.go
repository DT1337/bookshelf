package pages

import (
	"bookshelf/internal/dto"
	"bookshelf/internal/render"
)

type collectionsPageData struct {
	Collections []dto.ResolvedCollection
}

func RenderCollectionsPage(renderer *render.TemplateRenderer, bookshelf *dto.Bookshelf) error {
	data := collectionsPageData{
		Collections: bookshelf.BookCollections(),
	}

	return renderer.RenderToFile("collections", data)
}
