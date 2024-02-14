package controllers

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/server/payload"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/middleware"
	"github.com/zelas91/goph-keeper/internal/utils/validation"
)

var clientAppDir = "build/client"

type Controllers struct {
	auth       *auth
	card       *сreditCard
	credential *credential
	textData   *textData
	binary     *binaryFile
	log        logger.Logger
	valid      *validator.Validate
}

func New(log logger.Logger, options ...func(c *Controllers)) *Controllers {
	ctl := &Controllers{
		log:   log,
		valid: validation.NewValidator(log),
	}
	for _, opt := range options {
		opt(ctl)
	}
	return ctl
}

func WithAuthUseService(us userService) func(c *Controllers) {
	return func(c *Controllers) {
		c.auth = &auth{service: us, valid: c.valid, log: c.log}
	}
}

func WithCardUseService(cs cardService) func(c *Controllers) {
	return func(c *Controllers) {
		c.card = &сreditCard{service: cs, valid: c.valid, log: c.log}
	}
}

func WithUserCredentialUseService(cs credentialService) func(c *Controllers) {
	return func(c *Controllers) {
		c.credential = &credential{service: cs, valid: c.valid, log: c.log}
	}
}

func WithTextUseService(td textDataService) func(c *Controllers) {
	return func(c *Controllers) {
		c.textData = &textData{service: td, valid: c.valid, log: c.log}
	}
}

func WithBinaryFileUseService(bs binaryFileService) func(c *Controllers) {
	return func(c *Controllers) {
		c.binary = &binaryFile{service: bs, valid: c.valid, log: c.log}
	}
}
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(clientAppDir)
	if err != nil {
		payload.NewErrorResponse(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileList = append(fileList, file.Name())
	}

	html := "<h1>Files client:</h1><ul>"
	for _, file := range fileList {
		html += "<li><a href='/files/client/" + file + "'>" + file + "</a></li>"
	}
	html += "</ul>"
	if _, err := w.Write([]byte(html)); err != nil {
		payload.NewErrorResponse(w, "", http.StatusInternalServerError)
		return
	}
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	path := fmt.Sprintf("%s/%s", clientAppDir, fileName)
	_, err := os.Stat(path)
	if err != nil {
		payload.NewErrorResponse(w, "File not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}

func (c *Controllers) CreateRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/files/client", listFilesHandler)
	router.Get("/files/client/{fileName}", downloadFileHandler)
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.ContentTypeJSON(c.log), middleware2.Recoverer)
		r.Mount("/", c.auth.createRoutes())
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthorizationHandler(c.log, c.auth.service))
			r.Group(func(r chi.Router) {
				r.Mount("/card", c.card.createRoutes())
				r.Mount("/credential", c.credential.createRoutes())
				r.Mount("/text", c.textData.createRoutes())
				r.Mount("/file", c.binary.createRoutes())
			})
		})
	})
	return router
}
