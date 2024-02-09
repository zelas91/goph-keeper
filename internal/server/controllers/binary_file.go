package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/payload"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type binaryFile struct {
	log     logger.Logger
	valid   *validator.Validate
	service binaryFileService
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

//go:generate mockgen -package mocks -destination=./mocks/mock_binary_file_service.go -source=binary_file.go -package=mock
type binaryFileService interface {
	Upload(ctx context.Context, bf models.BinaryFile, reader <-chan []byte) error
	Download(ctx context.Context, bf models.BinaryFile, write chan<- []byte) error
	Delete(ctx context.Context, fileID int) error
	Files(ctx context.Context) ([]models.BinaryFile, error)
	File(ctx context.Context, fileID int) (models.BinaryFile, error)
}

func (b *binaryFile) upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			b.log.Errorf("Failed to upgrade connection:", err)
			return
		}

		defer func() {
			if err := conn.Close(); err != nil {
				b.log.Errorf("upload close websocket connect err: %v", err)
			}
		}()

		var bf models.BinaryFile

		if err = conn.ReadJSON(&bf); err != nil {
			b.log.Errorf("Failed to unmarshal JSON: %v", err)
			if err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData,
					fmt.Sprintf("upload: write message error err:%v", err))); err != nil {
				b.log.Errorf("upload: write message error err:%v", err)
				return
			}
			return
		}
		if err = b.valid.Struct(bf); err != nil {
			b.log.Errorf("upload: validation err:%v", err)
			if err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())); err != nil {
				b.log.Errorf("upload send message err:%v", err)
			}
			return
		}
		if err = conn.WriteJSON(models.AnswerBinaryFile{Confirm: true}); err != nil {
			b.log.Errorf("upload: answer send err: %v", err)
			return
		}
		reader := make(chan []byte)
		g, ctx := errgroup.WithContext(r.Context())

		g.Go(func() error {
			return b.service.Upload(ctx, bf, reader)
		})

		go func() {
			defer func() {
				close(reader)
			}()
			for {
				if ctx.Err() != nil {
					b.log.Errorf("download context err: %v", ctx.Err())
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
					mt, msg, err := conn.ReadMessage()
					if err != nil {
						b.log.Errorf("failed to read message: %v", err)
						return
					}
					if mt == websocket.BinaryMessage {
						reader <- msg
					} else {
						b.log.Debugf("message websocket text : %s", string(msg))
						return
					}
				}

			}
		}()

		if err = g.Wait(); err != nil {
			b.log.Errorf("upload service err: %v", err)

			body, err := json.Marshal(payload.ErrorMessage{
				Message:    fmt.Sprintf("service err :%v", err),
				StatusCode: websocket.CloseInternalServerErr,
			})
			if err == nil {
				if err = conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, string(body))); err != nil {
					b.log.Errorf("upload websocket send msg err: %v", err)
				}
			}
			return
		}

		if err = conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			b.log.Errorf("upload websocket send msg err: %v", err)
			return
		}

	}
}
func (b *binaryFile) download() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			b.log.Errorf("Failed to upgrade connection:", err)
			return
		}

		defer func() {
			if err := conn.Close(); err != nil {
				b.log.Errorf("download close websocket connect err: %v", err)
			}
		}()
		var bf models.BinaryFile

		if err = conn.ReadJSON(&bf); err != nil {
			b.log.Errorf("Failed to unmarshal JSON:", err)
			if err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData,
					fmt.Sprintf("download: write message error err:%v", err))); err != nil {
				b.log.Errorf("download: write message error err:%v", err)
				return
			}
			return
		}
		if err = b.valid.Struct(bf); err != nil {
			b.log.Errorf("download: validation err:%v", err)
			if err = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())); err != nil {
				b.log.Errorf("download websocket send msg err:%v", err)
			}
			return
		}
		if err = conn.WriteJSON(models.AnswerBinaryFile{Confirm: true}); err != nil {
			b.log.Errorf("download: answer send err: %v", err)
			return
		}
		writer := make(chan []byte)
		g, ctx := errgroup.WithContext(r.Context())
		g.Go(func() error {
			return b.service.Download(ctx, bf, writer)
		})
		for body := range writer {

			if err := conn.WriteMessage(websocket.BinaryMessage, body); err != nil {
				b.log.Errorf("download write binary websocket err: %v", err)
				return
			}
		}

		if err = g.Wait(); err != nil {
			b.log.Errorf("download service err: %v", err)

			body, err := json.Marshal(payload.ErrorMessage{
				Message:    fmt.Sprintf("service err :%v", err),
				StatusCode: websocket.CloseInternalServerErr,
			})
			if err == nil {
				if err = conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, string(body))); err != nil {
					b.log.Errorf("download websocket send msg err: %v", err)
				}
			}
			return
		}
		if err = conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			b.log.Errorf("download close websocket err:%v", err)
		}

	}
}
func (b *binaryFile) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			b.log.Errorf("delete get ID err: %v", err)
			payload.NewErrorResponse(w, "delete get id not found", http.StatusBadRequest)
			return
		}
		if err := b.service.Delete(r.Context(), id); err != nil {
			b.log.Errorf("delete err: %v", err)
			payload.NewErrorResponse(w, "delete err", http.StatusInternalServerError)
			return
		}

	}
}

func (b *binaryFile) Files() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := b.service.Files(r.Context())
		if err != nil {
			b.log.Errorf("files: get binary files err:%v ", err)
			payload.NewErrorResponse(w, "get binary files err", http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(files); err != nil {
			b.log.Errorf("files: files encode  err:%v ", err)
			payload.NewErrorResponse(w, "files encode  err", http.StatusInternalServerError)
			return
		}
	}
}

func (b *binaryFile) File() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			b.log.Errorf("file: get id from request err: %v", err)
			payload.NewErrorResponse(w, "file: get id from request err", http.StatusBadRequest)
			return
		}
		file, err := b.service.File(r.Context(), id)
		if err != nil {
			b.log.Errorf("file: get file err: %v", err)
			payload.NewErrorResponse(w, "file: get file err", http.StatusNotFound)
			return
		}

		if err = json.NewEncoder(w).Encode(file); err != nil {
			b.log.Errorf("file: encode err %v", err)
			payload.NewErrorResponse(w, "file: encode err", http.StatusInternalServerError)
			return
		}
	}
}

func (b *binaryFile) createRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Delete("/{id}", b.delete())
		r.Get("/upload", b.upload())
		r.Get("/download", b.download())
		r.Get("/", b.Files())
		r.Get("/{id}", b.File())
	})
	return router

}
