package router

import (
	"net/http"

	"covuavn/internal/article"
	"covuavn/internal/auth"
	"covuavn/internal/chat"
	"covuavn/internal/middleware"
	"covuavn/internal/notification"
	"covuavn/internal/user"

	"github.com/go-chi/chi/v5"
)

func New(
	userHandler *user.Handler,
	authHandler *auth.Handler,
	chatHandler *chat.Handler,
	notificationHandler *notification.Handler,
	articleHandler *article.Handler,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.JSON)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Get("/", userHandler.List)
	})

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Delete("/unregister", authHandler.Unregister)
	})

	r.Route("/v1/chat/messages", func(r chi.Router) {
		r.Post("/", chatHandler.Create)
		r.Get("/", chatHandler.List)
	})

	r.Route("/v1/notifications", func(r chi.Router) {
		r.Post("/", notificationHandler.Create)
		r.Get("/", notificationHandler.List)
	})

	r.Route("/v1/articles", func(r chi.Router) {
		r.Post("/", articleHandler.Create)
		r.Get("/", articleHandler.List)
	})

	return r
}
