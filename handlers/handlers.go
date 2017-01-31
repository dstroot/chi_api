package handler

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/dstroot/chi_api/models"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

// Article struct
type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var articles = []*Article{
	{ID: "1", Title: "Hi"},
	{ID: "2", Title: "sup"},
}

// https://github.com/golang/lint/pull/245
// Any package using context.WithValue and defining key types should either:
//
// - export their context key type
// - provide a constructor for it
// - provide some other helper func which does the context.WithValue behind the scenes.
//
// You really don't want to use strings as the keys. That means there's no
// isolated namespaces between packages, so different packages can collide
// and use the same keys, since everybody has access to the string type.
// But by requiring people to use their own types, there can't be conflicts.

// Key is the context key for this library.
type Key struct {
	article string
}

var key Key

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		article, err := dbGetArticle(articleID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, http.StatusText(http.StatusNotFound))
			return
		}

		// set context key
		key.article = "article"
		ctx := context.WithValue(r.Context(), key, article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SearchArticles searches the Articles data for a matching article.
// It's just a stub, but you get the idea.
func SearchArticles(w http.ResponseWriter, r *http.Request) {
	// Filter by query param, and search...
	render.JSON(w, r, articles)
}

// ListArticles returns an array of Articles.
func ListArticles(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, articles)
}

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	var data struct {
		*Article
		OmitID interface{} `json:"id,omitempty"` // prevents 'id' from being set
	}
	// ^ the above is a nifty trick for how to omit fields during json unmarshalling
	// through struct composition

	if err := render.Bind(r.Body, &data); err != nil {
		render.JSON(w, r, err.Error())
		return
	}

	article := data.Article
	dbNewArticle(article)

	render.JSON(w, r, article)
}

// GetArticle returns the specific Article. You'll notice it just
// fetches the Article right off the context, as its understood that
// if we made it this far, the Article must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetArticle(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value(key).(*Article)

	// chi provides a basic companion subpackage "github.com/pressly/chi/render", however
	// you can use any responder compatible with net/http.
	render.JSON(w, r, article)
}

// UpdateArticle updates an existing Article in our persistent store.
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value(key).(*Article)

	data := struct {
		*Article
		OmitID interface{} `json:"id,omitempty"` // prevents 'id' from being overridden
	}{Article: article}

	if err := render.Bind(r.Body, &data); err != nil {
		render.JSON(w, r, err)
		return
	}
	article = data.Article

	render.JSON(w, r, article)
}

// DeleteArticle removes an existing Article from our persistent store.
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value(key).(*Article)

	article, err = dbRemoveArticle(article.ID)
	if err != nil {
		render.JSON(w, r, err)
		return
	}

	// Respond with the deleted object, up to you.
	render.JSON(w, r, article)
}

// AdminRouter is a completely separate router for administrator routes
func AdminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: index"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: list accounts.."))
	})
	r.Get("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
	})
	return r
}

// AdminOnly middleware restricts access to just administrators.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

//--

// Below are a bunch of helper functions that mock some kind of storage

func dbNewArticle(article *Article) (string, error) {
	article.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	articles = append(articles, article)
	return article.ID, nil
}

func dbGetArticle(id string) (*Article, error) {
	for _, a := range articles {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("article not found")
}

func dbRemoveArticle(id string) (*Article, error) {
	for i, a := range articles {
		if a.ID == id {
			articles = append((articles)[:i], (articles)[i+1:]...)
			return a, nil
		}
	}
	return nil, errors.New("article not found")
}

// TaxPro Route
func TaxPro(w http.ResponseWriter, r *http.Request) {

	// Get params
	efin := chi.URLParam(r, "efin") // c.Param("efin")
	year := chi.URLParam(r, "year") // c.Param("year")

	// Get tax professionals
	results, err := models.GetTaxpro(year, efin)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, http.StatusText(http.StatusNotFound))
		render.JSON(w, r, err.Error())
		return
	}

	// Render results
	render.JSON(w, r, results)
}
