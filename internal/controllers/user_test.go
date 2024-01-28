package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/models"
	"github.com/zelas91/goph-keeper/internal/repository"
	"github.com/zelas91/goph-keeper/internal/repository/entities"
	"github.com/zelas91/goph-keeper/internal/service"
	mock "github.com/zelas91/goph-keeper/internal/service/mocks"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockBehaviorCreate func(s *mock.MockuserRepo, login, password string)
type mockBehaviorGet func(s *mock.MockuserRepo, user models.User)

type eqCreateUserParamsMatcher struct {
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {

	err := bcrypt.CompareHashAndPassword([]byte(x.(string)), []byte(e.password))
	return err == nil
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("password %v", e.password)
}
func TestSignUp(t *testing.T) {
	tests := []struct {
		name                   string
		want                   int
		url                    string
		content                string
		method                 string
		mockBehaviorCreateUser mockBehaviorCreate
		mockBehaviorGetUser    mockBehaviorGet
		login                  string
		password               string
		user                   models.User
	}{
		{
			name:    "#1 register OK",
			want:    http.StatusOK,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, login, password string) {
				s.EXPECT().CreateUser(gomock.Any(), login, eqCreateUserParamsMatcher{password: password}).Return(nil)
			},
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user models.User) {
				hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				if err != nil {
					return
				}
				us := entities.User{
					ID:       1,
					Login:    user.Login,
					Password: string(hash),
				}
				s.EXPECT().GetUser(gomock.Any(), user).Return(us, nil)
			},
			login:    "user",
			password: "12345678",
			user: models.User{
				Login:    "user",
				Password: "12345678"},
		},

		{
			name:     "#2 register bad request (validation)",
			want:     http.StatusBadRequest,
			method:   http.MethodPost,
			url:      "/api/signup",
			content:  "application/json",
			login:    "user",
			password: "12345678",
			user:     models.User{},
		},
		{
			name:    "#3 register conflict",
			want:    http.StatusConflict,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, login, password string) {
				s.EXPECT().CreateUser(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(repository.ErrDuplicate)
			},
			login:    "user",
			password: "12345678",
			user: models.User{
				Login:    "user",
				Password: "12345678",
			},
		},
		{
			name:    "#4 register Internal Server Error",
			want:    http.StatusInternalServerError,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, login, password string) {
				s.EXPECT().CreateUser(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(sql.ErrNoRows)
			},
			login:    "user",
			password: "12345678",
			user: models.User{
				Login:    "user",
				Password: "12345678",
			},
		},

		{
			name:    "#5 register Unauthorized",
			want:    http.StatusUnauthorized,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, login, password string) {
				s.EXPECT().CreateUser(gomock.Any(), login, eqCreateUserParamsMatcher{password: password}).Return(nil)
			},
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user models.User) {

				s.EXPECT().GetUser(gomock.Any(), user).Return(entities.User{}, sql.ErrNoRows)
			},
			login:    "user",
			password: "12345678",
			user: models.User{
				Login:    "user",
				Password: "12345678",
			},
		},
		{
			name:   "#6 register media type unsupported ",
			want:   http.StatusUnsupportedMediaType,
			method: http.MethodPost,
			url:    "/api/signup",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock.NewMockuserRepo(ctrl)
			if test.mockBehaviorCreateUser != nil {
				test.mockBehaviorCreateUser(repo, test.login, test.password)
			}
			if test.mockBehaviorGetUser != nil {
				test.mockBehaviorGetUser(repo, test.user)
			}

			handler := New(logger.New(), WithAuthUseService(service.New(service.WithAuthUseRepository(repo))))

			body, err := json.Marshal(test.user)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
			w := httptest.NewRecorder()

			h := handler.InitRoutes()
			request.Header.Set("Content-Type", test.content)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})

	}
}

func TestSignIn(t *testing.T) {
	tests := []struct {
		name                string
		want                int
		url                 string
		content             string
		method              string
		mockBehaviorGetUser mockBehaviorGet
		user                models.User
	}{
		{
			name:    "#1 OK authorization",
			want:    http.StatusOK,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user models.User) {
				hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				if err != nil {
					return
				}
				us := entities.User{
					ID:       1,
					Login:    user.Login,
					Password: string(hash),
				}
				s.EXPECT().GetUser(gomock.Any(), user).Return(us, nil)
			},
			user: models.User{
				Login:    "user",
				Password: "12345678",
			},
		},
		{
			name:    "#2 bad request validation",
			want:    http.StatusBadRequest,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			user: models.User{
				Login: "user",
			},
		},
		{
			name:    "#3 Unauthorized",
			want:    http.StatusUnauthorized,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user models.User) {

				s.EXPECT().GetUser(gomock.Any(), user).Return(entities.User{}, sql.ErrNoRows)
			},
			user: models.User{
				Login:    "user",
				Password: "12345678",
			},
		},
		{
			name:   "#4 register media type unsupported ",
			want:   http.StatusUnsupportedMediaType,
			method: http.MethodPost,
			url:    "/api/signin",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock.NewMockuserRepo(ctrl)

			if test.mockBehaviorGetUser != nil {
				test.mockBehaviorGetUser(repo, test.user)
			}
			handler := New(logger.New(), WithAuthUseService(service.New(service.WithAuthUseRepository(repo))))

			body, err := json.Marshal(test.user)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", test.content)
			h := handler.InitRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}
