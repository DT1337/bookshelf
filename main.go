package main

import (
	"log"

	"bookshelf/internal/dto"
	"bookshelf/internal/pages"
	"bookshelf/internal/render"
)

func main() {
	bookshelf, err := dto.LoadBookshelfFromFile("data/data.json")
	if err != nil {
		log.Fatal(err)
	}

	config := render.TemplateRendererConfig{
		TemplateType:           "html",
		TemplatesPath:          "templates",
		ComponentTemplatesPath: "components",
		PageTemplatesPath:      "pages",
		OutputPath:             "dist",
		BaseTemplateName:       "base",
	}

	renderer, err := render.New(config)
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}

	err = renderer.CopyStaticFiles("static", config.OutputPath)
	if err != nil {
		log.Fatalf("Failed to copy static files: %v", err)
	}

	err = pages.RenderIndexPage(renderer, bookshelf)
	if err != nil {
		log.Fatalf("Failed to render index page: %v", err)
	}

	err = pages.RenderBookshelfPage(renderer, bookshelf)
	if err != nil {
		log.Fatalf("Failed to render bookshelf page: %v", err)
	}

	err = pages.RenderCollectionsPage(renderer, bookshelf)
	if err != nil {
		log.Fatalf("Failed to render collections page: %v", err)
	}

	err = pages.RenderQuotesPage(renderer, bookshelf)
	if err != nil {
		log.Fatalf("Failed to render quotes page: %v", err)
	}

	err = pages.RenderWishlistPage(renderer, bookshelf)
	if err != nil {
		log.Fatalf("Failed to render wishlist page: %v", err)
	}

	err = pages.RenderBookPages(renderer, bookshelf)
	if err != nil {
		log.Fatalf("Failed to render book pages: %v", err)
	}

	log.Println("Pages and static files rendered successfully!")
}
