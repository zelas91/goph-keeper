package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/goph-keeper/internal/logger"
	mock2 "github.com/zelas91/goph-keeper/internal/server/controllers/mocks"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/service"
	"github.com/zelas91/goph-keeper/internal/server/service/mocks"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockBehaviorCreate func(s *mock.MockuserRepo, user entities.User)
type mockBehaviorGet func(s *mock.MockuserRepo, user entities.User)

type eqUserMatcher struct {
	login    string
	password string
}

func (eq eqUserMatcher) Matches(x interface{}) bool {
	user, ok := x.(entities.User)
	if !ok {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(eq.password))
	return err == nil
}

func (eq eqUserMatcher) String() string {
	return fmt.Sprintf("login: %s, password: %s", eq.login, eq.password)
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
		user                   entities.User
	}{
		{
			name:    "#1 register OK",
			want:    http.StatusOK,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, user entities.User) {
				eq := eqUserMatcher{login: user.Login, password: user.Password}
				s.EXPECT().CreateUser(gomock.Any(), eq).Return(nil)
			},
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user entities.User) {
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
			user: entities.User{
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
			user:     entities.User{},
		},
		{
			name:    "#3 register conflict",
			want:    http.StatusConflict,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, user entities.User) {

				s.EXPECT().CreateUser(gomock.Any(),
					gomock.Any()).Return(repository.ErrDuplicate)
			},
			login:    "user",
			password: "12345678",
			user: entities.User{
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
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, user entities.User) {
				s.EXPECT().CreateUser(gomock.Any(),
					gomock.Any()).Return(sql.ErrNoRows)
			},
			login:    "user",
			password: "12345678",
			user: entities.User{
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
			mockBehaviorCreateUser: func(s *mock.MockuserRepo, user entities.User) {
				eq := eqUserMatcher{login: user.Login, password: user.Password}
				s.EXPECT().CreateUser(gomock.Any(), eq).Return(nil)
			},
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user entities.User) {

				s.EXPECT().GetUser(gomock.Any(), user).Return(entities.User{}, sql.ErrNoRows)
			},
			login:    "user",
			password: "12345678",
			user: entities.User{
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
				test.mockBehaviorCreateUser(repo, entities.User{Login: test.login, Password: test.password})
			}
			if test.mockBehaviorGetUser != nil {
				test.mockBehaviorGetUser(repo, test.user)
			}

			handler := New(logger.New(), WithAuthUseService(service.New(service.WithAuthUseRepository(repo))))

			body, err := json.Marshal(test.user)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
			w := httptest.NewRecorder()

			h := handler.CreateRoutes()
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
		user                entities.User
	}{
		{
			name:    "#1 OK authorization",
			want:    http.StatusOK,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user entities.User) {
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
			user: entities.User{
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
			user: entities.User{
				Login: "user",
			},
		},
		{
			name:    "#3 Unauthorized",
			want:    http.StatusUnauthorized,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			mockBehaviorGetUser: func(s *mock.MockuserRepo, user entities.User) {

				s.EXPECT().GetUser(gomock.Any(), user).Return(entities.User{}, sql.ErrNoRows)
			},
			user: entities.User{
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
			h := handler.CreateRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

type mockBehaviorCreateService func(s *mock2.MockuserService, user models.User)
type mockBehaviorCreateTokenService func(s *mock2.MockuserService, user models.User)

func TestSignInService(t *testing.T) {

	tests := []struct {
		name            string
		want            int
		url             string
		content         string
		method          string
		mockCreateToken mockBehaviorCreateService
		user            models.User
	}{
		{
			name:    "#1 OK authorization",
			want:    http.StatusOK,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			user: models.User{
				Login:    "test",
				Password: "123456789",
			},
			mockCreateToken: func(s *mock2.MockuserService, user models.User) {
				s.EXPECT().CreateToken(gomock.Any(), user).Return("token", nil)
			},
		}, {
			name:    "#2 bad request validation",
			want:    http.StatusBadRequest,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			user: models.User{
				Login: "user",
			},
		}, {
			name:    "#3 Unauthorized",
			want:    http.StatusUnauthorized,
			url:     "/api/signin",
			content: "application/json",
			method:  http.MethodPost,
			mockCreateToken: func(s *mock2.MockuserService, user models.User) {

				s.EXPECT().CreateToken(gomock.Any(), user).Return("", errors.New("no user"))
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
			serv := mock2.NewMockuserService(ctrl)
			if test.mockCreateToken != nil {
				test.mockCreateToken(serv, test.user)
			}

			handler := New(logger.New(), WithAuthUseService(serv))

			body, err := json.Marshal(test.user)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", test.content)
			h := handler.CreateRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}

func TestSignUpService(t *testing.T) {
	tests := []struct {
		name            string
		want            int
		url             string
		content         string
		method          string
		mockCreate      mockBehaviorCreateService
		mockCreateToken mockBehaviorCreateTokenService
		user            models.User
	}{
		{
			name:    "#1 register OK",
			want:    http.StatusOK,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			user: models.User{
				Login:    "user",
				Password: "12345678"},
			mockCreate: func(s *mock2.MockuserService, user models.User) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return(nil)
			},
			mockCreateToken: func(s *mock2.MockuserService, user models.User) {
				s.EXPECT().CreateToken(gomock.Any(), user).Return("token", nil)
			},
		},
		{
			name:    "#2 register bad request (validation)",
			want:    http.StatusBadRequest,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			user:    models.User{Login: "test"},
		},
		{
			name:    "#3 register conflict",
			want:    http.StatusConflict,
			method:  http.MethodPost,
			url:     "/api/signup",
			content: "application/json",
			mockCreate: func(s *mock2.MockuserService, user models.User) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return(repository.ErrDuplicate)
			},
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
			mockCreate: func(s *mock2.MockuserService, user models.User) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return(sql.ErrNoRows)
			},
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
			mockCreate: func(s *mock2.MockuserService, user models.User) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return(nil)
			},
			mockCreateToken: func(s *mock2.MockuserService, user models.User) {
				s.EXPECT().CreateToken(gomock.Any(), user).Return("", errors.New("invalid token"))
			},
			user: models.User{
				Login:    "user",
				Password: "12345678",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serv := mock2.NewMockuserService(ctrl)
			if test.mockCreate != nil {
				test.mockCreate(serv, test.user)
			}
			if test.mockCreateToken != nil {
				test.mockCreateToken(serv, test.user)
			}

			handler := New(logger.New(), WithAuthUseService(serv))

			body, err := json.Marshal(test.user)
			assert.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", test.content)
			h := handler.CreateRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}
