package app

import (
	"log"
	"net/http"

	"github.com/fredbi/demo-api/app/assets"
	"github.com/fredbi/demo-api/pkg/repo/images"
	"github.com/go-chi/chi"
)

// Server holds all dependencies to run an API server.
//
// This is a barebone implementation. In a production environment,
// this "runtime" environment would be enriched with other much needed
// dependencies typically used to run this as a microservice, such as TLS
// material, logger, trace exporter...
//
// Further, the persistent store is currently backed by memory only, so
// as to provide a zero-config demo. Obviously, it is trivial to move
// this to some configurable permanent storage.
type Server struct {
	app *App
}

// New instance of an API server.
//
// In a typical production environment, this is where we gather paramenters and dependencies,
// such as tracing, logging, metrics, etc.
func New() *Server {
	return &Server{
		app: NewApp(images.New()), // inject persistent repository (badgerDB implementation)
	}
}

// Start the API server
func (s *Server) Start() error {
	r := chi.NewRouter()

	// webapp route
	r.Route("/", func(r chi.Router) {
		// in a typical production deployment, this would serve
		// locally built assets, mounted in the container.
		//
		// Alternatively, webapp serving could be handed over to
		// edge component, typically a CRD.
		r.Method(http.MethodGet, "/*", http.FileServer(assets.Static))
	})

	// register API routes
	//
	// NOTE: in a typical production environment, we inject here the appropriate middleware
	// to handle authentication, CORS, rate limiting, crash recovery, tracing etc.
	r.Route("/images", func(r chi.Router) {
		r.Post("/", s.app.CreateImage)
		r.Get("/", s.app.ListImages)
		r.Get("/{name}", s.app.GetImage)
		r.Patch("/{name}", s.app.UpdateImage)
		r.Delete("/{name}", s.app.DeleteImage)
	})

	log.Printf("listening on 0.0.0.0:3000")
	return http.ListenAndServe(":3000", r)
}

// NOTE: on a production environment, we should add there a graceful shutdown of the container,
// meaning essentially gracefully shutting down user connections and database connections.
