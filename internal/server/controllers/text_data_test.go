package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/goph-keeper/internal/logger"
	mock2 "github.com/zelas91/goph-keeper/internal/server/controllers/mocks"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTexts(t *testing.T) {

	tests := []struct {
		name                     string
		url                      string
		want                     int
		method                   string
		mockBehaviorTextsService func(s *mock2.MocktextDataService)
	}{
		{
			name:   "#1 ok get texts",
			url:    "/",
			want:   http.StatusOK,
			method: http.MethodGet,
			mockBehaviorTextsService: func(s *mock2.MocktextDataService) {
				s.EXPECT().Texts(gomock.Any()).Return([]models.TextData{
					{
						Text: `Prepared by experienced English teachers, 
							the texts, articles and conversations are brief and appropriate 
							to your level of proficiency`,
					},
					{
						Text: `English texts for beginners to practice reading and comprehension online and for free.`,
					},
				}, nil)
			},
		},
		{
			name:   "#2 nok get texts error",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodGet,
			mockBehaviorTextsService: func(s *mock2.MocktextDataService) {
				s.EXPECT().Texts(gomock.Any()).Return(nil, errors.New("repo error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMocktextDataService(ctrl)
			if test.mockBehaviorTextsService != nil {
				test.mockBehaviorTextsService(service)
			}

			handler := New(logger.New(""), WithTextUseService(service))

			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := handler.textData.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestText(t *testing.T) {

	tests := []struct {
		name                    string
		url                     string
		want                    int
		method                  string
		textIndex               int
		mockBehaviorTextService func(s *mock2.MocktextDataService, textID int)
	}{
		{
			name:      "#1 ok get text",
			textIndex: 2,
			url:       "/",
			want:      http.StatusOK,
			method:    http.MethodGet,
			mockBehaviorTextService: func(s *mock2.MocktextDataService, textID int) {
				s.EXPECT().Text(gomock.Any(), textID).Return(models.TextData{
					ID:      textID,
					Text:    `English texts for beginners to practice reading and comprehension online and for free.`,
					Version: 1,
				}, nil)
			},
		},
		{
			name:   "#2 nok text get id from request err",
			url:    "/a",
			want:   http.StatusBadRequest,
			method: http.MethodGet,
		},
		{
			name:      "#3 nok text: get text err",
			textIndex: 2,
			url:       "/",
			want:      http.StatusNotFound,
			method:    http.MethodGet,
			mockBehaviorTextService: func(s *mock2.MocktextDataService, textID int) {
				s.EXPECT().Text(gomock.Any(), textID).Return(models.TextData{},
					errors.New("text: get text err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMocktextDataService(ctrl)
			if test.mockBehaviorTextService != nil {
				test.mockBehaviorTextService(service, test.textIndex)
			}

			handler := New(logger.New(""), WithTextUseService(service))

			url := fmt.Sprintf("%s%d", test.url, test.textIndex)
			request := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()
			h := handler.textData.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestTextCreate(t *testing.T) {

	tests := []struct {
		name                      string
		url                       string
		want                      int
		method                    string
		mockBehaviorCreateService func(s *mock2.MocktextDataService, text models.TextData)
		body                      models.TextData
	}{
		{
			name:   "#1 ok create text",
			url:    "/",
			want:   http.StatusCreated,
			method: http.MethodPost,
			mockBehaviorCreateService: func(s *mock2.MocktextDataService, text models.TextData) {
				s.EXPECT().Create(gomock.Any(), text).Return(nil)
			},
			body: models.TextData{
				Text: `English texts for beginners to practice reading and comprehension online and for free.`,
			},
		},
		{
			name:   "#2 nok validation error text",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPost,
			body:   models.TextData{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMocktextDataService(ctrl)
			if test.mockBehaviorCreateService != nil {
				test.mockBehaviorCreateService(service, test.body)
			}

			handler := New(logger.New(""), WithTextUseService(service))

			body, err := json.Marshal(test.body)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, bytes.NewReader(body))
			w := httptest.NewRecorder()
			h := handler.textData.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestTextsUpdate(t *testing.T) {

	tests := []struct {
		name                      string
		url                       string
		want                      int
		method                    string
		mockBehaviorUpdateService func(s *mock2.MocktextDataService, text models.TextData)
		body                      models.TextData
	}{
		{
			name:   "#1 ok update",
			url:    "/",
			want:   http.StatusOK,
			method: http.MethodPut,
			body: models.TextData{
				ID:      1,
				Text:    `English texts for beginners to practice reading and comprehension online and for free.`,
				Version: 3,
			},
			mockBehaviorUpdateService: func(s *mock2.MocktextDataService, text models.TextData) {
				s.EXPECT().Update(gomock.Any(), text).Return(nil)
			},
		},

		{
			name:   "#2 nok get id text err",
			url:    "/qwe",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
		},

		{
			name:   "#3 nok validation error text",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.TextData{
				ID:      2,
				Version: 5,
			},
		},
		{
			name:   "#5 nok validation error version",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.TextData{
				ID:   2,
				Text: `English texts for beginners to practice reading and comprehension online and for free.`,
			},
		},

		{
			name:   "#6 nok text save err",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodPut,
			body: models.TextData{
				ID:      2,
				Text:    `English texts for beginners to practice reading and comprehension online and for free.`,
				Version: 9,
			},
			mockBehaviorUpdateService: func(s *mock2.MocktextDataService, text models.TextData) {
				s.EXPECT().Update(gomock.Any(), text).Return(errors.New("save repo err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMocktextDataService(ctrl)
			if test.mockBehaviorUpdateService != nil {
				test.mockBehaviorUpdateService(service, test.body)
			}

			handler := New(logger.New(""), WithTextUseService(service))

			body, err := json.Marshal(test.body)
			assert.NoError(t, err, "Body write error")

			path := fmt.Sprintf("%s%d", test.url, test.body.ID)
			request := httptest.NewRequest(test.method, path, bytes.NewReader(body))
			w := httptest.NewRecorder()
			h := handler.textData.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want, res.StatusCode)

		})
	}
}

func TestTextDelete(t *testing.T) {
	tests := []struct {
		name                      string
		want                      int
		method                    string
		url                       string
		textIndex                 int
		mockBehaviorDeleteService func(s *mock2.MocktextDataService, textID int)
	}{
		{
			name:      "#1 ok delete",
			want:      http.StatusOK,
			method:    http.MethodDelete,
			url:       "/",
			textIndex: 9,
			mockBehaviorDeleteService: func(s *mock2.MocktextDataService, textID int) {
				s.EXPECT().Delete(gomock.Any(), textID).Return(nil)

			},
		},

		{
			name:      "#2 nok text delete err",
			want:      http.StatusInternalServerError,
			method:    http.MethodDelete,
			url:       "/",
			textIndex: 9,
			mockBehaviorDeleteService: func(s *mock2.MocktextDataService, textID int) {
				s.EXPECT().Delete(gomock.Any(), textID).Return(errors.New("text delete error"))

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMocktextDataService(ctrl)
			if test.mockBehaviorDeleteService != nil {
				test.mockBehaviorDeleteService(service, test.textIndex)
			}

			url := fmt.Sprintf("%s%d", test.url, test.textIndex)

			handler := New(logger.New(""), WithTextUseService(service))

			r := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()

			handler.textData.createRoutes().ServeHTTP(w, r)
			result := w.Result()
			defer w.Result().Body.Close()
			assert.Equal(t, test.want, result.StatusCode)
		})
	}
}
