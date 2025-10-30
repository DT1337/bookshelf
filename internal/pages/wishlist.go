package pages

import (
	"bookshelf/internal/dto"
	"bookshelf/internal/render"
)

type wishlistPageData struct {
	Books []dto.Book
}

func RenderWishlistPage(renderer *render.TemplateRenderer, bookshelf *dto.Bookshelf) error {
	data := wishlistPageData{
		Books: bookshelf.WishlistedBooks(),
	}

	return renderer.RenderToFile("wishlist", data)
}
