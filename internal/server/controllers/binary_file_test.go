package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/goph-keeper/internal/logger"
	mock2 "github.com/zelas91/goph-keeper/internal/server/controllers/mocks"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"golang.org/x/net/context"
)

type upload struct {
	name                      string
	mockBehaviorUploadService func(s *mock2.MockbinaryFileService, binary models.BinaryFile)
	want                      int
	url                       string
	model                     models.BinaryFile
	buffer                    []byte
}

func TestUpload(t *testing.T) {

	tests := []upload{
		{
			name: "#1 ok",
			mockBehaviorUploadService: func(s *mock2.MockbinaryFileService, binary models.BinaryFile) {
				s.EXPECT().Upload(gomock.Any(), binary, gomock.Any()).DoAndReturn(func(ctx context.Context, bf models.BinaryFile, reader <-chan []byte) error {
					for range reader {

					}
					return nil
				})
			},
			want: websocket.CloseNormalClosure,
			url:  "/upload",
			model: models.BinaryFile{
				FileName: "config.yaml",
				Size:     12,
			},
			buffer: []byte(`English texts for beginners to practice reading and comprehension online and for free.`),
		},
		{
			name: "#2 nok service save err",
			mockBehaviorUploadService: func(s *mock2.MockbinaryFileService, binary models.BinaryFile) {
				s.EXPECT().Upload(gomock.Any(), binary, gomock.Any()).DoAndReturn(func(ctx context.Context, bf models.BinaryFile, reader <-chan []byte) error {
					for range reader {

					}
					return errors.New("save error")
				})
			},
			want: websocket.CloseInternalServerErr,
			url:  "/upload",
			model: models.BinaryFile{
				FileName: "config.yaml",
				Size:     12,
			},
			buffer: []byte(`English texts for beginners to practice reading and comprehension online and for free.`),
		},
		{
			name: "#3 valid err size",
			want: websocket.CloseUnsupportedData,
			url:  "/upload",
			model: models.BinaryFile{
				FileName: "config.yaml",
			},
			buffer: []byte(`English texts for beginners to practice reading and comprehension online and for free.`),
		},
		{
			name: "#4 valid err file name",
			want: websocket.CloseUnsupportedData,
			url:  "/upload",
			model: models.BinaryFile{
				Size: 86,
			},
			buffer: []byte(`English texts for beginners to practice reading and comprehension online and for free.`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockbinaryFileService(ctrl)
			if test.mockBehaviorUploadService != nil {
				test.mockBehaviorUploadService(service, test.model)
			}

			handler := New(logger.New(""), WithBinaryFileUseService(service))
			srv := httptest.NewServer(handler.binary.createRoutes())
			defer srv.Close()
			u := "ws" + strings.TrimPrefix(srv.URL, "http")
			err := clientWebsocketUpload(test, u)

			var webErr *websocket.CloseError
			if errors.As(err, &webErr) {
				assert.Equal(t, test.want, webErr.Code)
				return
			}

			t.Fatalf("no status websocket connect err:%v", err)

		})
	}
}
func clientWebsocketUpload(test upload, url string) error {
	c, _, err := websocket.DefaultDialer.Dial(url+test.url, nil)
	if err != nil {
		return errors.New("error open websocket client")
	}
	defer c.Close()

	err = c.WriteJSON(test.model)
	if err != nil {
		return fmt.Errorf("write msg file info err: %w", err)
	}
	_, msg, err := c.ReadMessage()
	if err != nil {
		return fmt.Errorf("read answer err: %w", err)
	}
	var answer models.AnswerBinaryFile
	err = json.Unmarshal(msg, &answer)
	if err != nil {
		return fmt.Errorf("unmarshal answer err:%w ", err)
	}

	err = c.WriteMessage(websocket.BinaryMessage, test.buffer)
	if err != nil {
		return fmt.Errorf("write binary err:%w ", err)
	}
	if err := c.WriteMessage(websocket.TextMessage,
		[]byte("Binary data transfer completed")); err != nil {
		return fmt.Errorf("websocket send msg end file err: %w", err)
	}
	mt, _, err := c.ReadMessage()
	if err != nil {
		return fmt.Errorf("read confirmation messages succsess err:%w ", err)
	}
	return fmt.Errorf("websocket no close connection return msgType=%d", mt)
}

