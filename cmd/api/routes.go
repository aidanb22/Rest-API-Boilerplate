package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.authenticate)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	mux.Get("/players", app.AllPlayers)
	mux.Get("/players/{id}", app.GetPlayer)

	mux.Get("/users", app.AllUsers)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/players", app.PlayerDatabase)
		mux.Get("/players/{id}", app.PlayerForEdit)
		mux.Put("/players/0", app.InsertPlayer)
		mux.Patch("/players/{id}", app.UpdatePlayer)
		mux.Delete("/players/{id}", app.DeletePlayer)
	})

	return mux
}
