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

func TestCards(t *testing.T) {

	tests := []struct {
		name                     string
		url                      string
		want                     int
		method                   string
		mockBehaviorCardsService func(s *mock2.MockcardService)
	}{
		{
			name:   "#1 ok get cards",
			url:    "/",
			want:   http.StatusOK,
			method: http.MethodGet,
			mockBehaviorCardsService: func(s *mock2.MockcardService) {
				s.EXPECT().Cards(gomock.Any()).Return([]models.Card{
					{
						Number:    "5500126132422715",
						Cvv:       "123",
						ExpiredAt: "12/26",
					},
					{
						Number:    "4361913530390185",
						Cvv:       "456",
						ExpiredAt: "09/27",
					},
					{
						Number:    "6011941309164282",
						Cvv:       "986",
						ExpiredAt: "05/26",
					},
					{
						Number:    "378193510265868",
						Cvv:       "345",
						ExpiredAt: "01/28",
					},
				}, nil)
			},
		},
		{
			name:   "#2 nok get cards error",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodGet,
			mockBehaviorCardsService: func(s *mock2.MockcardService) {
				s.EXPECT().Cards(gomock.Any()).Return(nil, errors.New("repo error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcardService(ctrl)
			if test.mockBehaviorCardsService != nil {
				test.mockBehaviorCardsService(service)
			}

			handler := New(logger.New(""), WithCardUseService(service))

			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := handler.card.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestCard(t *testing.T) {

	tests := []struct {
		name                    string
		url                     string
		want                    int
		method                  string
		cardIndex               int
		mockBehaviorCardService func(s *mock2.MockcardService, cardID int)
	}{
		{
			name:      "#1 ok get card",
			cardIndex: 2,
			url:       "/",
			want:      http.StatusOK,
			method:    http.MethodGet,
			mockBehaviorCardService: func(s *mock2.MockcardService, cardID int) {
				s.EXPECT().Card(gomock.Any(), cardID).Return(models.Card{
					ID:        cardID,
					Number:    "5500126132422715",
					Cvv:       "123",
					ExpiredAt: "12/26",
					Version:   1,
				}, nil)
			},
		},
		{
			name:   "#2 nok card get id from request err",
			url:    "/a",
			want:   http.StatusBadRequest,
			method: http.MethodGet,
		},
		{
			name:      "#3 nok card: get card err",
			cardIndex: 2,
			url:       "/",
			want:      http.StatusNotFound,
			method:    http.MethodGet,
			mockBehaviorCardService: func(s *mock2.MockcardService, cardID int) {
				s.EXPECT().Card(gomock.Any(), cardID).Return(models.Card{}, errors.New("card: get card err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcardService(ctrl)
			if test.mockBehaviorCardService != nil {
				test.mockBehaviorCardService(service, test.cardIndex)
			}

			handler := New(logger.New(""), WithCardUseService(service))

			url := fmt.Sprintf("%s%d", test.url, test.cardIndex)
			request := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()
			h := handler.card.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestCardCreate(t *testing.T) {

	tests := []struct {
		name                      string
		url                       string
		want                      int
		method                    string
		mockBehaviorCreateService func(s *mock2.MockcardService, card models.Card)
		body                      models.Card
	}{
		{
			name:   "#1 ok create card",
			url:    "/",
			want:   http.StatusCreated,
			method: http.MethodPost,
			mockBehaviorCreateService: func(s *mock2.MockcardService, card models.Card) {
				s.EXPECT().Create(gomock.Any(), card).Return(nil)
			},
			body: models.Card{
				Number:    "5500126132422715",
				Cvv:       "123",
				ExpiredAt: "12/26",
			},
		},
		{
			name:   "#2 nok validation error number",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPost,
			body: models.Card{
				Number:    "5500926132422715",
				Cvv:       "123",
				ExpiredAt: "12/26",
			},
		},
		{
			name:   "#3 nok validation error cvv",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPost,
			body: models.Card{
				Number:    "5500926132422715",
				Cvv:       "12",
				ExpiredAt: "12/26",
			},
		},
		{
			name:   "#4 nok validation error expired",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPost,
			body: models.Card{
				Number:    "5500926132422715",
				Cvv:       "123",
				ExpiredAt: "1226",
			},
		},

		{
			name:   "#5 nok card save err",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodPost,
			body: models.Card{
				Number:    "5500126132422715",
				Cvv:       "123",
				ExpiredAt: "12/26",
			},
			mockBehaviorCreateService: func(s *mock2.MockcardService, card models.Card) {
				s.EXPECT().Create(gomock.Any(), card).Return(errors.New("save card err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcardService(ctrl)
			if test.mockBehaviorCreateService != nil {
				test.mockBehaviorCreateService(service, test.body)
			}

			handler := New(logger.New(""), WithCardUseService(service))

			body, err := json.Marshal(test.body)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, bytes.NewReader(body))
			w := httptest.NewRecorder()
			h := handler.card.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestCardUpdate(t *testing.T) {

	tests := []struct {
		name                      string
		url                       string
		want                      int
		method                    string
		mockBehaviorUpdateService func(s *mock2.MockcardService, card models.Card)
		body                      models.Card
	}{
		{
			name:   "#1 ok update",
			url:    "/",
			want:   http.StatusOK,
			method: http.MethodPut,
			body: models.Card{
				ID:        1,
				Number:    "5500126132422715",
				Cvv:       "123",
				ExpiredAt: "12/26",
				Version:   12,
			},
			mockBehaviorUpdateService: func(s *mock2.MockcardService, card models.Card) {
				s.EXPECT().Update(gomock.Any(), card).Return(nil)
			},
		},

		{
			name:   "#2 nok get id card err",
			url:    "/qwe",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
		},

		{
			name:   "#3 nok validation error number",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.Card{
				Number:    "5500926132422715",
				Cvv:       "123",
				ExpiredAt: "12/26",
			},
		},
		{
			name:   "#4 nok validation error cvv",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.Card{
				Number:    "5500926132422715",
				Cvv:       "12",
				ExpiredAt: "12/26",
			},
		},
		{
			name:   "#5 nok validation error expired",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.Card{
				Number:    "5500926132422715",
				Cvv:       "123",
				ExpiredAt: "1226",
			},
		},

		{
			name:   "#5 nok validation error version",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.Card{
				ID:        1,
				Number:    "5500126132422715",
				Cvv:       "123",
				ExpiredAt: "12/26",
				Version:   0,
			},
		},

		{
			name:   "#6 nok card save err",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodPut,
			body: models.Card{
				ID:        1,
				Number:    "5500126132422715",
				Cvv:       "123",
				ExpiredAt: "12/26",
				Version:   12,
			},
			mockBehaviorUpdateService: func(s *mock2.MockcardService, card models.Card) {
				s.EXPECT().Update(gomock.Any(), card).Return(errors.New("save repo err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcardService(ctrl)
			if test.mockBehaviorUpdateService != nil {
				test.mockBehaviorUpdateService(service, test.body)
			}

			handler := New(logger.New(""), WithCardUseService(service))

			body, err := json.Marshal(test.body)
			assert.NoError(t, err, "Body write error")

			path := fmt.Sprintf("%s%d", test.url, test.body.ID)
			request := httptest.NewRequest(test.method, path, bytes.NewReader(body))
			w := httptest.NewRecorder()
			h := handler.card.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want, res.StatusCode)

		})
	}
}

func TestCardDelete(t *testing.T) {
	tests := []struct {
		name                      string
		want                      int
		method                    string
		url                       string
		cardIndex                 int
		mockBehaviorDeleteService func(s *mock2.MockcardService, cardID int)
	}{
		{
			name:      "#1 ok delete",
			want:      http.StatusOK,
			method:    http.MethodDelete,
			url:       "/",
			cardIndex: 5,
			mockBehaviorDeleteService: func(s *mock2.MockcardService, cardID int) {
				s.EXPECT().Delete(gomock.Any(), cardID).Return(nil)

			},
		},

		{
			name:      "#2 nok card delete err",
			want:      http.StatusInternalServerError,
			method:    http.MethodDelete,
			url:       "/",
			cardIndex: 5,
			mockBehaviorDeleteService: func(s *mock2.MockcardService, cardID int) {
				s.EXPECT().Delete(gomock.Any(), cardID).Return(errors.New("card delete error"))

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcardService(ctrl)
			if test.mockBehaviorDeleteService != nil {
				test.mockBehaviorDeleteService(service, test.cardIndex)
			}

			url := fmt.Sprintf("%s%d", test.url, test.cardIndex)

			handler := New(logger.New(""), WithCardUseService(service))

			r := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()

			handler.card.createRoutes().ServeHTTP(w, r)
			result := w.Result()
			defer w.Result().Body.Close()
			assert.Equal(t, test.want, result.StatusCode)
		})
	}
}
