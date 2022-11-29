package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.recoverPanic)

	mux.Get("/_/status", app.status)

	mux.Get("/products", app.productSearch)
	mux.Get("/products/facets", app.productFacets)

	return mux
}