type download struct {
	name                        string
	mockBehaviorDownloadService func(s *mock2.MockbinaryFileService, binary models.BinaryFile)
	want                        int
	url                         string
	model                       models.BinaryFile
}

func TestDownload(t *testing.T) {
	tests := []download{
		{
			name: "#1 ok download",
			want: websocket.CloseNormalClosure,
			url:  "/download",
			model: models.BinaryFile{
				FileName: "logger.log",
				Size:     86,
			},
			mockBehaviorDownloadService: func(s *mock2.MockbinaryFileService, binary models.BinaryFile) {
				s.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, bf models.BinaryFile, write chan<- []byte) error {
					close(write)
					return nil
				})
			},
		},
		{
			name: "#2 nok service error",
			want: websocket.CloseInternalServerErr,
			url:  "/download",
			model: models.BinaryFile{
				FileName: "logger.log",
				Size:     86,
			},
			mockBehaviorDownloadService: func(s *mock2.MockbinaryFileService, binary models.BinaryFile) {
				s.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, bf models.BinaryFile, write chan<- []byte) error {
					close(write)
					return errors.New("read file err")
				})
			},
		},
		{
			name: "#3 valid err size",
			want: websocket.CloseUnsupportedData,
			url:  "/upload",
			model: models.BinaryFile{
				FileName: "config.yaml",
			},
		},
		{
			name: "#4 valid err file name",
			want: websocket.CloseUnsupportedData,
			url:  "/upload",
			model: models.BinaryFile{
				Size: 86,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockbinaryFileService(ctrl)
			if test.mockBehaviorDownloadService != nil {
				test.mockBehaviorDownloadService(service, test.model)
			}

			handler := New(logger.New(""), WithBinaryFileUseService(service))
			srv := httptest.NewServer(handler.binary.createRoutes())
			defer srv.Close()
			u := "ws" + strings.TrimPrefix(srv.URL, "http")
			err := clientWebsocketDownload(test, u)

			var webErr *websocket.CloseError
			if errors.As(err, &webErr) {
				assert.Equal(t, test.want, webErr.Code)
				return
			}

			t.Fatalf("no status websocket connect err%v", err)

		})
	}
}
func clientWebsocketDownload(test download, url string) error {
	c, _, err := websocket.DefaultDialer.Dial(url+test.url, nil)
	if err != nil {
		return errors.New("error open websocket client")
	}
	defer c.Close()

	err = c.WriteJSON(test.model)
	if err != nil {
		return fmt.Errorf("write msg file info err: %w", err)
	}
	_, msg, err := c.ReadMessage()
	if err != nil {
		return fmt.Errorf("read answer err: %w", err)
	}
	var answer models.AnswerBinaryFile
	err = json.Unmarshal(msg, &answer)
	if err != nil {
		return fmt.Errorf("unmarshal answer err:%w ", err)
	}

	mt, _, err := c.ReadMessage()
	if err != nil {
		return err
	}
	return fmt.Errorf("websocket no close connection return msgType=%d", mt)
}

