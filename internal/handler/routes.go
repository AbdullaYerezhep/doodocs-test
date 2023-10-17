package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)



func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		
		r.Route("/archive", func(r chi.Router) {
			r.Post("/information", h.GetArchiveInfo)
			r.Post("/files", h.HandleCreateArchive)
		})

		r.Route("/mail", func(r chi.Router){
			r.Post("/file", h.HandleSendFileByEmail)
		})

	})
	return r
}
