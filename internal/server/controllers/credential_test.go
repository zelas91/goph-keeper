package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/goph-keeper/internal/logger"
	mock2 "github.com/zelas91/goph-keeper/internal/server/controllers/mocks"
	"github.com/zelas91/goph-keeper/internal/server/models"
)

func TestCredentials(t *testing.T) {

	tests := []struct {
		name                           string
		url                            string
		want                           int
		method                         string
		mockBehaviorCredentialsService func(s *mock2.MockcredentialService)
	}{
		{
			name:   "#1 ok get credentials",
			url:    "/",
			want:   http.StatusOK,
			method: http.MethodGet,
			mockBehaviorCredentialsService: func(s *mock2.MockcredentialService) {
				s.EXPECT().Credentials(gomock.Any()).Return([]models.UserCredentials{
					{
						Login:    "test",
						Password: "12345678",
					},
					{
						Login:    "users",
						Password: "qwertyuiop",
					},
					{
						Login:    "ololo",
						Password: "4567890qwe",
					},
					{
						Login:    "montgomery",
						Password: "CC771212cC",
					},
				}, nil)
			},
		},
		{
			name:   "#2 nok get credentials error",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodGet,
			mockBehaviorCredentialsService: func(s *mock2.MockcredentialService) {
				s.EXPECT().Credentials(gomock.Any()).Return(nil, errors.New("repo error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcredentialService(ctrl)
			if test.mockBehaviorCredentialsService != nil {
				test.mockBehaviorCredentialsService(service)
			}

			handler := New(logger.New(""), WithUserCredentialUseService(service))

			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := handler.credential.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestCredential(t *testing.T) {

	tests := []struct {
		name                          string
		url                           string
		want                          int
		method                        string
		credentialIndex               int
		mockBehaviorCredentialService func(s *mock2.MockcredentialService, credentialID int)
	}{
		{
			name:            "#1 ok get  credential",
			credentialIndex: 2,
			url:             "/",
			want:            http.StatusOK,
			method:          http.MethodGet,
			mockBehaviorCredentialService: func(s *mock2.MockcredentialService, credentialID int) {
				s.EXPECT().Credential(gomock.Any(), credentialID).Return(models.UserCredentials{
					ID:       credentialID,
					Login:    "montgomery",
					Password: "CC771212cC",
					Version:  1,
				}, nil)
			},
		},
		{
			name:   "#2 nok credential get id from request err",
			url:    "/a",
			want:   http.StatusBadRequest,
			method: http.MethodGet,
		},
		{
			name:            "#3 nok credential: get credential err",
			credentialIndex: 2,
			url:             "/",
			want:            http.StatusNotFound,
			method:          http.MethodGet,
			mockBehaviorCredentialService: func(s *mock2.MockcredentialService, credentialID int) {
				s.EXPECT().Credential(gomock.Any(), credentialID).Return(models.UserCredentials{},
					errors.New("credential: get credential err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcredentialService(ctrl)
			if test.mockBehaviorCredentialService != nil {
				test.mockBehaviorCredentialService(service, test.credentialIndex)
			}

			handler := New(logger.New(""), WithUserCredentialUseService(service))

			url := fmt.Sprintf("%s%d", test.url, test.credentialIndex)
			request := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()
			h := handler.credential.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestCredentialCreate(t *testing.T) {

	tests := []struct {
		name                      string
		url                       string
		want                      int
		method                    string
		mockBehaviorCreateService func(s *mock2.MockcredentialService, credential models.UserCredentials)
		body                      models.UserCredentials
	}{
		{
			name:   "#1 ok create credential",
			url:    "/",
			want:   http.StatusCreated,
			method: http.MethodPost,
			mockBehaviorCreateService: func(s *mock2.MockcredentialService, credential models.UserCredentials) {
				s.EXPECT().Create(gomock.Any(), credential).Return(nil)
			},
			body: models.UserCredentials{
				Login:    "montgomery",
				Password: "Cc771212cC",
			},
		},
		{
			name:   "#2 nok validation error login",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPost,
			body: models.UserCredentials{
				Login:    "Use",
				Password: "123456789",
			},
		},
		{
			name:   "#3 nok validation error password",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPost,
			body: models.UserCredentials{
				Login:    "montgomery",
				Password: "123456",
			},
		},

		{
			name:   "#5 nok credential save err",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodPost,
			body: models.UserCredentials{
				Login:    "montgomery",
				Password: "12345678",
			},
			mockBehaviorCreateService: func(s *mock2.MockcredentialService, credential models.UserCredentials) {
				s.EXPECT().Create(gomock.Any(), credential).Return(errors.New("save credential err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcredentialService(ctrl)
			if test.mockBehaviorCreateService != nil {
				test.mockBehaviorCreateService(service, test.body)
			}

			handler := New(logger.New(""), WithUserCredentialUseService(service))

			body, err := json.Marshal(test.body)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, bytes.NewReader(body))
			w := httptest.NewRecorder()
			h := handler.credential.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestCredentialUpdate(t *testing.T) {

	tests := []struct {
		name                      string
		url                       string
		want                      int
		method                    string
		mockBehaviorUpdateService func(s *mock2.MockcredentialService, credential models.UserCredentials)
		body                      models.UserCredentials
	}{
		{
			name:   "#1 ok update",
			url:    "/",
			want:   http.StatusOK,
			method: http.MethodPut,
			body: models.UserCredentials{
				ID:       1,
				Login:    "montgomery",
				Password: "Cc771212Cc",
				Version:  3,
			},
			mockBehaviorUpdateService: func(s *mock2.MockcredentialService, credential models.UserCredentials) {
				s.EXPECT().Update(gomock.Any(), credential).Return(nil)
			},
		},

		{
			name:   "#2 nok get id credential err",
			url:    "/qwe",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
		},

		{
			name:   "#3 nok validation error login",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.UserCredentials{
				ID:       2,
				Login:    "mon",
				Password: "12345678",
				Version:  5,
			},
		},
		{
			name:   "#4 nok validation error password",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.UserCredentials{
				ID:       2,
				Login:    "montgomery",
				Password: "1234",
				Version:  5,
			},
		},
		{
			name:   "#5 nok validation error version",
			url:    "/",
			want:   http.StatusBadRequest,
			method: http.MethodPut,
			body: models.UserCredentials{
				ID:       2,
				Login:    "montgomery",
				Password: "1234",
				Version:  0,
			},
		},

		{
			name:   "#6 nok credential save err",
			url:    "/",
			want:   http.StatusInternalServerError,
			method: http.MethodPut,
			body: models.UserCredentials{
				ID:       2,
				Login:    "montgomery",
				Password: "12345678",
				Version:  9,
			},
			mockBehaviorUpdateService: func(s *mock2.MockcredentialService, credential models.UserCredentials) {
				s.EXPECT().Update(gomock.Any(), credential).Return(errors.New("save repo err"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcredentialService(ctrl)
			if test.mockBehaviorUpdateService != nil {
				test.mockBehaviorUpdateService(service, test.body)
			}

			handler := New(logger.New(""), WithUserCredentialUseService(service))

			body, err := json.Marshal(test.body)
			assert.NoError(t, err, "Body write error")

			path := fmt.Sprintf("%s%d", test.url, test.body.ID)
			request := httptest.NewRequest(test.method, path, bytes.NewReader(body))
			w := httptest.NewRecorder()
			h := handler.credential.createRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want, res.StatusCode)

		})
	}
}

func TestCredentialDelete(t *testing.T) {
	tests := []struct {
		name                      string
		want                      int
		method                    string
		url                       string
		credentialIndex           int
		mockBehaviorDeleteService func(s *mock2.MockcredentialService, credentialID int)
	}{
		{
			name:            "#1 ok delete",
			want:            http.StatusOK,
			method:          http.MethodDelete,
			url:             "/",
			credentialIndex: 9,
			mockBehaviorDeleteService: func(s *mock2.MockcredentialService, credentialID int) {
				s.EXPECT().Delete(gomock.Any(), credentialID).Return(nil)

			},
		},

		{
			name:            "#2 nok credential delete err",
			want:            http.StatusInternalServerError,
			method:          http.MethodDelete,
			url:             "/",
			credentialIndex: 9,
			mockBehaviorDeleteService: func(s *mock2.MockcredentialService, credentialID int) {
				s.EXPECT().Delete(gomock.Any(), credentialID).Return(errors.New("credential delete error"))

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock2.NewMockcredentialService(ctrl)
			if test.mockBehaviorDeleteService != nil {
				test.mockBehaviorDeleteService(service, test.credentialIndex)
			}

			url := fmt.Sprintf("%s%d", test.url, test.credentialIndex)

			handler := New(logger.New(""), WithUserCredentialUseService(service))

			r := httptest.NewRequest(test.method, url, nil)
			w := httptest.NewRecorder()

			handler.credential.createRoutes().ServeHTTP(w, r)
			result := w.Result()
			defer w.Result().Body.Close()
			assert.Equal(t, test.want, result.StatusCode)
		})
	}
}
