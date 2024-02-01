package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	mock "github.com/zelas91/goph-keeper/internal/server/services/mocks"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
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

func TestParserToken(t *testing.T) {

}
func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		user       entities.User
		want       error
		mockCreate mockBehaviorCreate
	}{
		{
			name: "#1 create user ok",
			user: entities.User{
				Login:    "test",
				Password: "123456789",
			},
			mockCreate: func(s *mock.MockuserRepo, user entities.User) {
				eq := eqUserMatcher{login: user.Login, password: user.Password}
				s.EXPECT().Create(gomock.Any(), eq).Return(nil)
			},
			want: nil,
		},
		{
			name: "#2 create user nok(duplicate)",
			user: entities.User{
				Login:    "test",
				Password: "12345678",
			},
			mockCreate: func(s *mock.MockuserRepo, user entities.User) {
				eq := eqUserMatcher{login: user.Login, password: user.Password}
				s.EXPECT().Create(gomock.Any(), eq).Return(repository.ErrDuplicate)
			},
			want: repository.ErrDuplicate,
		},
		{
			name: "#3 create user nok",
			user: entities.User{
				Login:    "test",
				Password: "12345678",
			},
			mockCreate: func(s *mock.MockuserRepo, user entities.User) {
				eq := eqUserMatcher{login: user.Login, password: user.Password}
				s.EXPECT().Create(gomock.Any(), eq).Return(sql.ErrNoRows)
			},
			want: sql.ErrNoRows,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock.NewMockuserRepo(ctrl)
			if test.mockCreate != nil {
				test.mockCreate(repo, test.user)
			}
			serv := New(WithAuthUseRepository(repo))
			res := serv.Auth.CreateUser(context.TODO(), models.User{Login: test.user.Login, Password: test.user.Password})
			fmt.Println(res)
			assert.Equal(t, test.want, res)
		})
	}
}
func TestCreateToken(t *testing.T) {
	type want struct {
		login string
		err   error
	}
	tests := []struct {
		name    string
		user    entities.User
		want    want
		mockGet mockBehaviorGet
	}{
		{
			name: "#1 create token ok",
			user: entities.User{
				Login:    "test",
				Password: "test",
			},
			want: want{
				login: "test",
				err:   nil,
			},
			mockGet: func(s *mock.MockuserRepo, user entities.User) {
				hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				if err != nil {
					return
				}
				us := entities.User{
					ID:       1,
					Login:    user.Login,
					Password: string(hash),
				}
				s.EXPECT().FindUserByLogin(gomock.Any(), user).Return(us, nil)
			},
		},
		{
			name: "#2 create token user not found",
			user: entities.User{
				Login:    "test",
				Password: "test",
			},
			want: want{
				login: "test",
				err:   errors.New("user not found"),
			},
			mockGet: func(s *mock.MockuserRepo, user entities.User) {
				s.EXPECT().FindUserByLogin(gomock.Any(), user).Return(entities.User{}, errors.New("user not found"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock.NewMockuserRepo(ctrl)
			if test.mockGet != nil {
				test.mockGet(repo, test.user)
			}
			serv := New(WithAuthUseRepository(repo))
			tokenStr, err := serv.Auth.CreateToken(context.TODO(), models.User{Login: test.user.Login,
				Password: test.user.Password})

			assert.Equal(t, test.want.err, err)

			claims := &Claims{}

			_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error parsing jwt")
				}
				return secret, nil
			})
			if err == nil {
				assert.Equal(t, test.want.login, claims.Login)
			}

		})
	}
}
