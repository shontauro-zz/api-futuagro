package http

import (
	"net/http"
	"time"

	"futuagro.com/pkg/config"
	"futuagro.com/pkg/domain/services"
	"futuagro.com/pkg/http/rest"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	config          *config.Config
	supplierService *services.SupplierService
	countryService  *services.CountryService
	router          chi.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Run starts a http server
func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s,
	}
	httpServer.ListenAndServe()
}

// AllowOriginFunc Definie which origins our http servers accepts request from
func AllowOriginFunc(r *http.Request, origin string) bool {
	// accept all (*) origins
	return true
}

// NewServer returns a new HTTP server.
func NewServer(confPtr *config.Config, supplierServ *services.SupplierService, countryServ *services.CountryService) *Server {
	server := &Server{
		config:          confPtr,
		supplierService: supplierServ,
		countryService:  countryServ,
	}

	r := chi.NewRouter()
	// Setup CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"ETag", "Link", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset", "X-OAuth-Scopes", "X-Accepted-OAuth-Scopes"},
		AllowCredentials: true,
		MaxAge:           3600, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(5 * time.Minute))

	rSupplier := rest.SupplierHandler{Service: supplierServ}
	rCountry := rest.CountryHandler{Service: countryServ}

	r.Mount("/suppliers", rSupplier.NewRouter())
	r.Mount("/countries", rCountry.NewRouter())

	server.router = r
	return server
}
