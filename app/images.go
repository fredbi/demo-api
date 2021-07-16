package app

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/fredbi/demo-api/pkg/repo"
	"github.com/go-chi/chi"
)

// ErrRequired indicates a required request parameter
var ErrRequired = errors.New("param required")

// App implements the API handlers.
//
// In this crude implementation, the API directly renders content as HTML.
//
// A more flexible implementation would typically convey content as JSON.
//
// Notice that file content is uploaded (for POST and PATCH) using a multipart MIME form,
// which allows for an extensible schema for metadata. This is also compatible with
// standard API schema definitions, such as SWAGGER.
type App struct {
	repo repo.ImagesRepo
}

// NewApp builds an API with injected dependencies
func NewApp(repo repo.ImagesRepo) *App {
	return &App{repo: repo}
}

func asHTML(w http.ResponseWriter) {
	w.Header().Set("content-type", "text/html")
}

func errorResponse(w http.ResponseWriter, err error) {
	// errorResponse converts known error types into the appropriate http status code.
	//
	// If we want error pages or messages, we could write an explicit response here rather
	// than handing this over to the frontend.
	var code int

	log.Printf("ERROR: %v", err)

	switch err {
	case repo.ErrNotFound:
		code = http.StatusNotFound
	case ErrRequired:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
}

func byName(req *http.Request) (string, error) {
	key := chi.URLParam(req, "name")
	if key == "" {
		key = req.FormValue("name")
	}

	if key == "" {
		return "", ErrRequired
	}

	return key, nil
}

// CreateImage uploads a new image to the repository
func (app *App) CreateImage(w http.ResponseWriter, req *http.Request) {
	file, header, err := req.FormFile("file")
	if err != nil {
		errorResponse(w, err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	key, _ := byName(req)
	if key == "" {
		key = header.Filename
	}
	if key == "" {
		errorResponse(w, ErrRequired)
		return
	}

	err = app.repo.Create(key, file)
	if err != nil {
		errorResponse(w, err)
		return
	}

	asHTML(w)
}

// ListImages lists all images available in the repository, with their thumbnail.
//
// This is rendered directly as server-side HTML.
//
// Better interoperability is of course achieved by using JSON and a smarter frontend
// to display the response. In this demo, we wanted the frontend to remain as small
// as possible, hence the choice for server rendering.
func (app *App) ListImages(w http.ResponseWriter, req *http.Request) {
	images, err := app.repo.List()
	if err != nil {
		errorResponse(w, err)
		return
	}

	var buf bytes.Buffer
	err = list.Execute(&buf, images)
	if err != nil {
		errorResponse(w, err)
		return
	}

	asHTML(w)
	_, _ = w.Write(buf.Bytes())
}

// GetImage retrieves one specific image from the repository
func (app *App) GetImage(w http.ResponseWriter, req *http.Request) {
	key, err := byName(req)
	if err != nil {
		errorResponse(w, err)
		return
	}

	file, err := app.repo.Get(key)
	if err != nil {
		errorResponse(w, err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(w, file)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

// UpdateImage updates an image in the repository from a multipart form.
func (app *App) UpdateImage(w http.ResponseWriter, req *http.Request) {
	file, header, err := req.FormFile("file")
	if err != nil {
		errorResponse(w, err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	key, _ := byName(req)
	if key == "" {
		key = header.Filename
	}
	if key == "" {
		errorResponse(w, ErrRequired)
		return
	}

	err = app.repo.Update(key, file)
	if err != nil {
		errorResponse(w, err)
		return
	}

	asHTML(w)
}

// DeleteImage removes an image from the repository
func (app *App) DeleteImage(w http.ResponseWriter, req *http.Request) {
	key, err := byName(req)
	if err != nil {
		errorResponse(w, err)

		return
	}

	err = app.repo.Delete(key)
	if err != nil {
		errorResponse(w, err)

		return
	}

	asHTML(w)
}
