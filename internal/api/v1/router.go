package v1

import (
	"encoding/json"
	service "go-task-tracker/internal/services"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(r chi.Router, services *service.Services) {
	logWriter := setLogsFile()

	r.Use(chimw.RequestID) // TODO - ?
	r.Use(chimw.RealIP)    // TODO - ?
	r.Use(chimw.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ww := chimw.NewWrapResponseWriter(w, req.ProtoMajor)

			defer func() {
				entry := map[string]interface{}{
					"time":   time.Now().Format(time.RFC3339Nano),
					"method": req.Method,
					"uri":    req.RequestURI,
					"status": ww.Status(),
				}
				_ = json.NewEncoder(logWriter).Encode(entry)
			}()
			next.ServeHTTP(ww, req)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	// r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", func(cr chi.Router) {
		newAuthRoutes(cr, services.Auth)
	})

	authMiddleware := &AuthMiddleware{services.Auth}
	r.Route("/api/v1", func(api chi.Router) {
		api.Use(authMiddleware.UserIdentity)
		// Остальные группы хендлеров
	})
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		slog.Error("cannot open requests.log", "err", err)
		os.Exit(1)
	}
	return file
}
