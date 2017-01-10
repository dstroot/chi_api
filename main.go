//
// REST
// ====
// This example demonstrates a HTTP REST web service with some fixture data.
// Follow along the example and patterns.
//
// Also check routes.json for the generated docs from passing the -routes flag
//
// Boot the server:
// ----------------
// $ go run main.go
//
// Client requests:
// ----------------
// $ curl http://localhost:3333/
// root.
//
// $ curl http://localhost:3333/articles
// [{"id":"1","title":"Hi"},{"id":"2","title":"sup"}]
//
// $ curl http://localhost:3333/articles/1
// {"id":"1","title":"Hi"}
//
// $ curl -X DELETE http://localhost:3333/articles/1
// {"id":"1","title":"Hi"}
//
// $ curl http://localhost:3333/articles/1
// "Not Found"
//
// $ curl -X POST -d '{"id":"will-be-omitted","title":"awesomeness"}' http://localhost:3333/articles
// {"id":"97","title":"awesomeness"}
//
// $ curl http://localhost:3333/articles/97
// {"id":"97","title":"awesomeness"}
//
// $ curl http://localhost:3333/articles
// [{"id":"2","title":"sup"},{"id":"97","title":"awesomeness"}]
//
package main

import (
	"net/http"
	"time"

	"github.com/dstroot/chi_api/handlers"
	"github.com/goware/httpcoala"
	"github.com/pressly/chi"
	"github.com/pressly/chi/docgen"
	"github.com/pressly/chi/middleware"
	"github.com/russross/blackfriday"
)

func main() {

	r := chi.NewRouter()

	/**
	 * MIDDLEWARE
	 */

	// Injects a request ID into the context of each request.
	r.Use(middleware.RequestID)
	// Logs the start and end of each request with the elapsed processing time.
	r.Use(middleware.Logger)
	// Gracefully absorb panics and prints the stack trace.
	r.Use(middleware.Recoverer)
	// Stop processing after 2.5 seconds.
	r.Use(middleware.Timeout(2500 * time.Millisecond))
	// Only one request will be processed at a time.
	r.Use(middleware.Throttle(100))
	// Middleware handler that routes multiple requests for the
	// same URI (and routed methods) to be processed as a single request.
	r.Use(httpcoala.Route("HEAD", "GET")) // or, Route("*")
	// Health route for Heartbeat/load balancers
	r.Use(middleware.Heartbeat("/health"))

	/**
	 * ROUTES
	 */

	// RESTy routes for "articles" resource
	r.Route("/articles", func(r chi.Router) {
		r.With(handler.Paginate).Get("/", handler.ListArticles)
		r.Post("/", handler.CreateArticle)       // POST /articles
		r.Get("/search", handler.SearchArticles) // GET /articles/search

		r.Route("/:articleID", func(r chi.Router) {
			r.Use(handler.ArticleCtx)            // Load the *Article on the request context
			r.Get("/", handler.GetArticle)       // GET /articles/123
			r.Put("/", handler.UpdateArticle)    // PUT /articles/123
			r.Delete("/", handler.DeleteArticle) // DELETE /articles/123
		})
	})

	// Mount the admin sub-router, the same as a call to
	// Route("/admin", func(r chi.Router) { with routes here })
	r.Mount("/admin", handler.AdminRouter())

	// last so all routes are picked up in the docs
	md := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
		ProjectPath: "github.com/dstroot/chi_api",
		Intro:       "Welcome to the chi/_examples/rest generated docs.",
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		output := blackfriday.MarkdownCommon([]byte(md))
		w.Write([]byte(output))
	})

	/**
	 * SERVER
	 */

	http.ListenAndServe(":3333", r)
}