func TestBinaryFileDelete(t *testing.T) {
	tests := []struct {
		name                      string
		want                      int
		method                    string
		url                       string
		fileIndex                 int
		mockBehaviorDeleteService func(s *mock2.MockbinaryFileService, fileID int)
	}{
		{
			name:      "#1 ok delete",
			want:      http.StatusOK,
			method:    http.MethodDelete,
			url:       "/",
			fileIndex: 5,
			mockBehaviorDeleteService: func(s *mock2.MockbinaryFileService, fileID int) {
				s.EXPECT().Delete(gomock.Any(), fileID).Return(nil)

			},
		},

		{
			name:      "#2 nok file delete err",
			want:      http.StatusInternalServerError,
			method:    http.MethodDelete,
			url:       "/",
			fileIndex: 5,
			mockBehaviorDeleteService: func(s *mock2.MockbinaryFileService, fileID int) {
				s.EXPECT().Delete(gomock.Any(), fileID).Return(errors.New("file delete error"))

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockbinaryFileService(ctrl)
			if test.mockBehaviorDeleteService != nil {
				test.mockBehaviorDeleteService(service, test.fileIndex)
			}

			url := fmt.Sprintf("%s%d", test.url, test.fileIndex)

			handler := New(logger.New(""), WithBinaryFileUseService(service))

			r := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()

			handler.binary.createRoutes().ServeHTTP(w, r)
			result := w.Result()
			defer w.Result().Body.Close()
			assert.Equal(t, test.want, result.StatusCode)
		})
	}
}

func TestFiles(t *testing.T) {

	tests := []struct {
		name                     string
		url                      string
		want                     int
		method                   string
		mockBehaviorFilesService func(s *mock2.MockbinaryFileService)
	}{
		{
			name:   "#1 ok get files",
			url:    "/",
			want:   http.StatusOK,
			method: http.MethodGet,
			mockBehaviorFilesService: func(s *mock2.MockbinaryFileService) {
				s.EXPECT().Files(gomock.Any()).Return([]models.BinaryFile{
					{
						FileName: "logg.log",
						ID:       1,
					},
					{
						FileName: "go.pdf",
						ID:       2,
					},
					{
						FileName: "TCP.log",
						ID:       3,
					},
				}, nil)
			},
		},
		{
			name:   "#2 nok get files error",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodGet,
			mockBehaviorFilesService: func(s *mock2.MockbinaryFileService) {
				s.EXPECT().Files(gomock.Any()).Return(nil, errors.New("repo error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockbinaryFileService(ctrl)
			if test.mockBehaviorFilesService != nil {
				test.mockBehaviorFilesService(service)
			}

			handler := New(logger.New(""), WithBinaryFileUseService(service))

			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := handler.binary.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestFile(t *testing.T) {

	tests := []struct {
		name                    string
		url                     string
		want                    int
		method                  string
		fileIndex               int
		mockBehaviorFileService func(s *mock2.MockbinaryFileService, fileID int)
	}{
		{
			name:      "#1 ok get file",
			fileIndex: 2,
			url:       "/",
			want:      http.StatusOK,
			method:    http.MethodGet,
			mockBehaviorFileService: func(s *mock2.MockbinaryFileService, fileID int) {
				s.EXPECT().File(gomock.Any(), fileID).Return(models.BinaryFile{
					ID:       fileID,
					FileName: "myfile",
					Size:     86,
				}, nil)
			},
		},
		{
			name:   "#2 nok file get id from request err",
			url:    "/a",
			want:   http.StatusBadRequest,
			method: http.MethodGet,
		},
		{
			name:      "#3 nok file: get file err",
			fileIndex: 2,
			url:       "/",
			want:      http.StatusNotFound,
			method:    http.MethodGet,
			mockBehaviorFileService: func(s *mock2.MockbinaryFileService, fileID int) {
				s.EXPECT().File(gomock.Any(), fileID).Return(models.BinaryFile{}, errors.New("file: get file err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockbinaryFileService(ctrl)
			if test.mockBehaviorFileService != nil {
				test.mockBehaviorFileService(service, test.fileIndex)
			}

			handler := New(logger.New(""), WithBinaryFileUseService(service))

			url := fmt.Sprintf("%s%d", test.url, test.fileIndex)
			request := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()
			h := handler.binary.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}
